package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/arkiaconsulting/terraform-provider-akc/utils"
)

var (
	LabelNone     = "%00"
	FeaturePrefix = ".appconfig.featureflag/"
)

func toPrefixedFeature(key string) string {
	return fmt.Sprintf("%s%s", FeaturePrefix, key)
}

var (
	defaultContentType     = "application/vnd.microsoft.appconfig.kv+json"
	keyVaultRefContentType = "application/vnd.microsoft.appconfig.keyvaultref+json;charset=utf-8"
	featureContentType     = "application/vnd.microsoft.appconfig.ff+json;charset=utf-8"
)

type Client struct {
	*autorest.Client
	Endpoint string
}

type setKeyValuePayload struct {
	Value       string
	ContentType string `json:"content_type"`
}

type featurePayload struct {
	Id         string                   `json:"id"`
	Descripton string                   `json:"description"`
	Enabled    bool                     `json:"enabled"`
	Conditions featureConditionsPayload `json:"conditions"`
}

type featureConditionsPayload struct {
	ClientFilters []string `json:"client_filters"`
}

func NewClientCreds(endpoint string, clientID string, clientSecret string, tenantID string) (*Client, error) {

	ccc := auth.NewClientCredentialsConfig(clientID, clientSecret, tenantID)
	ccc.Resource = endpoint

	authorizer, err := ccc.Authorizer()
	if err != nil {
		return nil, err
	}
	return NewClient(endpoint, authorizer)
}

func NewClientCli(endpoint string) (*Client, error) {
	authorizer, err := auth.NewAuthorizerFromCLIWithResource(endpoint)
	if err != nil {
		return nil, err
	}
	return NewClient(endpoint, authorizer)
}

func NewClientMsi(endpoint string) (*Client, error) {
	msiConfig := auth.NewMSIConfig()
	msiConfig.Resource = endpoint

	authorizer, err := msiConfig.Authorizer()
	if err != nil {
		return nil, err
	}
	return NewClient(endpoint, authorizer)
}

func NewClient(endpoint string, authorizer autorest.Authorizer) (*Client, error) {
	client := autorest.NewClientWithUserAgent(userAgent())
	client.Authorizer = authorizer

	return &Client{
		Client:   &client,
		Endpoint: endpoint,
	}, nil
}

func (client *Client) GetKeyValue(label string, key string) (KeyValueResponse, error) {
	result := KeyValueResponse{}
	resp, err := client.send(
		label,
		url.QueryEscape(key),
		true,
		autorest.AsGet(),
	)

	if err != nil {
		return result, err
	}

	if utils.ResponseWasNotFound(resp) {
		return result, KVNotFoundError.with(key)
	}

	if err = getJSON(resp, &result); err != nil {
		return result, UnexpectedError.wrap(err)
	}

	return result, nil
}

func (client *Client) SetKeyValue(label string, key string, value string) (KeyValueResponse, error) {
	return client.setKeyValue(label, key, value, defaultContentType)
}

func (client *Client) SetKeyValueSecret(key string, secretID string, label string) (KeyValueResponse, error) {
	value := fmt.Sprintf("{\"uri\":\"%s\"}", secretID)
	return client.setKeyValue(label, key, value, keyVaultRefContentType)
}

func (client *Client) SetFeature(key string, label string, enabled bool, description string) (KeyValueResponse, error) {
	actualKey := fmt.Sprintf(".appconfig.featureflag%%2F%s", key)

	featurePayload := featurePayload{
		Id:         key,
		Descripton: description,
		Enabled:    enabled,
		Conditions: featureConditionsPayload{
			ClientFilters: []string{},
		},
	}

	b, err := json.Marshal(featurePayload)
	if err != nil {
		return KeyValueResponse{}, UnexpectedError.wrap(err)
	}

	return client.setKeyValue(label, actualKey, string(b), featureContentType)
}

func (client *Client) GetFeature(label string, key string) (FeatureResponse, error) {
	kvResponse, err := client.GetKeyValue(label, toPrefixedFeature(key))
	if err != nil {
		return FeatureResponse{}, err
	}

	details := featurePayload{}
	err = json.Unmarshal([]byte(kvResponse.Value), &details)
	if err != nil {
		return FeatureResponse{}, UnexpectedError.wrap(err)
	}

	resp := FeatureResponse{
		Key:          strings.TrimPrefix(kvResponse.Key, FeaturePrefix),
		Label:        kvResponse.Label,
		ContentType:  kvResponse.ContentType,
		LastModified: kvResponse.LastModified,
		Tags:         kvResponse.Tags,
		Description:  details.Descripton,
		Enabled:      details.Enabled,
	}

	return resp, nil
}

func (client *Client) DeleteFeature(label string, key string) (bool, error) {
	return client.DeleteKeyValue(label, toPrefixedFeature(key))
}

func (client *Client) setKeyValue(label string, key string, value string, contentType string) (KeyValueResponse, error) {
	result := KeyValueResponse{}
	payload := setKeyValuePayload{
		Value:       value,
		ContentType: contentType,
	}

	resp, err := client.send(
		label,
		key,
		false,
		autorest.AsContentType(defaultContentType),
		autorest.AsPut(),
		autorest.WithJSON(payload),
	)
	if err != nil {
		return result, err
	}

	if utils.ResponseWasStatusCode(resp, http.StatusForbidden) {
		return result, UnexpectedError.wrap(fmt.Errorf("Forbidden"))
	}

	err = getJSON(resp, &result)
	if err != nil {
		return result, UnexpectedError.wrap(err)
	}

	return result, nil
}

func (client *Client) DeleteKeyValue(label string, key string) (bool, error) {
	resp, err := client.send(
		label,
		url.QueryEscape(key),
		false,
		autorest.AsDelete(),
	)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return false, nil
	}

	return true, nil
}

func userAgent() string {
	return "Azure-SDK-For-Go/" + "1.0" + " akc-key-value/2020-03-01"
}

func getJSON(response *http.Response, target interface{}) error {
	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)

	err := dec.Decode(target)

	return err
}

func (client *Client) send(label string, key string, retry bool, additionalDecorator ...autorest.PrepareDecorator) (*http.Response, error) {
	req, err := client.getPreparer(
		label,
		key,
		additionalDecorator...,
	).Prepare(&http.Request{})

	if err != nil {
		return nil, UnexpectedError.wrap(err)
	}

	var resp *http.Response
	var retryDecorator autorest.SendDecorator
	if retry {
		retryDecorator = autorest.DoRetryForStatusCodes(5, 5*time.Second/10, 404)
		resp, err = client.Send(req, retryDecorator)
	} else {
		resp, err = client.Send(req)
	}

	if err != nil {
		return nil, UnexpectedError.wrap(err)
	}

	if utils.ResponseWasThrottled(resp) {
		return nil, UnexpectedError.with("Requests are throttled")
	}

	if utils.ResponseWasForbidden(resp) {
		return nil, UnexpectedError.with("Forbidden")
	}

	if utils.ResponseWasUnauthorized(resp) {
		return nil, UnexpectedError.with("Unauthorized")
	}

	return resp, err
}

func (client *Client) getPreparer(label string, key string, additionalDecorators ...autorest.PrepareDecorator) autorest.Preparer {
	const apiVersion = "1.0"
	queryParameters := map[string]interface{}{
		"label":       label,
		"api-version": apiVersion,
	}

	pathParameters := map[string]interface{}{
		"key": key,
	}

	decorators := []autorest.PrepareDecorator{
		autorest.WithBaseURL(client.Endpoint),
		autorest.WithPathParameters("/kv/{key}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		client.Client.WithAuthorization(),
	}

	decorators = append(decorators, additionalDecorators...)

	return autorest.CreatePreparer(decorators...)
}

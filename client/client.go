package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-akc/utils"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var (
	LabelNone = "%00"
)

var (
	defaultContentType     = "application/vnd.microsoft.appconfig.kv+json"
	keyVaultRefContentType = "application/vnd.microsoft.appconfig.keyvaultref+json;charset=utf-8"
)

type Client struct {
	*autorest.Client
	Endpoint string
}

type setKeyValuePayload struct {
	Value       string
	ContentType string `json:"content_type"`
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

func NewClientEnv(endpoint string) (*Client, error) {
	authorizer, err := auth.NewAuthorizerFromEnvironmentWithResource(endpoint)
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
		key,
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

func (client *Client) setKeyValue(label string, key string, value string, contentType string) (KeyValueResponse, error) {
	result := KeyValueResponse{}
	payload := setKeyValuePayload{
		Value:       value,
		ContentType: contentType,
	}

	resp, err := client.send(
		label,
		key,
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
		key,
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

func (client *Client) send(label string, key string, additionalDecorator ...autorest.PrepareDecorator) (*http.Response, error) {
	req, err := client.getPreparer(
		label,
		key,
		additionalDecorator...,
	).Prepare(&http.Request{})

	if err != nil {
		return nil, UnexpectedError.wrap(err)
	}

	resp, err := client.Send(req)
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

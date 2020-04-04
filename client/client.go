package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/arkiaconsulting/terraform-provider-akc/utils"
)

var (
	// LabelNone Represents an empty label
	LabelNone = "%00"
)

var (
	defaultContentType     = "application/vnd.microsoft.appconfig.kv+json"
	keyVaultRefContentType = "application/vnd.microsoft.appconfig.keyvaultref+json;charset=utf-8"
)

// AppConfigClient holds all of the information required to connect to a server
type AppConfigClient struct {
	*autorest.Client
	AppConfigURI string
}

type setKeyValuePayload struct {
	Value       string
	ContentType string `json:"content_type"`
}

// NewAppConfigurationClient instantiate a new client
func NewAppConfigurationClient(appConfigURI string) (*AppConfigClient, error) {
	client := autorest.NewClientWithUserAgent(userAgent())
	authorizer, err := getAuthorizer(appConfigURI)
	if err != nil {
		return nil, err
	}

	client.Authorizer = authorizer

	return &AppConfigClient{
		Client:       &client,
		AppConfigURI: appConfigURI,
	}, nil
}

// GetKeyValue get a given App Configuration key-value
func (client *AppConfigClient) GetKeyValue(label string, key string) (KeyValueResponse, error) {
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

// SetKeyValue creates a key with the given value
func (client *AppConfigClient) SetKeyValue(label string, key string, value string) (KeyValueResponse, error) {
	return client.setKeyValue(label, key, value, defaultContentType)
}

// SetKeyValueSecret creates a key with the given KeyVault secret ID
func (client *AppConfigClient) SetKeyValueSecret(key string, secretID string, label string) (KeyValueResponse, error) {
	value := fmt.Sprintf("{\"uri\":\"%s\"}", secretID)
	return client.setKeyValue(label, key, value, keyVaultRefContentType)
}

func (client *AppConfigClient) setKeyValue(label string, key string, value string, contentType string) (KeyValueResponse, error) {
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

// DeleteKeyValue get a given App Configuration key-value
func (client *AppConfigClient) DeleteKeyValue(label string, key string) (bool, error) {
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

func (client *AppConfigClient) send(label string, key string, additionalDecorator ...autorest.PrepareDecorator) (*http.Response, error) {
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

func getAuthorizer(resource string) (autorest.Authorizer, error) {
	authorizer, err := auth.NewAuthorizerFromCLIWithResource(resource)
	if err != nil {
		authorizer, err = auth.NewAuthorizerFromEnvironmentWithResource(resource)
		if err != nil {
			return nil, fmt.Errorf("Unable to find a suitable authorizer")
		}
	}

	return authorizer, nil
}

func (client *AppConfigClient) getPreparer(label string, key string, additionalDecorators ...autorest.PrepareDecorator) autorest.Preparer {
	const apiVersion = "1.0"
	queryParameters := map[string]interface{}{
		"label":       label,
		"api-version": apiVersion,
	}

	pathParameters := map[string]interface{}{
		"key": key,
	}

	decorators := []autorest.PrepareDecorator{
		autorest.WithBaseURL(client.AppConfigURI),
		autorest.WithPathParameters("/kv/{key}", pathParameters),
		autorest.WithQueryParameters(queryParameters),
		client.Client.WithAuthorization(),
	}

	decorators = append(decorators, additionalDecorators...)

	return autorest.CreatePreparer(decorators...)
}

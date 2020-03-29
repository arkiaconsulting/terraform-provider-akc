package client

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/arkiaconsulting/terraform-provider-akc/utils"
)

// AppConfigClient holds all of the information required to connect to a server
type AppConfigClient struct {
	*autorest.Client
	AppConfigURI string
}

type setKeyValuePayload struct {
	Value string
}

// NewAppConfigurationClient instantiate a new client
func NewAppConfigurationClient(appConfigURI string) *AppConfigClient {
	client := autorest.NewClientWithUserAgent(userAgent())
	client.Authorizer, _ = auth.NewAuthorizerFromEnvironmentWithResource(appConfigURI)

	return &AppConfigClient{
		Client:       &client,
		AppConfigURI: appConfigURI,
	}
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

// GetKeyValue get a given App Configuration key-value
func (client *AppConfigClient) GetKeyValue(key string) (KeyValueResponse, error) {
	return client.GetKeyValueWithLabel(key, "%00")
}

// GetKeyValueWithLabel get a given App Configuration key-value with label
func (client *AppConfigClient) GetKeyValueWithLabel(key string, label string) (KeyValueResponse, error) {
	result := KeyValueResponse{}
	resp, err := client.send(
		label,
		key,
		autorest.AsGet(),
	)

	if err != nil {
		return result, UnexpectedError.wrap(err)
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
func (client *AppConfigClient) SetKeyValue(key string, value string) (KeyValueResponse, error) {
	return client.SetKeyValueWithLabel(key, value, "%00")
}

// SetKeyValueWithLabel creates a key with the given value and label
func (client *AppConfigClient) SetKeyValueWithLabel(key string, value string, label string) (KeyValueResponse, error) {
	result := KeyValueResponse{}
	payload := setKeyValuePayload{
		Value: value,
	}

	resp, err := client.send(
		label,
		key,
		autorest.AsContentType("application/vnd.microsoft.appconfig.kv+json"),
		autorest.AsPut(),
		autorest.WithJSON(payload),
	)
	if err != nil {
		return result, UnexpectedError.wrap(err)
	}

	err = getJSON(resp, &result)
	if err != nil {
		return result, UnexpectedError.wrap(err)
	}

	return result, nil
}

// DeleteKeyValue get a given App Configuration key-value
func (client *AppConfigClient) DeleteKeyValue(key string) (bool, error) {
	return client.DeleteKeyValueWithLabel(key, "%00")
}

// DeleteKeyValueWithLabel get a given App Configuration key-value with label
func (client *AppConfigClient) DeleteKeyValueWithLabel(key string, label string) (bool, error) {
	resp, err := client.send(
		label,
		key,
		autorest.AsDelete(),
	)
	if err != nil {
		return false, UnexpectedError.wrap(err)
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

	return resp, err
}

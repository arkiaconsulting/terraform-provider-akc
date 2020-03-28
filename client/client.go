package client

import (
	"encoding/json"
	"log"
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

func (client *AppConfigClient) getPreparer(label string, additionalDecorators ...autorest.PrepareDecorator) autorest.Preparer {
	const apiVersion = "1.0"
	queryParameters := map[string]interface{}{
		"label":       label,
		"api-version": apiVersion,
	}

	decorators := []autorest.PrepareDecorator{
		autorest.WithBaseURL(client.AppConfigURI),
	}

	queryDecorator := []autorest.PrepareDecorator{
		autorest.WithQueryParameters(queryParameters),
	}

	decorators = append(decorators, additionalDecorators...)
	decorators = append(decorators, queryDecorator...)
	decorators = append([]autorest.PrepareDecorator{client.Client.WithAuthorization()}, decorators...)

	return autorest.CreatePreparer(decorators...)
}

// GetKeyValue get a given App Configuration key-value
func (client *AppConfigClient) GetKeyValue(key string) (KeyValueResponse, error) {
	logHere("GetKeyValue", "%00", key)

	return client.GetKeyValueWithLabel(key, "%00")
}

// GetKeyValueWithLabel get a given App Configuration key-value with label
func (client *AppConfigClient) GetKeyValueWithLabel(key string, label string) (KeyValueResponse, error) {
	logHere("GetKeyValueWithLabel", label, key)

	result := KeyValueResponse{}
	pathParameters := map[string]interface{}{
		"key":   key,
		"label": label,
	}

	req, err := client.getPreparer(
		label,
		autorest.AsGet(),
		autorest.WithPathParameters("/kv/{key}", pathParameters),
	).Prepare(&http.Request{})

	if err != nil {
		return result, kVError(label, key, "Fail preparing request")
	}

	resp, err := client.Send(req)
	if err != nil {
		return result, kVError(label, key, "Fail sending request")
	}

	if utils.ResponseWasNotFound(resp) {
		return result, kVError(label, key, "Not found")
	}

	err = getJSON(resp, &result)
	if err != nil {
		return result, kVError(label, key, "Fail reading json response")
	}

	return result, nil
}

// SetKeyValue creates a key with the given value
func (client *AppConfigClient) SetKeyValue(key string, value string) (KeyValueResponse, error) {
	logHere("SetKeyValue", "%00", key)

	return client.SetKeyValueWithLabel(key, value, "%00")
}

// SetKeyValueWithLabel creates a key with the given value and label
func (client *AppConfigClient) SetKeyValueWithLabel(key string, value string, label string) (KeyValueResponse, error) {
	logHere("SetKeyValueWithLabel", label, key)

	result := KeyValueResponse{}
	pathParameters := map[string]interface{}{
		"key": key,
	}

	payload := setKeyValuePayload{
		Value: value,
	}

	req, err := client.getPreparer(
		label,
		autorest.AsContentType("application/vnd.microsoft.appconfig.kv+json"),
		autorest.AsPut(),
		autorest.WithPathParameters("/kv/{key}", pathParameters),
		autorest.WithJSON(payload),
	).Prepare(&http.Request{})
	if err != nil {
		return result, kVError(label, key, "Fail preparing request")
	}

	resp, err := client.Send(req)
	if err != nil {
		return result, kVError(label, key, "Fail sending request")
	}

	err = getJSON(resp, &result)
	if err != nil {
		return result, kVError(label, key, "Fail reading json response")
	}

	return result, nil
}

// DeleteKeyValue get a given App Configuration key-value
func (client *AppConfigClient) DeleteKeyValue(key string) (bool, error) {
	logHere("DeleteKeyValue", "%00", key)

	return client.DeleteKeyValueWithLabel(key, "%00")
}

// DeleteKeyValueWithLabel get a given App Configuration key-value with label
func (client *AppConfigClient) DeleteKeyValueWithLabel(key string, label string) (bool, error) {
	logHere("DeleteKeyValueWithLabel", label, key)

	pathParameters := map[string]interface{}{
		"key": key,
	}

	req, err := client.getPreparer(
		label,
		autorest.AsDelete(),
		autorest.WithPathParameters("/kv/{key}", pathParameters),
	).Prepare(&http.Request{})

	if err != nil {
		return false, kVError(label, key, "Fail preparing request")
	}

	resp, err := client.Send(req)
	if err != nil {
		return false, kVError(label, key, "Fail sending request")
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

func logHere(method string, label string, key string) {
	log.Printf("[INFO] %s: (%s/%s)", method, label, key)
}

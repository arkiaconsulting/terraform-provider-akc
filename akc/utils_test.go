package akc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func preCheck(t *testing.T) {
	log.Printf("[DEBUG] - testPreCheck\n")

	variables := []string{
		"ARM_CLIENT_ID",
		"ARM_CLIENT_SECRET",
		"ARM_SUBSCRIPTION_ID",
		"ARM_TENANT_ID",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

// BuildTerraformConfig build terraform config
func buildTerraformConfigWithoutLabel(key string, value string) string {
	return fmt.Sprintf(`
resource "akc_key_value" "test" {
  endpoint     = "%s"
  key = "%s"
  value = "%s"
}
`, endpointUnderTest, key, value)
}

func buildTerraformConfigWithLabel(label string, key string, value string) string {
	return fmt.Sprintf(`
resource "akc_key_value" "test" {
  endpoint     = "%s"
  label = "%s"
  key = "%s"
  value = "%s"
}
`, endpointUnderTest, label, key, value)
}

func buildTerraformConfigSecret(label string, key string, secretID string) string {
	return fmt.Sprintf(`
resource "akc_key_secret" "test" {
  endpoint     = "%s"
  label = "%s"
  key = "%s"
  secret_id = "%s"
}
`, endpointUnderTest, label, key, secretID)
}

func buildTerraformConfigSecretLatestVersion(label string, key string, secretID string, latestVersion bool) string {
	return fmt.Sprintf(`
resource "akc_key_secret" "test" {
  endpoint     = "%s"
  label = "%s"
  key = "%s"
  secret_id = "%s"
  latest_version = %t
}
`, endpointUnderTest, label, key, secretID, latestVersion)
}

func buildTerraformConfigDataSourceKeyValue(key string) string {
	return fmt.Sprintf(`
data "akc_key_value" "test" {
  endpoint     = "%s"
  key = "%s"
}
`, endpointUnderTest, key)
}

func buildTerraformConfigDataSourceKeyValueLabel(label string, key string) string {
	return fmt.Sprintf(`
data "akc_key_value" "test" {
  endpoint     = "%s"
  label = "%s"
  key = "%s"
}
`, endpointUnderTest, label, key)
}

func buildTerraformConfigDataSourceKeySecret(key string) string {
	return fmt.Sprintf(`
data "akc_key_secret" "test" {
  endpoint     = "%s"
  key = "%s"
}
`, endpointUnderTest, key)
}

func buildTerraformConfigDataSourceKeySecretLabel(label string, key string) string {
	return fmt.Sprintf(`
data "akc_key_secret" "test" {
  endpoint     = "%s"
  label = "%s"
  key = "%s"
}
`, endpointUnderTest, label, key)
}

func buildLabeledFeature(name string, label string, description string, enabled bool) string {
	return fmt.Sprintf(`
resource "akc_feature" "test" {
  endpoint     = "%s"
  name = "%s"
  label = "%s"
  description = "%s"
  enabled = %t
}
`, endpointUnderTest, name, label, description, enabled)
}

func buildFeature(name string, description string, enabled bool) string {
	return fmt.Sprintf(`
resource "akc_feature" "test" {
  endpoint     = "%s"
  name = "%s"
  description = "%s"
  enabled = %t
}
`, endpointUnderTest, name, description, enabled)
}

func testCheckKeyValueDestroy(state *terraform.State) error {
	log.Printf("[INFO] Entering Destroy")
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "akc_key_value" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		key := rs.Primary.Attributes["key"]
		label := rs.Primary.Attributes["label"]
		endpoint := rs.Primary.Attributes["endpoint"]
		log.Printf("[INFO] Checking that KV is destroyed %s/%s/%s", endpoint, label, key)

		cl, err := client.NewClientCli(endpoint)
		if err != nil {
			return err
		}

		_, err = cl.GetKeyValue(label, key)

		if !errors.Is(err, client.AppConfigClientError{Message: client.KVNotFoundError.Message, Info: key}) {
			return fmt.Errorf("expected %s, got %s", client.KVNotFoundError, err)
		}

		log.Printf("[INFO] KV is destroyed %s/%s/%s", endpoint, label, key)

		return nil
	}

	return nil
}

func testCheckFeatureDestroy(state *terraform.State) error {
	log.Printf("[INFO] Entering Destroy")
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "akc_feature" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		name := rs.Primary.Attributes["name"]
		label := rs.Primary.Attributes["label"]
		endpoint := rs.Primary.Attributes["endpoint"]
		log.Printf("[INFO] Checking that feature is destroyed %s/%s/%s", endpoint, label, name)

		cl, err := client.NewClientCli(endpoint)
		if err != nil {
			return err
		}

		// eventual consistency
		time.Sleep(5 * time.Second)

		_, err = cl.GetFeature(label, name)

		if !errors.Is(err, client.AppConfigClientError{Message: client.KVNotFoundError.Message, Info: fmt.Sprintf("%s%s", client.FeaturePrefix, name)}) {
			return fmt.Errorf("expected %s, got %s", client.KVNotFoundError, err)
		}

		log.Printf("[INFO] KV is destroyed %s/%s/%s", endpoint, label, name)

		return nil
	}

	return nil
}

func testCheckKeyValueExists(resource string, kv *client.KeyValueResponse) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("[INFO] Checking that KV exist (%s)", resource)

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		endpoint := rs.Primary.Attributes["endpoint"]
		key := rs.Primary.Attributes["key"]
		value := rs.Primary.Attributes["value"]
		label := rs.Primary.Attributes["label"]
		cl, err := client.NewClientCli(endpoint)
		if err != nil {
			return err
		}

		result, err := cl.GetKeyValue(label, key)
		if errors.Is(err, client.KVNotFoundError) {
			return fmt.Errorf("Cannot find resource %s", resource)
		}

		if err != nil {
			return fmt.Errorf("Error fetching KV with resource %s. %s", resource, err)
		}

		if result.Value != value {
			return fmt.Errorf("The given key %s exists but its value does not match", key)
		}

		*kv = result

		return nil
	}
}

func testCheckKeyValueSecretExists(resource string, kv *client.KeyValueResponse) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("[INFO] Checking that KV secret exist (%s)", resource)

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		endpoint := rs.Primary.Attributes["endpoint"]
		key := rs.Primary.Attributes["key"]
		label := rs.Primary.Attributes["label"]
		cl, err := client.NewClientCli(endpoint)
		if err != nil {
			return err
		}

		result, err := cl.GetKeyValue(label, key)
		if errors.Is(err, client.KVNotFoundError) {
			return fmt.Errorf("Cannot find resource %s", resource)
		}

		if err != nil {
			return fmt.Errorf("Error fetching KV with resource %s. %s", resource, err)
		}

		log.Printf("[INFO] Key %s found with value %s", key, result.Value)

		*kv = result

		return nil
	}
}

func testCheckFeatureExists(resource string, kv *client.FeatureResponse) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("[INFO] Checking that feature exist (%s)", resource)

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		endpoint := rs.Primary.Attributes["endpoint"]
		name := rs.Primary.Attributes["name"]
		label := rs.Primary.Attributes["label"]
		cl, err := client.NewClientCli(endpoint)
		if err != nil {
			return err
		}

		result, err := cl.GetFeature(label, name)
		if errors.Is(err, client.KVNotFoundError) {
			return fmt.Errorf("Cannot find resource %s", resource)
		}

		if err != nil {
			return fmt.Errorf("Error fetching feature with resource %s. %s", resource, err)
		}

		*kv = result

		return nil
	}
}

func testCheckStoredSecretID(kv *client.KeyValueResponse, expectedSecretID string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Print("[INFO] Checking stored secret ID")

		var actualStoredValue keyVaultReferenceValue
		err := json.Unmarshal([]byte(kv.Value), &actualStoredValue)
		if err != nil {
			return fmt.Errorf("Error while deserializing key-secret value")
		}

		if actualStoredValue.URI != expectedSecretID {
			return fmt.Errorf("Stored secret ID '%s' does not match expected one '%s'", actualStoredValue.URI, expectedSecretID)
		}

		log.Printf("[INFO] Right secret ID was stored: '%s'", expectedSecretID)

		return nil
	}
}

func testCheckStoredValue(kv *client.KeyValueResponse, expectedValue string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Print("[INFO] Checking stored value")

		if kv.Value != expectedValue {
			return fmt.Errorf("Stored value '%s' does not match expected one '%s'", kv.Value, expectedValue)
		}

		log.Printf("[INFO] Right value was stored: '%s'", expectedValue)

		return nil
	}
}

func randBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

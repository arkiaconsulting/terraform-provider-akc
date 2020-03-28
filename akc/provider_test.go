package akc

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// AzureProvider dsfdsfd
var AzureProvider *schema.Provider

// SupportedProviders dsfdsfd
var SupportedProviders map[string]terraform.ResourceProvider

var testProviders = map[string]terraform.ResourceProvider{
	"akc": Provider(),
}
var testAccProvider *schema.Provider

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// PreCheck ...
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
func BuildTerraformConfig() string {
	return fmt.Sprintf(`
resource "akc_key_value" "test" {
  endpoint     = "https://testlg.azconfig.io"
  key = "myKey"
  value = "myValue"
}
`)
}

// BuildTerraformConfigUpdateValue build terraform config
func BuildTerraformConfigUpdateValue() string {
	return fmt.Sprintf(`
resource "akc_key_value" "test" {
  endpoint     = "https://testlg.azconfig.io"
  key = "myKey"
  value = "myValueUpdated"
}
`)
}

func BuildTerraformConfigWithLabel() string {
	return fmt.Sprintf(`
resource "akc_key_value" "test" {
  endpoint     = "https://testlg.azconfig.io"
  key = "myKey"
  value = "myValue"
  label = "myLabel"
}
`)
}

// BuildTerraformConfigUpdateValue build terraform config
func BuildTerraformConfigUpdateValueWithLabel() string {
	return fmt.Sprintf(`
resource "akc_key_value" "test" {
  endpoint     = "https://testlg.azconfig.io"
  key = "myKey"
  value = "myValueUpdated"
  label = "myLabel"
}
`)
}

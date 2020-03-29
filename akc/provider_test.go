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

const endpointUnderTest = "https://testlg.azconfig.io"

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
func BuildTerraformConfigWithoutLabel(key string, value string) string {
	return fmt.Sprintf(`
resource "akc_key_value" "test" {
  endpoint     = "%s"
  key = "%s"
  value = "%s"
}
`, endpointUnderTest, key, value)
}

func BuildTerraformConfigWithLabel(label string, key string, value string) string {
	return fmt.Sprintf(`
resource "akc_key_value" "test" {
  endpoint     = "%s"
  label = "%s"
  key = "%s"
  value = "%s"
}
`, endpointUnderTest, label, key, value)
}

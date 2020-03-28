package akc

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCreateKeyValue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", "https://testlg.azconfig.io"),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", "myKey"),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", "myValue"),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", "%00"),
				),
			},
		},
	})
}

func TestAccCreateKeyValueWithLabel(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigWithLabel(),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", "https://testlg.azconfig.io"),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", "myKey"),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", "myValue"),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", "myLabel"),
				),
			},
		},
	})
}

func TestAccUpdateKeyValueUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", "https://testlg.azconfig.io"),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", "myKey"),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", "myValue"),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", "%00"),
				),
			},
			{
				Config: BuildTerraformConfigUpdateValue(),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", "https://testlg.azconfig.io"),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", "myKey"),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", "myValueUpdated"),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", "%00"),
				),
			},
		},
	})
}

func TestAccUpdateKeyValueUpdateWithLabel(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigWithLabel(),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", "https://testlg.azconfig.io"),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", "myKey"),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", "myValue"),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", "myLabel"),
				),
			},
			{
				Config: BuildTerraformConfigUpdateValueWithLabel(),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", "https://testlg.azconfig.io"),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", "myKey"),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", "myValueUpdated"),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", "myLabel"),
				),
			},
		},
	})
}

func testKeyValueDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		log.Printf("[INFO] Checking that KV is destroyed")

		if rs.Type != "akc_key_value" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		key := rs.Primary.Attributes["key"]
		label := rs.Primary.Attributes["label"]
		endpoint := rs.Primary.Attributes["endpoint"]
		client := client.NewAppConfigurationClient(endpoint)

		_, err := client.GetKeyValueWithLabel(key, label)

		notFoundErr := "Not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}

func testCheckKeyValueExists(resource string) resource.TestCheckFunc {
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
		client := client.NewAppConfigurationClient(endpoint)

		result, err := client.GetKeyValueWithLabel(key, label)
		if err != nil {
			return fmt.Errorf("Error fetching KV with resource %s. %s", resource, err)
		}

		if result.Value != value {
			return fmt.Errorf("The given key %s exists but its value does not match", key)
		}

		return nil
	}
}

package akc

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCreateKeyValue(t *testing.T) {
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigWithoutLabel(key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", "%00"),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
		},
	})
}

func TestAccCreateKeyValueWithLabel(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigWithLabel(label, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
		},
	})
}

func TestAccUpdateKeyValueValue(t *testing.T) {
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigWithoutLabel(key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", "%00"),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
			{
				Config: BuildTerraformConfigWithoutLabel(key, newValue),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", "%00"),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", newValue),
				),
			},
		},
	})
}

func TestAccUpdateKeyValueKey(t *testing.T) {
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	newKey := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigWithoutLabel(key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", "%00"),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
			{
				Config: BuildTerraformConfigWithoutLabel(newKey, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", "%00"),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", newKey),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
		},
	})
}

func TestAccUpdateKeyValueValueWithLabel(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigWithLabel(label, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
			{
				Config: BuildTerraformConfigWithLabel(label, key, newValue),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", newValue),
				),
			},
		},
	})
}

func TestAccUpdateKeyValueKeyWithLabel(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	newKey := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigWithLabel(label, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
			{
				Config: BuildTerraformConfigWithLabel(label, newKey, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", newKey),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
		},
	})
}

func TestAccUpdateKeyValueLabelWithLabel(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	newLabel := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigWithLabel(label, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
			{
				Config: BuildTerraformConfigWithLabel(newLabel, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", newLabel),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
		},
	})
}

func testKeyValueDestroy(state *terraform.State) error {
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

		cl := client.NewAppConfigurationClient(endpoint)

		_, err := cl.GetKeyValueWithLabel(key, label)

		if !errors.Is(err, client.KVNotFoundError) {
			return fmt.Errorf("expected %s, got %s", client.KVNotFoundError, err)
		}

		log.Printf("[INFO] KV is destroyed %s/%s/%s", endpoint, label, key)
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
		cl := client.NewAppConfigurationClient(endpoint)

		result, err := cl.GetKeyValueWithLabel(key, label)
		if errors.Is(err, client.KVNotFoundError) {
			return fmt.Errorf("Cannot find resource %s", resource)
		}

		if err != nil {
			return fmt.Errorf("Error fetching KV with resource %s. %s", resource, err)
		}

		if result.Value != value {
			return fmt.Errorf("The given key %s exists but its value does not match", key)
		}

		return nil
	}
}

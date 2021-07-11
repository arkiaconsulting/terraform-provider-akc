package akc

import (
	"testing"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKeyValue_createNoLabel(t *testing.T) {
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	var kv client.KeyValueResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithoutLabel(key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test", &kv),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", client.LabelNone),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
					testCheckStoredValue(&kv, value),
				),
			},
			{
				ResourceName:      "akc_key_value.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKeyValue_create(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	var kv client.KeyValueResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithLabel(label, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test", &kv),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
					testCheckStoredValue(&kv, value),
				),
			},
			{
				ResourceName:      "akc_key_value.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKeyValue_updateValue(t *testing.T) {
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	var kv client.KeyValueResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithoutLabel(key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test", &kv),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", client.LabelNone),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
					testCheckStoredValue(&kv, value),
				),
			},
			{
				Config: buildTerraformConfigWithoutLabel(key, newValue),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test", &kv),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", client.LabelNone),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", newValue),
					testCheckStoredValue(&kv, newValue),
				),
			},
		},
	})
}

func TestAccKeyValue_updateKey(t *testing.T) {
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	newKey := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	var kv client.KeyValueResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithoutLabel(key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test", &kv),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", client.LabelNone),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
					testCheckStoredValue(&kv, value),
				),
			},
			{
				Config: buildTerraformConfigWithoutLabel(newKey, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test", &kv),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", client.LabelNone),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", newKey),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
					testCheckStoredValue(&kv, value),
				),
			},
		},
	})
}

func TestAccKeyValue_updateValueWithLabel(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	newValue := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	var kv client.KeyValueResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithLabel(label, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test", &kv),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
					testCheckStoredValue(&kv, value),
				),
			},
			{
				Config: buildTerraformConfigWithLabel(label, key, newValue),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test", &kv),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", newValue),
					testCheckStoredValue(&kv, newValue),
				),
			},
		},
	})
}

func TestAccKeyValue_updateKeyWithLabel(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	newKey := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	var kv client.KeyValueResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithLabel(label, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test", &kv),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
					testCheckStoredValue(&kv, value),
				),
			},
			{
				Config: buildTerraformConfigWithLabel(label, newKey, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test", &kv),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", newKey),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
					testCheckStoredValue(&kv, value),
				),
			},
		},
	})
}

func TestAccKeyValue_updateLabelWithLabel(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	newLabel := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	var kv client.KeyValueResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithLabel(label, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test", &kv),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
					testCheckStoredValue(&kv, value),
				),
			},
			{
				Config: buildTerraformConfigWithLabel(newLabel, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test", &kv),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", newLabel),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
					testCheckStoredValue(&kv, value),
				),
			},
		},
	})
}

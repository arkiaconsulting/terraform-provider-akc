package akc

import (
	"testing"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCreateKeyValue(t *testing.T) {
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	value := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithoutLabel(key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", client.LabelNone),
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
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithLabel(label, key, value),
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
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithoutLabel(key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", client.LabelNone),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
			{
				Config: buildTerraformConfigWithoutLabel(key, newValue),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", client.LabelNone),
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
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithoutLabel(key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", client.LabelNone),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
			{
				Config: buildTerraformConfigWithoutLabel(newKey, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", client.LabelNone),
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
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithLabel(label, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
			{
				Config: buildTerraformConfigWithLabel(label, key, newValue),
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
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithLabel(label, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
			{
				Config: buildTerraformConfigWithLabel(label, newKey, value),
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
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigWithLabel(label, key, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueExists("akc_key_value.test"),
					resource.TestCheckResourceAttr("akc_key_value.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_value.test", "value", value),
				),
			},
			{
				Config: buildTerraformConfigWithLabel(newLabel, key, value),
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

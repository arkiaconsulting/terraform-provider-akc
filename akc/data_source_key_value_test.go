package akc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceKeyValue_basicNoLabel(t *testing.T) {
	key := "myKey"
	value := "myValuesdfsdfds"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigDataSourceKeyValue(key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.akc_key_value.test", "id"),
					resource.TestCheckResourceAttrSet("data.akc_key_value.test", "endpoint"),
					resource.TestCheckResourceAttrSet("data.akc_key_value.test", "label"),
					resource.TestCheckResourceAttr("data.akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("data.akc_key_value.test", "value", value),
				),
			},
		},
	})
}

func TestAccDataSourceKeyValue_basicLabel(t *testing.T) {
	key := "myKey2"
	value := "myValue2"
	label := "myLabel"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigDataSourceKeyValueLabel(label, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.akc_key_value.test", "id"),
					resource.TestCheckResourceAttrSet("data.akc_key_value.test", "endpoint"),
					resource.TestCheckResourceAttr("data.akc_key_value.test", "key", key),
					resource.TestCheckResourceAttr("data.akc_key_value.test", "label", label),
					resource.TestCheckResourceAttr("data.akc_key_value.test", "value", value),
				),
			},
		},
	})
}

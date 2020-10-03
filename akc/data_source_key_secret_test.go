package akc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceKeySecret_basicNoLabel(t *testing.T) {
	key := "my-secret1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigDataSourceKeySecret(key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.akc_key_secret.test", "id"),
					resource.TestCheckResourceAttrSet("data.akc_key_secret.test", "endpoint"),
					resource.TestCheckResourceAttr("data.akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttrSet("data.akc_key_secret.test", "secret_id"),
				),
			},
		},
	})
}

func TestAccDataSourceKeySecret_basicLabel(t *testing.T) {
	key := "my-secret12"
	label := "myLabel"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigDataSourceKeySecretLabel(label, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.akc_key_secret.test", "id"),
					resource.TestCheckResourceAttrSet("data.akc_key_secret.test", "endpoint"),
					resource.TestCheckResourceAttr("data.akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("data.akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttrSet("data.akc_key_secret.test", "secret_id"),
				),
			},
		},
	})
}

func TestAccDataSourceKeySecret_basicLabelWithLatestVersion(t *testing.T) {
	key := "my-secret12-latest"
	label := "myLabel"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigDataSourceKeySecretLabel(label, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.akc_key_secret.test", "id"),
					resource.TestCheckResourceAttrSet("data.akc_key_secret.test", "endpoint"),
					resource.TestCheckResourceAttr("data.akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("data.akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttrSet("data.akc_key_secret.test", "secret_id"),
				),
			},
		},
	})
}

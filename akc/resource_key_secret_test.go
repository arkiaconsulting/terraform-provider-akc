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

func TestAccCreateKeySecret(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	secretName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	secretID := fmt.Sprintf("https://toto/%s", secretName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigSecret(label, key, secretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test"),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
				),
			},
		},
	})
}

func TestAccUpdateKeySecretKey(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	newKey := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	secretName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	secretID := fmt.Sprintf("https://toto/%s", secretName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigSecret(label, key, secretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test"),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
				),
			},
			{
				Config: BuildTerraformConfigSecret(label, newKey, secretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test"),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", newKey),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
				),
			},
		},
	})
}

func TestAccUpdateKeySecretSecretID(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	secretName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	secretID := fmt.Sprintf("https://toto/%s", secretName)
	newSecretID := fmt.Sprintf("%s%s", secretID, "new")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigSecret(label, key, secretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test"),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
				),
			},
			{
				Config: BuildTerraformConfigSecret(label, key, newSecretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test"),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", newSecretID),
				),
			},
		},
	})
}

func TestAccUpdateKeySecretSecretLabel(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	newLabel := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	secretName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	secretID := fmt.Sprintf("https://toto/%s", secretName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: BuildTerraformConfigSecret(label, key, secretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test"),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
				),
			},
			{
				Config: BuildTerraformConfigSecret(newLabel, key, secretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test"),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", newLabel),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
				),
			},
		},
	})
}

func testCheckKeyValueSecretExists(resource string) resource.TestCheckFunc {
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
		secretID := rs.Primary.Attributes["secret_id"]
		label := rs.Primary.Attributes["label"]
		cl, err := client.NewAppConfigurationClient(endpoint)
		if err != nil {
			panic(err)
		}

		result, err := cl.GetKeyValue(label, key)
		if errors.Is(err, client.KVNotFoundError) {
			return fmt.Errorf("Cannot find resource %s", resource)
		}

		if err != nil {
			return fmt.Errorf("Error fetching KV with resource %s. %s", resource, err)
		}

		expectedValue := fmt.Sprintf("{\"uri\":\"%s\"}", secretID)
		if result.Value != expectedValue {
			return fmt.Errorf("The given key %s exists but its value does not match", key)
		}

		return nil
	}
}

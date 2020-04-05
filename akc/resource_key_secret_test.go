package akc

import (
	"fmt"
	"testing"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCreateKeySecret(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	secretName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	secretID := fmt.Sprintf("https://toto/%s/version", secretName)
	var kv client.KeyValueResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigSecret(label, key, secretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "false"),
					testCheckStoredSecretID(&kv, secretID),
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
	secretID := fmt.Sprintf("https://toto/%s/version", secretName)
	var kv client.KeyValueResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigSecret(label, key, secretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "false"),
					testCheckStoredSecretID(&kv, secretID),
				),
			},
			{
				Config: buildTerraformConfigSecret(label, newKey, secretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", newKey),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "false"),
					testCheckStoredSecretID(&kv, secretID),
				),
			},
		},
	})
}

func TestAccUpdateKeySecretSecretID(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	secretName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	secretID := fmt.Sprintf("https://toto/%s/version", secretName)
	newSecretID := fmt.Sprintf("%s%s", secretID, "new")
	var kv client.KeyValueResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigSecret(label, key, secretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "false"),
					testCheckStoredSecretID(&kv, secretID),
				),
			},
			{
				Config: buildTerraformConfigSecret(label, key, newSecretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", newSecretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "false"),
					testCheckStoredSecretID(&kv, newSecretID),
				),
			},
		},
	})
}

func TestAccUpdateKeySecretLabel(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	newLabel := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	secretName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	secretID := fmt.Sprintf("https://toto/%s/version", secretName)
	var kv client.KeyValueResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigSecret(label, key, secretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "false"),
					testCheckStoredSecretID(&kv, secretID),
				),
			},
			{
				Config: buildTerraformConfigSecret(newLabel, key, secretID),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", newLabel),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "false"),
					testCheckStoredSecretID(&kv, secretID),
				),
			},
		},
	})
}

func TestAccCreateKeySecretWithVersionLatestVersionTrue(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	secretName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	secretID := fmt.Sprintf("https://toto/secrets/%s/version", secretName)
	var kv client.KeyValueResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigSecretLatestVersion(label, key, secretID, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "true"),
					testCheckStoredSecretID(&kv, fmt.Sprintf("https://toto/secrets/%s", secretName)),
				),
			},
		},
	})
}

func TestAccCreateKeySecretWithoutVersionLatestVersionTrue(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	secretName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	secretID := fmt.Sprintf("https://toto/secrets/%s", secretName)
	var kv client.KeyValueResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigSecretLatestVersion(label, key, secretID, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "true"),
					testCheckStoredSecretID(&kv, secretID),
				),
			},
		},
	})
}

func TestAccUpdateLatestVersionKeySecret(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	secretName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	secretID := fmt.Sprintf("https://toto/secrets/%s/version", secretName)
	var kv client.KeyValueResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigSecretLatestVersion(label, key, secretID, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "false"),
					testCheckStoredSecretID(&kv, secretID),
				),
			},
			{
				Config: buildTerraformConfigSecretLatestVersion(label, key, secretID, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "true"),
					testCheckStoredSecretID(&kv, fmt.Sprintf("https://toto/secrets/%s", secretName)),
				),
			},
		},
	})
}

func TestAccUpdateSecretIdKeySecretWithVersion(t *testing.T) {
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	key := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	secretName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	secretID := fmt.Sprintf("https://toto/secrets/%s/version", secretName)
	newSecretID := fmt.Sprintf("https://toto/secrets/%s/otherversion", secretName)
	var kv client.KeyValueResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckKeyValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildTerraformConfigSecretLatestVersion(label, key, secretID, false),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", secretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "false"),
					testCheckStoredSecretID(&kv, secretID),
				),
			},
			{
				Config: buildTerraformConfigSecretLatestVersion(label, key, newSecretID, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckKeyValueSecretExists("akc_key_secret.test", &kv),
					resource.TestCheckResourceAttr("akc_key_secret.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_key_secret.test", "label", label),
					resource.TestCheckResourceAttr("akc_key_secret.test", "key", key),
					resource.TestCheckResourceAttr("akc_key_secret.test", "secret_id", newSecretID),
					resource.TestCheckResourceAttr("akc_key_secret.test", "latest_version", "true"),
					testCheckStoredSecretID(&kv, fmt.Sprintf("https://toto/secrets/%s", secretName)),
				),
			},
		},
	})
}

package akc

import (
	"strconv"
	"testing"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFeature_create(t *testing.T) {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	label := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	enabled := randBool()

	var kv client.FeatureResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckFeatureDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildLabeledFeature(name, label, description, enabled),
				Check: resource.ComposeTestCheckFunc(
					testCheckFeatureExists("akc_feature.test", &kv),
					resource.TestCheckResourceAttr("akc_feature.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_feature.test", "label", label),
					resource.TestCheckResourceAttr("akc_feature.test", "name", name),
					resource.TestCheckResourceAttr("akc_feature.test", "enabled", strconv.FormatBool(enabled)),
					resource.TestCheckResourceAttr("akc_feature.test", "description", description),
				),
			},
			{
				ResourceName:      "akc_feature.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceFeatureNoLabel_createNoLabel(t *testing.T) {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	enabled := randBool()

	var kv client.FeatureResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckFeatureDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildFeature(name, description, enabled),
				Check: resource.ComposeTestCheckFunc(
					testCheckFeatureExists("akc_feature.test", &kv),
					resource.TestCheckResourceAttr("akc_feature.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_feature.test", "name", name),
					resource.TestCheckResourceAttr("akc_feature.test", "label", client.LabelNone),
					resource.TestCheckResourceAttr("akc_feature.test", "enabled", strconv.FormatBool(enabled)),
					resource.TestCheckResourceAttr("akc_feature.test", "description", description),
				),
			},
			{
				ResourceName:      "akc_feature.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceFeature_updateNoLabel(t *testing.T) {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	enabled := randBool()
	newEnabled := !enabled

	var kv client.FeatureResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckFeatureDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildFeature(name, description, enabled),
				Check: resource.ComposeTestCheckFunc(
					testCheckFeatureExists("akc_feature.test", &kv),
					resource.TestCheckResourceAttr("akc_feature.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_feature.test", "name", name),
					resource.TestCheckResourceAttr("akc_feature.test", "label", client.LabelNone),
					resource.TestCheckResourceAttr("akc_feature.test", "enabled", strconv.FormatBool(enabled)),
					resource.TestCheckResourceAttr("akc_feature.test", "description", description),
				),
			},
			{
				Config: buildFeature(name, newDescription, newEnabled),
				Check: resource.ComposeTestCheckFunc(
					testCheckFeatureExists("akc_feature.test", &kv),
					resource.TestCheckResourceAttr("akc_feature.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_feature.test", "name", name),
					resource.TestCheckResourceAttr("akc_feature.test", "label", client.LabelNone),
					resource.TestCheckResourceAttr("akc_feature.test", "enabled", strconv.FormatBool(newEnabled)),
					resource.TestCheckResourceAttr("akc_feature.test", "description", newDescription),
				),
			},
		},
	})
}

func TestAccResourceFeature_updateLabel(t *testing.T) {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	label := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	newDescription := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	enabled := randBool()
	newEnabled := !enabled

	var kv client.FeatureResponse

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { preCheck(t) },
		Providers:    testProviders,
		CheckDestroy: testCheckFeatureDestroy,
		Steps: []resource.TestStep{
			{
				Config: buildLabeledFeature(name, label, description, enabled),
				Check: resource.ComposeTestCheckFunc(
					testCheckFeatureExists("akc_feature.test", &kv),
					resource.TestCheckResourceAttr("akc_feature.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_feature.test", "name", name),
					resource.TestCheckResourceAttr("akc_feature.test", "label", label),
					resource.TestCheckResourceAttr("akc_feature.test", "enabled", strconv.FormatBool(enabled)),
					resource.TestCheckResourceAttr("akc_feature.test", "description", description),
				),
			},
			{
				Config: buildLabeledFeature(name, label, newDescription, newEnabled),
				Check: resource.ComposeTestCheckFunc(
					testCheckFeatureExists("akc_feature.test", &kv),
					resource.TestCheckResourceAttr("akc_feature.test", "endpoint", endpointUnderTest),
					resource.TestCheckResourceAttr("akc_feature.test", "name", name),
					resource.TestCheckResourceAttr("akc_feature.test", "label", label),
					resource.TestCheckResourceAttr("akc_feature.test", "enabled", strconv.FormatBool(newEnabled)),
					resource.TestCheckResourceAttr("akc_feature.test", "description", newDescription),
				),
			},
		},
	})
}

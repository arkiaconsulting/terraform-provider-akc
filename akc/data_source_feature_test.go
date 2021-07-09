package akc

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFeature_noLabel(t *testing.T) {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	enabled := randBool()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: buildTemplateNoLabel(name, description, enabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.akc_feature.test", "id"),
					resource.TestCheckResourceAttrSet("data.akc_feature.test", "endpoint"),
					resource.TestCheckResourceAttr("data.akc_feature.test", "name", name),
					resource.TestCheckResourceAttr("data.akc_feature.test", "label", client.LabelNone),
					resource.TestCheckResourceAttr("data.akc_feature.test", "description", description),
					resource.TestCheckResourceAttr("data.akc_feature.test", "enabled", strconv.FormatBool(enabled)),
				),
			},
		},
	})
}

func TestAccDataSourceFeature_label(t *testing.T) {
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	label := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	enabled := randBool()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preCheck(t) },
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: buildLabeledFeature(name, label, description, enabled),
			},
			{
				Config: buildTemplateLabel(name, label, description, enabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.akc_feature.test", "id"),
					resource.TestCheckResourceAttrSet("data.akc_feature.test", "endpoint"),
					resource.TestCheckResourceAttr("data.akc_feature.test", "name", name),
					resource.TestCheckResourceAttr("data.akc_feature.test", "label", label),
					resource.TestCheckResourceAttr("data.akc_feature.test", "description", description),
					resource.TestCheckResourceAttr("data.akc_feature.test", "enabled", strconv.FormatBool(enabled)),
				),
			},
		},
	})
}

func buildTemplateNoLabel(name string, description string, enabled bool) string {
	return fmt.Sprintf(`
%s

data "akc_feature" "test" {
	endpoint     = akc_feature.test.endpoint
	name = akc_feature.test.name
  }
`, buildFeature(name, description, enabled))
}

func buildTemplateLabel(name string, label string, description string, enabled bool) string {
	return fmt.Sprintf(`
%s

data "akc_feature" "test" {
	endpoint     = akc_feature.test.endpoint
	name = akc_feature.test.name
	label = akc_feature.test.label
  }
`, buildLabeledFeature(name, label, description, enabled))
}

package akc

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider coucouc
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{
			"akc_key_value":  resourceKeyValue(),
			"akc_key_secret": resourceKeySecret(),
		},
	}
}

// ProviderConfigure is only here for starting with a good template
func ProviderConfigure(d *schema.ResourceData) error {
	log.Printf("[DEBUG] configure")

	return nil
}

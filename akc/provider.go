package akc

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider akc
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{
			"akc_key_value":  resourceKeyValue(),
			"akc_key_secret": resourceKeySecret(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akc_key_value":  dataSourceKeyValue(),
			"akc_key_secret": dataSourceKeySecret(),
		},
	}
}

// ProviderConfigure is only here for starting with a good template
func ProviderConfigure(d *schema.ResourceData) error {
	log.Printf("[DEBUG] configure")

	return nil
}

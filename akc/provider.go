package akc

import (
	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider akc
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Description: "App Configuration url to target",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKC_ENDPOINT", nil),
			},
			"client_id": {
				Type:        schema.TypeString,
				Description: "Azure AD Client Id",
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Description: "Azure AD Client Secret",
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", nil),
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Description: "Azure AD Tenant ID",
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", nil),
			},
			"msi": {
				Type:        schema.TypeBool,
				Description: "Use msi if available, will fail if not in a MSI context",
				Optional:    true,
				Default:     false,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"akc_key_value":  resourceKeyValue(),
			"akc_key_secret": resourceKeySecret(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akc_key_value":  dataSourceKeyValue(),
			"akc_key_secret": dataSourceKeySecret(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	endpoint := d.Get("endpoint").(string)

	if d.Get("msi").(bool) {
		return client.NewClientMsi(endpoint)
	}

	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	tenantId := d.Get("tenant_id").(string)

	if (clientId != "") && (clientSecret != "") && (tenantId != "") {
		return client.NewClientCreds(endpoint, clientId, clientSecret, tenantId)
	}

	return client.NewClientCli(endpoint)
}

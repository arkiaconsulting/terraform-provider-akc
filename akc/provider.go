package akc

import (
	"terraform-provider-akc/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider akc
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZAPPCONF_ENDPOINT", nil),
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", nil),
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", nil),
			},
			"auth_method": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AZAPPCONF_AUTH_METHOD", nil),
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
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	tenantId := d.Get("tenant_id").(string)
	authMethod := d.Get("auth_method").(string)

	if (clientId != "") && (clientSecret != "") {
		authMethod = "creds"
	}

	var c *client.Client
	var err error

	switch authMethod {
	case "creds":
		c, err = client.NewClientCreds(endpoint, clientId, clientSecret, tenantId)
	case "cli":
		c, err = client.NewClientCli(endpoint)
	case "msi":
		c, err = client.NewClientMsi(endpoint)
	case "env":
	default:
		c, err = client.NewClientEnv(endpoint)
	}
	return c, err
}

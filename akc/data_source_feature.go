package akc

import (
	"fmt"
	"log"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFeature() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFeatureRead,

		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  client.LabelNone,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(readTimeout),
		},
	}
}

func dataSourceFeatureRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading resource %s", d.Id())

	endpoint := d.Get("endpoint").(string)
	label := d.Get("label").(string)
	name := d.Get("name").(string)

	cl, err := getClient(endpoint, meta.(func(endpoint string) (*client.Client, error)))
	if err != nil {
		return fmt.Errorf("error building client for endpoint %s: %+v", endpoint, err)
	}

	var feature client.FeatureResponse
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		feature, err = cl.GetFeature(label, name)

		if err != nil {
			if client.IsNotFound(err) {
				log.Printf("[INFO] retrying to get feature '%s:%s' because: %s", label, name, err)

				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error getting App Configuration feature %s/%s: %+v", label, name, err)
	}

	id, err := formatFeatureID(endpoint, label, name)
	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("label", label)
	d.Set("name", name)
	d.Set("description", feature.Description)
	d.Set("enabled", feature.Enabled)

	log.Printf("[INFO] KV has been fetched %s/%s/%s=%s", endpoint, label, name, feature.Description)

	return nil
}

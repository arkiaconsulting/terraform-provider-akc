package akc

import (
	"fmt"
	"log"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKeyValue() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeyValueRead,

		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  client.LabelNone,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKeyValueRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading resource %s", d.Id())

	endpoint := d.Get("endpoint").(string)
	label := d.Get("label").(string)
	key := d.Get("key").(string)

	cl := meta.(*client.Client)
	cl.Endpoint = endpoint

	log.Printf("[INFO] Fetching KV %s/%s/%s", endpoint, label, key)
	kv, err := cl.GetKeyValue(label, key)
	if err != nil {
		return fmt.Errorf("error getting App Configuration key %s/%s: %+v", label, key, err)
	}

	if kv.Label == "" {
		kv.Label = client.LabelNone
	}

	id, err := formatID(endpoint, label, key)
	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("key", key)
	d.Set("value", kv.Value)
	d.Set("label", kv.Label)

	log.Printf("[INFO] KV has been fetched %s/%s/%s=%s", endpoint, label, key, kv.Value)

	return nil
}

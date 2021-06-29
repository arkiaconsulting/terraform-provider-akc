package akc

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKeySecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeySecretRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  client.LabelNone,
			},
			"secret_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKeySecretRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading resource %s", d.Id())

	label := d.Get("label").(string)
	key := d.Get("key").(string)

	c := meta.(*client.Client)
	endpoint := c.Endpoint

	log.Printf("[INFO] Fetching KV %s/%s/%s", endpoint, label, key)
	kv, err := c.GetKeyValue(label, key)
	if err != nil {
		return fmt.Errorf("Error getting App Configuration key %s/%s: %+v", label, key, err)
	}

	if kv.Label == "" {
		kv.Label = client.LabelNone
	}

	id, err := formatID(endpoint, label, key)
	if err != nil {
		return err
	}

	var wrapper keyVaultReferenceValue
	err = json.Unmarshal([]byte(kv.Value), &wrapper)

	d.SetId(id)
	d.Set("key", key)
	d.Set("secret_id", wrapper.URI)
	d.Set("label", kv.Label)

	log.Printf("[INFO] KV has been fetched %s/%s/%s=%s", endpoint, label, key, wrapper.URI)

	return nil
}

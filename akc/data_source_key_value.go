package akc

import (
	"context"
	"fmt"
	"log"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKeyValue() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeyValueRead,

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
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKeyValueRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading resource %s", d.Id())
	var diags diag.Diagnostics

	label := d.Get("label").(string)
	key := d.Get("key").(string)

	cl := meta.(*client.Client)
	endpoint := cl.Endpoint

	log.Printf("[INFO] Fetching KV %s/%s/%s", endpoint, label, key)
	kv, err := cl.GetKeyValue(label, key)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting App Configuration key %s/%s: %+v", label, key, err))
	}

	if kv.Label == "" {
		kv.Label = client.LabelNone
	}

	id, err := formatID(endpoint, label, key)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	d.Set("key", key)
	d.Set("value", kv.Value)
	d.Set("label", kv.Label)

	log.Printf("[INFO] KV has been fetched %s/%s/%s=%s", endpoint, label, key, kv.Value)

	return diags
}

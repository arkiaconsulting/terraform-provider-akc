package akc

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeyValue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeyValueCreate,
		ReadContext:   resourceKeyValueRead,
		UpdateContext: resourceKeyValueUpdate,
		DeleteContext: resourceKeyValueDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  client.LabelNone,
				ForceNew: true,
			},
		},
	}
}

func resourceKeyValueCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Print("[INFO] Creating resource")

	c := m.(*client.Client)
	endpoint := c.Endpoint

	key := d.Get("key").(string)
	value := d.Get("value").(string)
	label := d.Get("label").(string)

	if d.IsNewResource() {
		_, err := c.GetKeyValue(label, key)
		if err == nil {
			return diag.FromErr(fmt.Errorf("The resource needs to be imported: %s", "akc_key_value"))
		}
	}

	_, err := c.SetKeyValue(label, key, value)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := formatID(endpoint, label, key)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return resourceKeyValueRead(ctx, d, m)
}

func resourceKeyValueRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading resource %s", d.Id())
	var diags diag.Diagnostics

	_, label, key := parseID(d.Id())
	c := m.(*client.Client)
	endpoint := c.Endpoint

	log.Printf("[INFO] Fetching KV %s/%s/%s", endpoint, label, key)
	result, err := c.GetKeyValue(label, key)
	if err != nil {
		log.Printf("[INFO] KV not found, removing from state: %s/%s/%s", endpoint, label, key)
		d.SetId("")
		return nil
	}

	if result.Label == "" {
		result.Label = client.LabelNone
	}

	d.Set("key", key)
	d.Set("value", result.Value)
	d.Set("label", result.Label)

	log.Printf("[INFO] KV has been fetched %s/%s/%s=%s", endpoint, label, key, result.Value)

	return diags
}

func resourceKeyValueUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating resource %s", d.Id())

	_, label, key := parseID(d.Id())
	c := m.(*client.Client)
	endpoint := c.Endpoint
	value := d.Get("value").(string)

	_, err := c.SetKeyValue(label, key, value)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := formatID(endpoint, label, key)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return resourceKeyValueRead(ctx, d, m)
}

func resourceKeyValueDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting resource %s", d.Id())
	var diags diag.Diagnostics

	_, label, key := parseID(d.Id())

	c := m.(*client.Client)

	_, err := c.DeleteKeyValue(label, key)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func formatID(endpoint string, label string, key string) (string, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("Unable to parse the given endpoint %s", endpoint)
	}

	name := url.Host

	return fmt.Sprintf("%s/%s/%s", name, label, key), nil
}

func parseID(id string) (endpoint string, label string, key string) {
	split := strings.Split(id, "/")

	endpoint = fmt.Sprintf("https://%s", split[0])
	label = split[1]
	key = split[2]

	return endpoint, label, key
}

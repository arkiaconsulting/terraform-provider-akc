package akc

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeySecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeySecretCreate,
		ReadContext:   resourceKeySecretRead,
		UpdateContext: resourceKeySecretUpdate,
		DeleteContext: resourceKeyValueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  client.LabelNone,
				ForceNew: true,
			},
			"latest_version": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKeySecretCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	endpoint := c.Endpoint

	key := d.Get("key").(string)
	value := d.Get("secret_id").(string)
	label := d.Get("label").(string)
	trim := d.Get("latest_version").(bool)

	if trim {
		value = trimVersion(value)
	}

	_, err = c.SetKeyValueSecret(key, value, label)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := formatID(endpoint, label, key)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	d.Set("value", value)

	return resourceKeySecretRead(ctx, d, meta)
}

func resourceKeySecretUpdate(d *schema.ResourceData, m interface{}) error {
	_, label, key := parseID(d.Id())
	c := m.(*client.Client)

	endpoint := c.Endpoint
	value := d.Get("secret_id").(string)
	trim := d.Get("latest_version").(bool)

	if trim {
		value = trimVersion(value)
	}

	_, err := c.SetKeyValueSecret(key, value, label)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := formatID(endpoint, label, key)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	d.Set("value", value)

	return resourceKeySecretRead(ctx, d, m)
}

func resourceKeySecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading resource %s", d.Id())
	var diags diag.Diagnostics

	_, label, key := parseID(d.Id())
	c := m.(*client.Client)
	endpoint := c.Endpoint

	log.Printf("[INFO] Fetching KV %s/%s/%s", endpoint, label, key)
	kv, err := c.GetKeyValue(label, key)
	if err != nil {
		log.Printf("[INFO] KV not found, removing from state: %s/%s/%s", endpoint, label, key)
		d.SetId("")
		return nil
	}

	if kv.Label == "" {
		kv.Label = client.LabelNone
	}

	var wrapper keyVaultReferenceValue
	err = json.Unmarshal([]byte(kv.Value), &wrapper)

	d.Set("key", key)
	d.Set("value", wrapper.URI)
	d.Set("label", kv.Label)

	log.Printf("[INFO] KV has been fetched %s/%s/%s=%s", endpoint, label, key, wrapper.URI)

	return diags
}

func trimVersion(s string) string {
	url, _ := url.Parse(s)

	pathParts := strings.Split(url.Path, "/")

	if len(pathParts) == 4 {
		url.Path = strings.Join(pathParts[:len(pathParts)-1], "/")
	}

	return url.String()
}

type keyVaultReferenceValue struct {
	URI string
}

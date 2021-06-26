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
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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

func resourceKeySecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	endpoint := d.Get("endpoint").(string)
	key := d.Get("key").(string)
	value := d.Get("secret_id").(string)
	label := d.Get("label").(string)
	trim := d.Get("latest_version").(bool)

	cl, err := client.BuildAppConfigurationClient(ctx, endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if trim {
		value = trimVersion(value)
	}

	_, err = cl.SetKeyValueSecret(key, value, label)
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

func resourceKeySecretUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	endpoint, label, key := parseID(d.Id())
	value := d.Get("secret_id").(string)
	trim := d.Get("latest_version").(bool)

	cl, err := client.BuildAppConfigurationClient(ctx, endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	if trim {
		value = trimVersion(value)
	}

	_, err = cl.SetKeyValueSecret(key, value, label)
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

	endpoint, label, key := parseID(d.Id())

	cl, err := client.BuildAppConfigurationClient(ctx, endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Fetching KV %s/%s/%s", endpoint, label, key)
	kv, err := cl.GetKeyValue(label, key)
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

	d.Set("endpoint", endpoint)
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

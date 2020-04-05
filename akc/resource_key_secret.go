package akc

import (
	"net/url"
	"strings"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKeySecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeySecretCreate,
		Read:   resourceKeyValueRead,
		Update: resourceKeySecretUpdate,
		Delete: resourceKeyValueDelete,
		Exists: resourceKeyValueExists,

		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  client.LabelNone,
				ForceNew: true,
			},
			"latest_version": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceKeySecretCreate(d *schema.ResourceData, m interface{}) error {
	endpoint := d.Get("endpoint").(string)
	client, err := client.NewAppConfigurationClient(endpoint)
	key := d.Get("key").(string)
	value := d.Get("secret_id").(string)
	label := d.Get("label").(string)
	trim := d.Get("latest_version").(bool)

	if trim {
		value = trimVersion(value)
	}

	_, err = client.SetKeyValueSecret(key, value, label)
	if err != nil {
		return err
	}

	id, err := formatID(endpoint, label, key)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceKeyValueRead(d, m)
}

func resourceKeySecretUpdate(d *schema.ResourceData, m interface{}) error {
	endpoint, label, key := parseID(d.Id())
	client, err := client.NewAppConfigurationClient(endpoint)
	value := d.Get("secret_id").(string)
	trim := d.Get("latest_version").(bool)

	if trim {
		value = trimVersion(value)
	}

	_, err = client.SetKeyValueSecret(key, value, label)
	if err != nil {
		return err
	}

	id, err := formatID(endpoint, label, key)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceKeyValueRead(d, m)
}

func trimVersion(s string) string {
	url, _ := url.Parse(s)

	pathParts := strings.Split(url.Path, "/")

	if len(pathParts) == 4 {
		url.Path = strings.Join(pathParts[:len(pathParts)-1], "/")
	}

	return url.String()
}

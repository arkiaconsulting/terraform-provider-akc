package akc

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKeyValue() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyValueCreate,
		Read:   resourceKeyValueRead,
		Update: resourceKeyValueUpdate,
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
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  client.LabelNone,
				ForceNew: true,
			},
		},
	}
}

func resourceKeyValueCreate(d *schema.ResourceData, m interface{}) error {
	endpoint := d.Get("endpoint").(string)
	cl, err := client.NewAppConfigurationClient(endpoint)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	label := d.Get("label").(string)

	_, err = cl.SetKeyValue(label, key, value)
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

func resourceKeyValueRead(d *schema.ResourceData, m interface{}) error {
	log.Print("[INFO] Reading...")

	endpoint, label, key := parseID(d.Id())
	cl, err := client.NewAppConfigurationClient(endpoint)

	result, err := cl.GetKeyValue(label, key)
	if err != nil {
		return err
	}

	id, err := formatID(endpoint, label, key)
	if err != nil {
		return err
	}

	if result.Label == "" {
		result.Label = client.LabelNone
	}

	d.SetId(id)
	d.Set("endpoint", endpoint)
	d.Set("key", key)
	d.Set("value", result.Value)
	d.Set("label", result.Label)

	return nil
}

func resourceKeyValueUpdate(d *schema.ResourceData, m interface{}) error {
	endpoint, label, key := parseID(d.Id())
	cl, err := client.NewAppConfigurationClient(endpoint)
	value := d.Get("value").(string)

	_, err = cl.SetKeyValue(label, key, value)
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

func resourceKeyValueDelete(d *schema.ResourceData, m interface{}) error {
	endpoint, label, key := parseID(d.Id())

	cl, err := client.NewAppConfigurationClient(endpoint)

	_, err = cl.DeleteKeyValue(label, key)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceKeyValueExists(d *schema.ResourceData, m interface{}) (bool, error) {
	endpoint, label, key := parseID(d.Id())
	cl, err := client.NewAppConfigurationClient(endpoint)

	_, err = cl.GetKeyValue(label, key)
	if err != nil {
		return false, err
	}

	return true, nil
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
	log.Printf("[INFO] Parsing id %s", id)
	split := strings.Split(id, "/")

	endpoint = fmt.Sprintf("https://%s", split[0])
	label = split[1]
	key = split[2]

	log.Printf("[INFO] Parsed %s with label %s and key: %s", endpoint, label, key)

	return endpoint, label, key
}

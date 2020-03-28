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
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
		},
	}
}

func resourceKeyValueCreate(d *schema.ResourceData, m interface{}) error {
	endpoint := d.Get("endpoint").(string)
	client := client.NewAppConfigurationClient(endpoint)
	key := d.Get("key").(string)
	value := d.Get("value").(string)

	_, err := client.SetKeyValue(key, value)
	if err != nil {
		return err
	}

	id, err := formatID(endpoint, key)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceKeyValueRead(d, m)
}

func resourceKeyValueRead(d *schema.ResourceData, m interface{}) error {
	log.Print("[INFO] Reading...")

	endpoint, key := parseID(d.Id())
	client := client.NewAppConfigurationClient(endpoint)

	result, err := client.GetKeyValue(key)
	if err != nil {
		return err
	}

	id, err := formatID(endpoint, key)
	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("endpoint", endpoint)
	d.Set("key", key)
	d.Set("value", result.Value)

	return nil
}

func resourceKeyValueUpdate(d *schema.ResourceData, m interface{}) error {
	endpoint := d.Get("endpoint").(string)
	client := client.NewAppConfigurationClient(endpoint)
	key := d.Get("key").(string)
	value := d.Get("value").(string)

	_, err := client.SetKeyValue(key, value)
	if err != nil {
		return err
	}

	id, err := formatID(endpoint, key)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceKeyValueRead(d, m)
}

func resourceKeyValueDelete(d *schema.ResourceData, m interface{}) error {
	log.Print("[INFO] Destroying...")

	endpoint := d.Get("endpoint").(string)
	client := client.NewAppConfigurationClient(endpoint)
	key := d.Get("key").(string)

	isDeleted, err := client.DeleteKeyValue(key)
	if err != nil {
		return err
	}

	if !isDeleted {
		log.Printf("[WARNING] KV was not deleted")
	}

	return nil
}

func resourceKeyValueExists(d *schema.ResourceData, m interface{}) (bool, error) {
	endpoint, key := parseID(d.Id())
	client := client.NewAppConfigurationClient(endpoint)

	_, err := client.GetKeyValue(key)
	if err != nil {
		return false, err
	}

	return true, nil
}

func formatID(endpoint string, key string) (string, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("Unable to parse the given endpoint %s", endpoint)
	}

	name := url.Host

	return fmt.Sprintf("%s/%s", name, key), nil
}

func parseID(id string) (endpoint string, key string) {
	log.Printf("[INFO] Parsing id %s", id)
	split := strings.Split(id, "/")

	endpoint = fmt.Sprintf("https://%s", split[0])
	key = split[1]

	log.Printf("[INFO] Parsed %s with key: %s", endpoint, key)

	return endpoint, key
}

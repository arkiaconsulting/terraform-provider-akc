package akc

import (
	"errors"
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
	log.Print("[INFO] Creating resource")

	endpoint := d.Get("endpoint").(string)
	cl, err := client.NewAppConfigurationClient(endpoint)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	label := d.Get("label").(string)

	if d.IsNewResource() {
		_, err := cl.GetKeyValue(label, key)
		if err == nil {
			return fmt.Errorf("The resource needs to be imported: %s", "akc_key_value")
		}
	}

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
	log.Printf("[INFO] Reading resource %s", d.Id())

	endpoint, label, key := parseID(d.Id())
	cl, err := client.NewAppConfigurationClient(endpoint)

	log.Printf("[INFO] Fetching KV %s/%s/%s", endpoint, label, key)
	result, err := cl.GetKeyValue(label, key)
	if err != nil {
		log.Printf("[INFO] KV not found, removing from state: %s/%s/%s", endpoint, label, key)
		d.SetId("")
		return nil
	}

	if result.Label == "" {
		result.Label = client.LabelNone
	}

	d.Set("endpoint", endpoint)
	d.Set("key", key)
	d.Set("value", result.Value)
	d.Set("label", result.Label)

	log.Printf("[INFO] KV has been fetched %s/%s/%s=%s", endpoint, label, key, result.Value)

	return nil
}

func resourceKeyValueUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Updating resource %s", d.Id())

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
	log.Printf("[INFO] Deleting resource %s", d.Id())

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
	log.Printf("[INFO] Does resource exist %s", d.Id())

	endpoint, label, key := parseID(d.Id())
	cl, err := client.NewAppConfigurationClient(endpoint)

	_, err = cl.GetKeyValue(label, key)
	if err != nil {
		if errors.Is(err, client.AppConfigClientError{Message: client.KVNotFoundError.Message, Info: key}) {
			return false, nil
		} else {
			return false, err
		}
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
	split := strings.Split(id, "/")

	endpoint = fmt.Sprintf("https://%s", split[0])
	label = split[1]
	key = split[2]

	return endpoint, label, key
}

package akc

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeyValue() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyValueCreate,
		Read:   resourceKeyValueRead,
		Update: resourceKeyValueUpdate,
		Delete: resourceKeyValueDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
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

func resourceKeyValueCreate(d *schema.ResourceData, meta interface{}) error {
	log.Print("[INFO] Creating resource")

	endpoint := d.Get("endpoint").(string)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	label := d.Get("label").(string)

	cl, err := getOrReuseClient(endpoint, meta.(func(endpoint string) (*client.Client, error)))
	if err != nil {
		return fmt.Errorf("error building client for endpoint %s: %+v", endpoint, err)
	}

	if d.IsNewResource() {
		_, err := cl.GetKeyValue(label, key)
		if err == nil {
			return fmt.Errorf("the resource needs to be imported: %s", "akc_key_value")
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

	return resourceKeyValueRead(d, meta)
}

func resourceKeyValueRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading resource %s", d.Id())

	endpoint, label, key := parseID(d.Id())

	cl, err := getOrReuseClient(endpoint, meta.(func(endpoint string) (*client.Client, error)))
	if err != nil {
		return fmt.Errorf("error building client for endpoint %s: %+v", endpoint, err)
	}

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

	d.Set("key", key)
	d.Set("value", result.Value)
	d.Set("label", result.Label)
	d.Set("endpoint", endpoint)

	log.Printf("[INFO] KV has been fetched %s/%s/%s=%s", endpoint, label, key, result.Value)

	return nil
}

func resourceKeyValueUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating resource %s", d.Id())

	endpoint, label, key := parseID(d.Id())

	value := d.Get("value").(string)

	cl, err := getOrReuseClient(endpoint, meta.(func(endpoint string) (*client.Client, error)))
	if err != nil {
		return fmt.Errorf("error building client for endpoint %s: %+v", endpoint, err)
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

	return resourceKeyValueRead(d, meta)
}

func resourceKeyValueDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting resource %s", d.Id())

	endpoint, label, key := parseID(d.Id())

	cl, err := getOrReuseClient(endpoint, meta.(func(endpoint string) (*client.Client, error)))
	if err != nil {
		return fmt.Errorf("error building client for endpoint %s: %+v", endpoint, err)
	}

	_, err = cl.DeleteKeyValue(label, key)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func formatID(endpoint string, label string, key string) (string, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("unable to parse the given endpoint %s", endpoint)
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

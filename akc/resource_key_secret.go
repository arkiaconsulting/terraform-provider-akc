package akc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeySecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeySecretCreate,
		Read:   resourceKeySecretRead,
		Update: resourceKeySecretUpdate,
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
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(readTimeout),
		},
	}
}

func resourceKeySecretCreate(d *schema.ResourceData, meta interface{}) error {

	endpoint := d.Get("endpoint").(string)
	key := d.Get("key").(string)
	value := d.Get("secret_id").(string)
	label := d.Get("label").(string)
	trim := d.Get("latest_version").(bool)

	cl, err := getClient(endpoint, meta.(func(endpoint string) (*client.Client, error)))
	if err != nil {
		return fmt.Errorf("error building client for endpoint %s: %+v", endpoint, err)
	}

	if trim {
		value = trimVersion(value)
	}

	if d.IsNewResource() {
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			_, err = cl.GetKeyValue(label, key)

			if err != nil {
				if client.IsNotFound(err) {
					log.Printf("[INFO] retrying to get key-secret '%s/%s' because: %s", label, key, err)

					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err == nil {
			return fmt.Errorf("the resource needs to be imported: %s", "akc_key_value")
		}
	}

	_, err = cl.SetKeyValueSecret(key, value, label)
	if err != nil {
		return err
	}

	id, err := formatID(endpoint, label, key)
	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("value", value)

	return resourceKeySecretRead(d, meta)
}

func resourceKeySecretUpdate(d *schema.ResourceData, meta interface{}) error {
	endpoint, label, key := parseID(d.Id())

	value := d.Get("secret_id").(string)
	trim := d.Get("latest_version").(bool)

	cl, err := getClient(endpoint, meta.(func(endpoint string) (*client.Client, error)))
	if err != nil {
		return fmt.Errorf("error building client for endpoint %s: %+v", endpoint, err)
	}

	if trim {
		value = trimVersion(value)
	}

	_, err = cl.SetKeyValueSecret(key, value, label)
	if err != nil {
		return err
	}

	id, err := formatID(endpoint, label, key)
	if err != nil {
		return err
	}

	d.SetId(id)
	d.Set("value", value)

	return resourceKeySecretRead(d, meta)
}

func resourceKeySecretRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading resource %s", d.Id())

	endpoint, label, key := parseID(d.Id())

	cl, err := getClient(endpoint, meta.(func(endpoint string) (*client.Client, error)))
	if err != nil {
		return fmt.Errorf("error building client for endpoint %s: %+v", endpoint, err)
	}

	var kv client.KeyValueResponse
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		kv, err = cl.GetKeyValue(label, key)

		if err != nil {
			if client.IsNotFound(err) {
				log.Printf("[INFO] retrying to get key-secret '%s/%s' because: %s", label, key, err)

				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		log.Printf("[INFO] key-secret not found, removing from state: %s/%s/%s\n", endpoint, label, key)
		d.SetId("")
		return nil
	}

	var wrapper keyVaultReferenceValue
	err = json.Unmarshal([]byte(kv.Value), &wrapper)
	if err != nil {
		return err
	}

	d.Set("key", key)
	d.Set("value", wrapper.URI)
	d.Set("label", label)
	d.Set("endpoint", endpoint)

	log.Printf("[INFO] the key-secret '%s/%s/%s=%s' was read successfuly\n", endpoint, label, key, wrapper.URI)

	return nil
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

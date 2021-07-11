package akc

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFeature() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeatureCreate,
		Read:   resourceFeatureRead,
		Update: resourceFeatureUpdate,
		Delete: resourceFeatureDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"label": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  client.LabelNone,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(readTimeout),
		},
	}
}

func resourceFeatureCreate(d *schema.ResourceData, meta interface{}) error {

	endpoint := d.Get("endpoint").(string)
	name := d.Get("name").(string)
	label := d.Get("label").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)

	cl, err := getClient(endpoint, meta.(func(endpoint string) (*client.Client, error)))
	if err != nil {
		return fmt.Errorf("error building client for endpoint %s: %+v", endpoint, err)
	}

	if d.IsNewResource() {
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			_, err = cl.GetFeature(label, name)

			if err != nil {
				if client.IsNotFound(err) {
					log.Printf("[INFO] retrying to get feature '%s/%s' because: %s", label, name, err)

					return resource.RetryableError(err)
				}

				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err == nil {
			return fmt.Errorf("the resource needs to be imported: %s", d.Id())
		}
	}

	_, err = cl.SetFeature(name, label, enabled, description)
	if err != nil {
		return err
	}

	id, err := formatFeatureID(endpoint, label, name)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceFeatureRead(d, meta)
}

func resourceFeatureRead(d *schema.ResourceData, meta interface{}) error {
	endpoint, label, name := parseFeatureID(d.Id())

	cl, err := getClient(endpoint, meta.(func(endpoint string) (*client.Client, error)))
	if err != nil {
		return fmt.Errorf("error building client for endpoint %s: %+v", endpoint, err)
	}

	var feature client.FeatureResponse
	err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		feature, err = cl.GetFeature(label, name)

		if err != nil {
			if client.IsNotFound(err) {
				log.Printf("[INFO] retrying to get feature '%s:%s' because: %s", label, name, err)

				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		log.Printf("[INFO] KV not found, removing from state: %s/%s/%s", endpoint, label, name)
		d.SetId("")
		return nil
	}

	d.Set("endpoint", endpoint)
	d.Set("name", name)
	d.Set("label", label)
	d.Set("description", feature.Description)
	d.Set("enabled", feature.Enabled)

	log.Printf("[INFO] KV has been fetched %s/%s/%s", endpoint, label, name)

	return nil
}

func resourceFeatureUpdate(d *schema.ResourceData, meta interface{}) error {

	endpoint, label, name := parseFeatureID(d.Id())
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)

	cl, err := getClient(endpoint, meta.(func(endpoint string) (*client.Client, error)))
	if err != nil {
		return fmt.Errorf("error building client for endpoint %s: %+v", endpoint, err)
	}

	_, err = cl.SetFeature(name, label, enabled, description)
	if err != nil {
		return err
	}

	return resourceFeatureRead(d, meta)
}

func resourceFeatureDelete(d *schema.ResourceData, meta interface{}) error {
	endpoint, label, name := parseFeatureID(d.Id())

	cl, err := getClient(endpoint, meta.(func(endpoint string) (*client.Client, error)))
	if err != nil {
		return fmt.Errorf("error building client for endpoint %s: %+v", endpoint, err)
	}

	_, err = cl.DeleteFeature(label, name)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func formatFeatureID(endpoint string, label string, name string) (string, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("unable to parse the given endpoint %s", endpoint)
	}

	host := url.Host

	return fmt.Sprintf("%s/feature/%s/%s", host, label, name), nil
}

func parseFeatureID(id string) (endpoint string, label string, name string) {
	split := strings.Split(id, "/")

	endpoint = fmt.Sprintf("https://%s", split[0])
	label = split[2]
	name = split[3]

	return endpoint, label, name
}

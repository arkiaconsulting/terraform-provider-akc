package akc

import (
	"time"

	"github.com/arkiaconsulting/terraform-provider-akc/client"
)

func getClient(endpoint string, clientBuilder func(endpoint string) (*client.Client, error)) (*client.Client, error) {
	cl, err := clientBuilder(endpoint)
	if err != nil {
		return nil, err
	}
	return cl, nil
}

const readTimeout = 20 * time.Second

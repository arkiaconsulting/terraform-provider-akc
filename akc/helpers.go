package akc

import "github.com/arkiaconsulting/terraform-provider-akc/client"

var clientsCache map[string]*client.Client = make(map[string]*client.Client)

func getOrReuseClient(endpoint string, clientBuilder func(endpoint string) (*client.Client, error)) (*client.Client, error) {

	cl := clientsCache[endpoint]

	// client is in cache
	if cl != nil {
		return cl, nil
	}

	// no client for the given endpoint in cache
	cl, err := clientBuilder(endpoint)
	if err != nil {
		return nil, err
	}

	clientsCache[endpoint] = cl

	return cl, nil
}

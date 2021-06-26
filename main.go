package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"terraform-provider-akc/akc"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: akc.Provider,
	})
}

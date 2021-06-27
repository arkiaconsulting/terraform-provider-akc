package akc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testProviders = map[string]*schema.Provider{
	"akc": Provider(),
}

const appConfigHost = "arkia.azconfig.io"
const endpointUnderTest = "https://" + appConfigHost

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/vr-ibm/terraform-quantum-aerospace/internal/provider"
)

func testAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"quantum": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}

func testAccPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("IONQ_API_KEY") == "" {
		t.Skip("IONQ_API_KEY not set — skipping acceptance tests")
	}
}

// Ensure the test helper compiles even when no test in this file uses resource directly.
var _ = resource.Test

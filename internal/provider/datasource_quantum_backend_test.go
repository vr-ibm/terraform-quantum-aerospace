package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuantumBackendDataSource(t *testing.T) {
	testAccPreCheck(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccQuantumBackendDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.quantum_backend.simulator", "name", "simulator"),
					resource.TestCheckResourceAttr("data.quantum_backend.simulator", "cloud_provider", "ionq"),
					resource.TestCheckResourceAttrSet("data.quantum_backend.simulator", "qubits"),
					resource.TestCheckResourceAttrSet("data.quantum_backend.simulator", "status"),
				),
			},
		},
	})
}

func testAccQuantumBackendDataSourceConfig() string {
	return `
provider "quantum" {}

data "quantum_backend" "simulator" {
  name = "simulator"
}
`
}

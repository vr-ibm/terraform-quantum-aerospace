package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQuantumJobResource(t *testing.T) {
	testAccPreCheck(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccQuantumJobResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("quantum_job.route_optimization", "backend", "simulator"),
					resource.TestCheckResourceAttr("quantum_job.route_optimization", "shots", "100"),
					resource.TestCheckResourceAttrSet("quantum_job.route_optimization", "id"),
					resource.TestCheckResourceAttr("quantum_job.route_optimization", "status", "completed"),
					resource.TestCheckResourceAttrSet("quantum_job.route_optimization", "results"),
				),
			},
		},
	})
}

func TestAccQuantumJobResource_Import(t *testing.T) {
	testAccPreCheck(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccQuantumJobResourceConfig(),
			},
			{
				ResourceName:            "quantum_job.route_optimization",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"circuit", "shots", "timeout"},
			},
		},
	})
}

func testAccQuantumJobResourceConfig() string {
	return `
provider "quantum" {}

resource "quantum_job" "route_optimization" {
  backend = "simulator"
  shots   = 100
  timeout = "5m"
  circuit = <<-QASM
    OPENQASM 3.0;
    include "stdgates.inc";
    qubit[2] q;
    bit[2] c;
    // Minimal 2-qubit route segment interaction
    // Encodes fuel cost conflict between two candidate flight paths
    h q[0];
    cx q[0], q[1];
    c = measure q;
  QASM
}
`
}

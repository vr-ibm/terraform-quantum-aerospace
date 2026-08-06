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
			// Create and poll until complete
			{
				Config: testAccQuantumJobResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("quantum_job.bell_state", "backend", "simulator"),
					resource.TestCheckResourceAttr("quantum_job.bell_state", "shots", "100"),
					resource.TestCheckResourceAttrSet("quantum_job.bell_state", "id"),
					resource.TestCheckResourceAttr("quantum_job.bell_state", "status", "completed"),
					resource.TestCheckResourceAttrSet("quantum_job.bell_state", "results"),
				),
			},
			// Read back from state
			{
				ResourceName:      "quantum_job.bell_state",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"circuit", "timeout", "shots"},
			},
		},
	})
}

func testAccQuantumJobResourceConfig() string {
	return `
provider "quantum" {}

resource "quantum_job" "bell_state" {
  backend = "simulator"
  shots   = 100
  timeout = "5m"
  circuit = <<-QASM
    OPENQASM 3.0;
    include "stdgates.inc";
    qubit[2] q;
    bit[2] c;
    h q[0];
    cx q[0], q[1];
    c = measure q;
  QASM
}
`
}

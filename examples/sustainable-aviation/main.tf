terraform {
  required_providers {
    quantum = {
      source  = "vr-ibm/quantum"
      version = "~> 0.1"
    }
  }
}

provider "quantum" {
  backend = "ionq"
}

data "quantum_backend" "simulator" {
  name = "simulator"
}

output "backend_qubits" {
  value = data.quantum_backend.simulator.qubits
}

output "backend_status" {
  value = data.quantum_backend.simulator.status
}

# Reusable circuit for sustainable route selection
resource "quantum_circuit" "route_qaoa" {
  name        = "qaoa-route-emission-4q"
  description = "4-qubit QAOA ansatz for commercial flight route optimization minimizing fuel burn and CO2 emissions"
  qubits      = 4
  body        = <<-QASM
    OPENQASM 3.0;
    include "stdgates.inc";
    qubit[4] q;
    bit[4] c;
    // Each qubit represents a candidate route segment
    // q[0]: direct route (shorter, higher altitude)
    // q[1]: wind-optimized route (longer, lower fuel burn)
    // q[2]: SAF-available corridor (sustainable fuel access)
    // q[3]: congestion-avoidance route (reduced holding patterns)
    // Initial superposition over all route combinations
    h q[0];
    h q[1];
    h q[2];
    h q[3];
    // Problem unitary: ZZ interactions encoding emission costs
    // between conflicting route segments
    cx q[0], q[1];
    rz(0.4) q[1];
    cx q[0], q[1];
    cx q[1], q[2];
    rz(0.6) q[2];
    cx q[1], q[2];
    cx q[2], q[3];
    rz(0.3) q[3];
    cx q[2], q[3];
    // Mixer unitary: enables exploration of route combinations
    rx(0.35) q[0];
    rx(0.35) q[1];
    rx(0.35) q[2];
    rx(0.35) q[3];
    // Measure optimal route selection
    c = measure q;
  QASM
}

# Submit route optimization job
resource "quantum_job" "route_optimization" {
  backend = data.quantum_backend.simulator.name
  shots   = 1000
  timeout = "5m"
  circuit = quantum_circuit.route_qaoa.body
}

output "circuit_hash" {
  value = quantum_circuit.route_qaoa.hash
}

output "job_id" {
  value = quantum_job.route_optimization.id
}

output "job_status" {
  value = quantum_job.route_optimization.status
}

output "job_results" {
  value = quantum_job.route_optimization.results
}

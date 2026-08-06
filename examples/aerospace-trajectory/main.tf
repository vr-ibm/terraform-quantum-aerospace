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

# Look up available backend metadata before submitting
data "quantum_backend" "simulator" {
  name = "simulator"
}

output "backend_qubits" {
  value = data.quantum_backend.simulator.qubits
}

output "backend_status" {
  value = data.quantum_backend.simulator.status
}

# Define circuit as a reusable, named resource
resource "quantum_circuit" "trajectory_qaoa" {
  name        = "qaoa-trajectory-4q"
  description = "4-qubit QAOA ansatz for satellite trajectory waypoint optimization"
  qubits      = 4
  body        = <<-QASM
    OPENQASM 3.0;
    include "stdgates.inc";
    qubit[4] q;
    bit[4] c;
    // Initial superposition
    h q[0];
    h q[1];
    h q[2];
    h q[3];
    // Problem unitary (ZZ interactions encoding trajectory constraints)
    cx q[0], q[1];
    rz(0.5) q[1];
    cx q[0], q[1];
    cx q[1], q[2];
    rz(0.5) q[2];
    cx q[1], q[2];
    cx q[2], q[3];
    rz(0.5) q[3];
    cx q[2], q[3];
    // Mixer unitary
    rx(0.3) q[0];
    rx(0.3) q[1];
    rx(0.3) q[2];
    rx(0.3) q[3];
    // Measure
    c = measure q;
  QASM
}

# Submit the circuit as a job — references the circuit resource
resource "quantum_job" "trajectory_optimization" {
  backend = data.quantum_backend.simulator.name
  shots   = 1000
  timeout = "5m"
  circuit = quantum_circuit.trajectory_qaoa.body
}

output "circuit_hash" {
  value = quantum_circuit.trajectory_qaoa.hash
}

output "job_id" {
  value = quantum_job.trajectory_optimization.id
}

output "job_status" {
  value = quantum_job.trajectory_optimization.status
}

output "job_results" {
  value = quantum_job.trajectory_optimization.results
}

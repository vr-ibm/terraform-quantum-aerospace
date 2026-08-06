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
data "quantum_backend" "aria" {
  name = "ionq.aria"
}

output "backend_qubits" {
  value = data.quantum_backend.aria.qubits
}

output "backend_status" {
  value = data.quantum_backend.aria.status
}

# Aerospace use case: QAOA circuit for satellite trajectory optimization
# Circuit encoded as OpenQASM 3.0
resource "quantum_job" "trajectory_optimization" {
  backend = data.quantum_backend.aria.name
  shots   = 1000
  circuit = <<-QASM
    OPENQASM 3.0;
    include "stdgates.inc";
    // 4-qubit QAOA ansatz for trajectory waypoint optimization
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

output "job_id" {
  value = quantum_job.trajectory_optimization.id
}

output "job_status" {
  value = quantum_job.trajectory_optimization.status
}

output "job_results" {
  value = quantum_job.trajectory_optimization.results
}

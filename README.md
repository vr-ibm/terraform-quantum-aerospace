# terraform-quantum-aerospace

[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](LICENSE)

Terraform provider for quantum cloud backends (IonQ on GCP, AWS Braket) applied to aerospace
optimization problems. Explores IaC design patterns for non-CRUD resource lifecycles — job-based
execution, probabilistic state, and ephemeral circuits — as quantum services expand across major
cloud platforms.

## Motivation

Quantum cloud services (IonQ via GCP Marketplace, AWS Braket, Azure Quantum) are fundamentally
different from traditional CRUD APIs:

- **Jobs are immutable** — you submit a circuit, you don't update it in place
- **State is probabilistic** — results are distributions, not deterministic values
- **Backends are ephemeral** — calibration data and availability change without user action
- **Execution is async** — jobs queue, run, and complete on timescales Terraform doesn't expect

These characteristics challenge standard Terraform provider design patterns. This project explores
how IaC tooling should model non-traditional cloud resources, using aerospace optimization workloads
as a reference domain.

## Provider Resources

| Type | Kind | Description |
|------|------|-------------|
| `quantum_job` | Resource | Submit a quantum circuit to a cloud backend |
| `quantum_backend` | Data Source | Query backend metadata (qubits, status, shot limits) |

## Usage

```hcl
provider "quantum" {
  backend = "ionq"
}

data "quantum_backend" "aria" {
  name = "ionq.aria"
}

resource "quantum_job" "trajectory" {
  backend = data.quantum_backend.aria.name
  shots   = 1000
  circuit = file("circuits/qaoa_trajectory.qasm")
}
```

See [`examples/aerospace-trajectory/`](examples/aerospace-trajectory/) for a complete working
configuration.

## Building

```sh
go build ./...
```

## Research

This provider serves as the implementation artifact for a research paper exploring IaC design
patterns for quantum cloud services. Key research questions:

- How should Terraform model resources with non-deterministic state?
- What identity patterns apply to job-based (submit-and-poll) APIs vs. CRUD APIs?
- How do list/import patterns work for ephemeral quantum resources?
- What Plugin Framework gaps emerge when modeling non-traditional lifecycles?

## Aerospace workloads

The reference workloads focus on problems where quantum approaches show promise:

- **Satellite trajectory optimization** — QAOA for waypoint selection under orbital constraints
- **Constellation scheduling** — Quantum annealing for multi-satellite communication scheduling
- **Structural analysis** — VQE for materials simulation in aerospace composites

## Roadmap

- [ ] IonQ API client (GCP Marketplace auth)
- [ ] AWS Braket API client
- [ ] Async job polling with configurable timeout
- [ ] `quantum_circuit` resource (reusable circuit definitions)
- [ ] `quantum_job` list resource for `terraform query`
- [ ] Acceptance tests against simulator backends
- [ ] Benchmark: classical vs. quantum solver for trajectory problems

## License

This project is licensed under the [Mozilla Public License 2.0](LICENSE).
Source files may be used, modified, and distributed under the terms of the MPL-2.0.
Each source file must retain the Exhibit A header notice; see the [`LICENSE`](LICENSE) file for the full terms.

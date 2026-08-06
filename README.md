# terraform-quantum-aerospace

[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](LICENSE)

A Terraform provider for quantum cloud infrastructure applied to sustainable commercial aviation.

## Research Focus

Commercial aviation accounts for ~2.5% of global CO2 emissions. Route optimization, fleet
scheduling, and fuel allocation are combinatorial problems where quantum computing may offer
advantages at scale. This project models quantum optimization jobs as Terraform-managed
infrastructure, applying QAOA to flight route selection that minimizes fuel burn and emissions.

## Design Patterns

This provider implements four IaC design patterns for non-CRUD cloud services:

| Pattern | Resource | Description |
|---------|----------|-------------|
| Immutable job lifecycle | `quantum_job` | `RequiresReplace` on circuit/backend, no Update |
| Async polling | `quantum_job` | Exponential backoff, configurable timeout |
| Partial import | `quantum_job` | ID resolves status/results, circuit is write-only |
| Config-only resource | `quantum_circuit` | No backing API, state-only, content-addressed |

## Application: Sustainable Route Optimization

The example encodes a flight route selection problem as QAOA:

- 4 qubits represent candidate route segments (direct, wind-optimized, SAF corridor, congestion-avoidance)
- ZZ interactions encode emission cost conflicts between segments
- Measurement outcomes approximate the lowest-emission route combination

## Quick Start

```bash
export IONQ_API_KEY=your_key
cd examples/sustainable-aviation
terraform init
terraform apply
```

## Provider Resources

- `quantum_job` — Submit and poll a quantum circuit execution
- `quantum_circuit` — Define a reusable, named circuit (config-only)
- `data.quantum_backend` — Read backend metadata (qubits, status)

## Testing

```bash
IONQ_API_KEY=your_key TF_ACC=1 go test ./internal/provider/... -v -timeout 15m
```

## License

This project is licensed under the [Mozilla Public License 2.0](LICENSE).
Source files may be used, modified, and distributed under the terms of the MPL-2.0.
Each source file must retain the Exhibit A header notice; see the [`LICENSE`](LICENSE) file for the full terms.

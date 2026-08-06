# Aerospace Trajectory Optimization Example

This example demonstrates using the `quantum` provider to submit a QAOA
circuit for satellite trajectory waypoint optimization to IonQ Aria.

## Usage

```hcl
terraform init
terraform plan
terraform apply
```

## What this does

- Queries IonQ Aria backend metadata (qubit count, status, shot limits)
- Submits a 4-qubit QAOA ansatz encoding trajectory constraints
- Returns job ID, status, and measurement results

## Research context

QAOA (Quantum Approximate Optimization Algorithm) is a candidate algorithm
for aerospace trajectory optimization problems, which are NP-hard in the
general case. This example encodes a simplified 4-waypoint problem as a
QAOA circuit and submits it to a real quantum backend via Terraform.

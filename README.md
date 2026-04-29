# driftwatch

> Detect configuration drift between running Kubernetes workloads and their source Helm charts.

---

## Overview

`driftwatch` compares the live state of Kubernetes workloads against the expected state defined in their source Helm charts, helping you catch unintended changes before they cause problems.

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/driftwatch/releases).

## Usage

```bash
# Check for drift in a specific namespace
driftwatch --namespace production --chart ./charts/myapp

# Compare against a remote Helm chart
driftwatch --namespace staging --repo https://charts.example.com --chart myapp --version 1.4.2

# Output results as JSON
driftwatch --namespace production --chart ./charts/myapp --output json
```

### Example Output

```
[DRIFT DETECTED] deployment/api-server
  - spec.replicas: expected 3, got 1
  - spec.template.spec.containers[0].image: expected myapp:1.4.2, got myapp:1.3.0

[OK] deployment/worker
[OK] service/api-server
```

## Configuration

| Flag | Description | Default |
|------|-------------|---------|
| `--namespace` | Kubernetes namespace to inspect | `default` |
| `--chart` | Path or name of the Helm chart | required |
| `--repo` | Helm chart repository URL | — |
| `--version` | Chart version to compare against | latest |
| `--output` | Output format (`text`, `json`) | `text` |
| `--kubeconfig` | Path to kubeconfig file | `~/.kube/config` |

## Requirements

- Go 1.21+
- `kubectl` configured with access to your cluster
- Helm 3

## License

This project is licensed under the [MIT License](LICENSE).
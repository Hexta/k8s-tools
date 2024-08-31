# k8s-tools

A toolset for inspecting Kubernetes clusters.

## Getting Started

### Installation

```shell
go install github.com/Hexta/k8s-tools/cmd/k8s-tools@latest
```

## Usage

### Examples

* Print nodes utilisation in zone eu-central-1a
    ```shell
    k8s-tools node utilisation -l topology.kubernetes.io/zone=eu-central-1a
    ```

## Contributing

Please see our [contributing guidelines](CONTRIBUTING.md).

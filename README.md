# k8s-tools

A toolset for inspecting Kubernetes clusters.

## Getting Started

### Installation

> [!IMPORTANT]
> CGO_ENABLED=1 is required for DuckDB.

```shell
CGO_ENABLED=1 go install github.com/Hexta/k8s-tools/cmd/k8s-tools@latest
```

## Usage

### Examples

#### Ad-hoc tools

* Print nodes utilisation in zone eu-central-1a
    ```shell
    k8s-tools node utilisation -l topology.kubernetes.io/zone=eu-central-1a
    ```

#### Internal DB

1. Init DB.
    ```shell
    k8s-tools db init
    ```
2. Query data.
    ```shell
    k8s-tools db query "select name, cpu_utilisation, age(creation_ts)::varchar as age from k8s.nodes order by cpu_utilisation asc limit 10"
    ```

## Contributing

Please see our [contributing guidelines](CONTRIBUTING.md).

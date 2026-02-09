# k8s-tools

A toolset for querying and inspecting Kubernetes clusters via a built-in SQL interface.

k8s-tools provide a SQL interface for querying information about various Kubernetes resources, including:
* Containers
* Custom Resources
* DaemonSets
* Deployments
* EndpointSlices
* Events
* HorizontalPodAutoscalers
* Nodes
* PersistentVolumeClaims
* PersistentVolumes
* Pods
* Services
* StatefulSets

## Features

- SQL-based querying of Kubernetes resources
- Interactive TUI for database exploration
- Support for multiple output formats (JSON, Table, Vertical)

## Database Schema

The internal database contains the following schemas:

- `k8s`: Main schema containing all Kubernetes resource tables
  - `nodes`
  - `pods`
  - `deployments`
  - (etc...)

For detailed schema information, see [DB Documentation](docs/db/index.md)

## Getting Started

### Prerequisites

- Go 1.25 or later
- CGO-enabled environment
- Kubernetes cluster access configured

### Installation

> [!IMPORTANT]
> CGO_ENABLED=1 is required for DuckDB.

```shell
CGO_ENABLED=1 go install github.com/Hexta/k8s-tools/cmd/k8s-tools@latest
```

## Usage

[CLI Documentation](docs/cli/k8s-tools.md)

### Examples

#### Ad-hoc tools

* Print nodes utilisation in zone eu-central-1a
    ```shell
    k8s-tools node utilisation -l topology.kubernetes.io/zone=eu-central-1a
    ```

1. Init DB.
    ```shell
    k8s-tools db init
    ```
2. Query data.
    * Top10 CPU underutilized nodes
      ```shell
      k8s-tools db query "select name, cpu_utilisation, labels['karpenter.sh/nodepool'] from k8s.nodes order by cpu_utilisation asc limit 10"
      ```
3. Run SQL TUI.
    ```shell
    k8s-tools db tui
    ```
   ![DB TUI](docs/db-tui-0.png)

##### Useful queries
1. List tables
   ```shell
   k8s-tools db query "select table_schema, table_name, table_type from information_schema.tables where table_schema='k8s' order by table_name""
   ```

## Contributing

Contributions are welcome! Please see our [contributing guidelines](CONTRIBUTING.md).

## License

This project is licensed under the Apache License 2.0.

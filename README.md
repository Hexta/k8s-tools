# k8s-tools

A toolset for inspecting Kubernetes clusters.

Provides SQL interface for querying information about K8s:
* Containers
* Deployments
* Pods
* Nodes

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

k8s-tools can save K8s cluster state in internal DB to ease analyzing.
[DB Documentation](docs/db/index.md)

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

## Contributing

Please see our [contributing guidelines](CONTRIBUTING.md).

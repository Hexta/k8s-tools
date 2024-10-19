## k8s-tools db

DB with K8s information

### Options

```
  -h, --help   help for db
```

### Options inherited from parent commands

```
      --cache-dir string                      cache directory
      --k8s-retry-initial-interval duration   Initial interval for Kubernetes API retry (default 1s)
      --k8s-retry-jitter-percent uint         Jitter percent for Kubernetes API retry (default 50)
      --k8s-retry-max-attempts uint           Maximum number of attempts for Kubernetes API retry (default 5)
      --k8s-retry-max-interval duration       Maximum interval between retries for Kubernetes API (default 10s)
      --kubeconfig string                     kubeconfig file
  -o, --output Format                         output format (json, table, vertical) (default table)
  -v, --verbose                               verbose
```

### SEE ALSO

* [k8s-tools](k8s-tools.md)	 - K8s toolbox
* [k8s-tools db init](k8s-tools_db_init.md)	 - Init DB
* [k8s-tools db query](k8s-tools_db_query.md)	 - Query DB
* [k8s-tools db tui](k8s-tools_db_tui.md)	 - TUI for DB


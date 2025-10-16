## k8s-tools db query

Query DB

```
k8s-tools db query QUERY [flags]
```

### Examples

```
query "SELECT * FROM k8s.nodes LIMIT 10"
```

### Options

```
  -h, --help   help for query
```

### Options inherited from parent commands

```
      --cache-dir string                      cache directory
      --context string                        context
      --k8s-retry-initial-interval duration   Initial interval for Kubernetes API retry (default 1s)
      --k8s-retry-jitter-percent uint         Jitter percent for Kubernetes API retry (default 50)
      --k8s-retry-max-attempts uint           Maximum number of attempts for Kubernetes API retry (default 5)
      --k8s-retry-max-interval duration       Maximum interval between retries for Kubernetes API (default 10s)
      --kubeconfig string                     kubeconfig file
      --no-headers                            do not print headers
  -o, --output Format                         output format (json, table, vertical) (default table)
  -v, --verbose                               verbose
```

### SEE ALSO

* [k8s-tools db](k8s-tools_db.md)	 - DB with K8s information


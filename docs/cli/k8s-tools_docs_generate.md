## k8s-tools docs generate

Generate documentation

```
k8s-tools docs generate [output dir] [flags]
```

### Options

```
  -h, --help   help for generate
```

### Options inherited from parent commands

```
      --cache-dir string                      cache directory
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

* [k8s-tools docs](k8s-tools_docs.md)	 - Documentation for k8s-tools


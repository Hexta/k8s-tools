## k8s-tools completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	k8s-tools completion fish | source

To load completions for every new session, execute once:

	k8s-tools completion fish > ~/.config/fish/completions/k8s-tools.fish

You will need to start a new shell for this setup to take effect.


```
k8s-tools completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
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

* [k8s-tools completion](k8s-tools_completion.md)	 - Generate the autocompletion script for the specified shell


## k8s-tools completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(k8s-tools completion zsh)

To load completions for every new session, execute once:

#### Linux:

	k8s-tools completion zsh > "${fpath[1]}/_k8s-tools"

#### macOS:

	k8s-tools completion zsh > $(brew --prefix)/share/zsh/site-functions/_k8s-tools

You will need to start a new shell for this setup to take effect.


```
k8s-tools completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
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

* [k8s-tools completion](k8s-tools_completion.md)	 - Generate the autocompletion script for the specified shell


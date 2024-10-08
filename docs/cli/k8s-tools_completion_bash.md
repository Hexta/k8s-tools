## k8s-tools completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(k8s-tools completion bash)

To load completions for every new session, execute once:

#### Linux:

	k8s-tools completion bash > /etc/bash_completion.d/k8s-tools

#### macOS:

	k8s-tools completion bash > $(brew --prefix)/etc/bash_completion.d/k8s-tools

You will need to start a new shell for this setup to take effect.


```
k8s-tools completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --cache-dir string    cache directory
      --kubeconfig string   kubeconfig file
  -v, --verbose             verbose
```

### SEE ALSO

* [k8s-tools completion](k8s-tools_completion.md)	 - Generate the autocompletion script for the specified shell


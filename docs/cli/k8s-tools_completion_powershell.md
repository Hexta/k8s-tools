## k8s-tools completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	k8s-tools completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
k8s-tools completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
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

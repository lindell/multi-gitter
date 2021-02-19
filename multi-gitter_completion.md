## multi-gitter completion

Print shell autocompletion scripts for multi-gitter

### Synopsis

To load completions:
Bash:
$ source <(multi-gitter completion bash)
# To load completions for each session, execute once:
Linux:
  $ multi-gitter completion bash > /etc/bash_completion.d/multi-gitter
MacOS:
  $ multi-gitter completion bash > /usr/local/etc/bash_completion.d/multi-gitter
Zsh:
# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:
$ echo "autoload -U compinit; compinit" >> ~/.zshrc
# To load completions for each session, execute once:
$ multi-gitter completion zsh > "${fpath[1]}/_multi-gitter"
# You will need to start a new shell for this setup to take effect.
Fish:
$ multi-gitter completion fish | source
# To load completions for each session, execute once:
$ multi-gitter completion fish > ~/.config/fish/completions/multi-gitter.fish


```
multi-gitter completion [bash|zsh|fish]
```

### Options

```
  -h, --help   help for completion
```

### SEE ALSO

* [multi-gitter](multi-gitter.md)	 - Multi gitter is a tool for making changes into multiple git repositories


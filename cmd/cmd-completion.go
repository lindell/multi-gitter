package cmd

import "github.com/spf13/cobra"

// CompletionCmd generates completion scripts
func CompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish]",
		Short: "Print shell autocompletion scripts for multi-gitter.",
		Long: `To load completions:
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
`,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish"},
		Args:                  cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			switch args[0] {
			case "bash":
				err = cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				err = cmd.Root().GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				err = cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			}

			return err
		},
	}
}

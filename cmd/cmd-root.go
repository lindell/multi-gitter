package cmd

import (
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

// RootCmd is the root command containing all subcommands
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi-gitter",
		Short: "Multi gitter is a tool for making changes into multiple git repositories.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd) // Bind configs that are not flags
		},
	}

	cmd.AddCommand(RunCmd())
	cmd.AddCommand(StatusCmd())
	cmd.AddCommand(MergeCmd())
	cmd.AddCommand(CloseCmd())
	cmd.AddCommand(PrintCmd())
	cmd.AddCommand(VersionCmd())

	return cmd
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

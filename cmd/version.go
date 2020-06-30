package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionCmd prints the version of multi-gitter
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of multi-gitter.",
	Long:  "Get the version of multi-gitter.",
	Args:  cobra.NoArgs,
	Run:   version,
}

// Version is the current version of multigitter (set by main.go)
var Version string

func version(cmd *cobra.Command, args []string) {
	fmt.Println(Version)
}

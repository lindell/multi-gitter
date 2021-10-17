package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

// VersionCmd prints the version of multi-gitter
func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Get the version of multi-gitter.",
		Long:  "Get the version of multi-gitter.",
		Args:  cobra.NoArgs,
		Run:   version,
	}

	return cmd
}

// Version is the current version of multigitter (set by main.go)
var Version string

// BuildDate is the time the build was made (set by main.go)
var BuildDate time.Time

// Commit is the commit the build was made on (set by main.go)
var Commit string

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("multi-gitter version: %s\n", Version)
	fmt.Printf("Release-Date: %s\n", BuildDate.Format("2006-01-02"))
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("Commit: %s\n", Commit)
}

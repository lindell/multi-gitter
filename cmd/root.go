package cmd

import (
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

// RootCmd is the root command containing all subcommands
var RootCmd = &cobra.Command{
	Use:   "multi-gitter",
	Short: "Multi gitter is a tool for making changes into multiple git repositories",
}

func init() {
	RootCmd.PersistentFlags().StringP("gh-base-url", "g", "https://api.github.com/", "Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.")
	RootCmd.PersistentFlags().StringP("token", "T", "", "The GitHub personal access token. Can also be set using the GITHUB_TOKEN environment variable.")

	RootCmd.AddCommand(RunCmd)

	rand.Seed(time.Now().UTC().UnixNano())
}

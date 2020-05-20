package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "multi-gitter",
	Short: "Multi gitter is a tool for making changes into multiple git repositories",
}

func init() {
	RootCmd.PersistentFlags().StringP("gh-base-url", "g", "https://api.github.com/", "Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.")
	RootCmd.PersistentFlags().StringP("token", "T", "", "The GitHub personal access token. Can also be set using the GITHUB_TOKEN enviroment variable.")

	RootCmd.AddCommand(RunCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

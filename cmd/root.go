package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "multi-gitter",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
}

func init() {
	rootCmd.PersistentFlags().StringP("gh-base-url", "g", "https://api.github.com/", "Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.")
	rootCmd.PersistentFlags().StringP("token", "T", "", "The GitHub personal access token. Can also be set using the GITHUB_TOKEN enviroment variable.")

	rootCmd.AddCommand(runCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

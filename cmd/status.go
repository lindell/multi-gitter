package cmd

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/lindell/multi-gitter/internal/github"
	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/spf13/cobra"
)

// StatusCmd gets statuses of pull requests
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of pull requests",
	Long:  "Get the status of all pull requests with a specified branch name in an organization.",
	Args:  cobra.NoArgs,
	RunE:  status,
}

func init() {
	StatusCmd.Flags().StringP("branch", "B", "multi-gitter-branch", "The name of the branch where changes are committed.")
	StatusCmd.Flags().StringP("org", "o", "", "The name of the GitHub organization.")
}

func status(cmd *cobra.Command, args []string) error {
	flag := cmd.Flags()

	ghBaseURL, _ := flag.GetString("gh-base-url")
	token, _ := flag.GetString("token")
	branchName, _ := flag.GetString("branch")
	org, _ := flag.GetString("org")

	if token != "" {
		if ght := os.Getenv("GITHUB_TOKEN"); ght != "" {
			token = ght
		}
	}

	if token == "" {
		return errors.New("either the --token flag or the GITHUB_TOKEN environment variable has to be set")
	}

	if org == "" {
		return errors.New("no organization set")
	}

	vc, err := github.New(token, ghBaseURL)
	if err != nil {
		return err
	}

	statuser := multigitter.Statuser{
		VersionController: vc,

		FeatureBranch: branchName,
		OrgName:       org,
	}

	err = statuser.Statuses(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

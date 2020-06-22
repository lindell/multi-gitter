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

// MergeCmd merges pull requests
var MergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge pull requests.",
	Long:  "Merge pull requests with a specified branch name in an organization and with specified conditions.",
	Args:  cobra.NoArgs,
	RunE:  merge,
}

func init() {
	MergeCmd.Flags().StringP("branch", "B", "multi-gitter-branch", "The name of the branch where changes are committed.")
	MergeCmd.Flags().StringP("org", "o", "", "The name of the GitHub organization.")
}

func merge(cmd *cobra.Command, args []string) error {
	flag := cmd.Flags()

	ghBaseURL, _ := flag.GetString("gh-base-url")
	token, _ := flag.GetString("token")
	branchName, _ := flag.GetString("branch")
	org, _ := flag.GetString("org")

	if token == "" {
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

	statuser := multigitter.Merger{
		VersionController: vc,

		FeatureBranch: branchName,
		OrgName:       org,
	}

	err = statuser.Merge(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

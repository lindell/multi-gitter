package cmd

import (
	"context"
	"errors"

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

	branchName, _ := flag.GetString("branch")
	org, _ := flag.GetString("org")

	if org == "" {
		return errors.New("no organization set")
	}

	vc, err := getVersionController(flag)
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
		return err
	}

	return nil
}

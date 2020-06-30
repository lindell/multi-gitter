package cmd

import (
	"context"

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
}

func merge(cmd *cobra.Command, args []string) error {
	flag := cmd.Flags()

	branchName, _ := flag.GetString("branch")

	vc, err := getVersionController(flag)
	if err != nil {
		return err
	}

	statuser := multigitter.Merger{
		VersionController: vc,

		FeatureBranch: branchName,
	}

	err = statuser.Merge(context.Background())
	if err != nil {
		return err
	}

	return nil
}

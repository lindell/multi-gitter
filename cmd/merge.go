package cmd

import (
	"context"

	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/spf13/cobra"
)

// MergeCmd merges pull requests
func MergeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "merge",
		Short:   "Merge pull requests.",
		Long:    "Merge pull requests with a specified branch name in an organization and with specified conditions.",
		Args:    cobra.NoArgs,
		PreRunE: logFlagInit,
		RunE:    merge,
	}

	cmd.Flags().StringP("branch", "B", "multi-gitter-branch", "The name of the branch where changes are committed.")
	cmd.Flags().StringSliceP("merge-type", "", []string{"merge", "squash", "rebase"}, "The type of merge that should be done (GitHub). Multiple types can be used as backup strategies if the first one is not allowed.")
	configurePlatform(cmd)
	configureLogging(cmd, "-")

	return cmd
}

func merge(cmd *cobra.Command, args []string) error {
	flag := cmd.Flags()

	branchName, _ := flag.GetString("branch")

	vc, err := getVersionController(flag, true)
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

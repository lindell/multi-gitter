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
	configureMergeType(cmd, false)
	configurePlatform(cmd)
	configureRunPlatform(cmd, false)
	configureLogging(cmd, "-")
	configureConfig(cmd)

	return cmd
}

func merge(cmd *cobra.Command, _ []string) error {
	flag := cmd.Flags()

	branchName, _ := flag.GetString("branch")

	vc, err := getVersionController(flag, true, false)
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

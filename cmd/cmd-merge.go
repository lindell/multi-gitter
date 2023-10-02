package cmd

import (
	"context"

	"github.com/lindell/multi-gitter/cmd/namedflag"
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

	fss := namedflag.New(cmd)
	flags := fss.FlagSet("Merge")

	flags.StringP("branch", "B", "multi-gitter-branch", "The name of the branch where changes are committed.")
	flags.StringSliceP("merge-type", "", []string{"merge", "squash", "rebase"},
		"The type of merge that should be done (GitHub). Multiple types can be used as backup strategies if the first one is not allowed.")
	configurePlatform(cmd, fss.FlagSet("Platform"))
	configureRunPlatform(fss.FlagSet("Platform"), false)
	configureLogging(fss.FlagSet("Logging"), "-")
	configureConfig(fss.FlagSet("Config"))

	namedflag.SetUsageAndHelpFunc(cmd, fss)

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

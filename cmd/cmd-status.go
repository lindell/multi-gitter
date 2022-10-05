package cmd

import (
	"context"
	"os"

	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/spf13/cobra"
)

// StatusCmd gets statuses of pull requests
func StatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Short:   "Get the status of pull requests.",
		Long:    "Get the status of all pull requests with a specified branch name in an organization.",
		Args:    cobra.NoArgs,
		PreRunE: logFlagInit,
		RunE:    status,
	}

	cmd.Flags().StringP("branch", "B", "multi-gitter-branch", "The name of the branch where changes are committed.")
	configurePlatform(cmd)
	configureRunPlatform(cmd, false)
	configureLogging(cmd, "-")
	configureConfig(cmd)
	cmd.Flags().AddFlagSet(outputFlag())

	return cmd
}

func status(cmd *cobra.Command, args []string) error {
	flag := cmd.Flags()

	branchName, _ := flag.GetString("branch")
	strOutput, _ := flag.GetString("output")

	vc, err := getVersionController(flag, true, false)
	if err != nil {
		return err
	}

	output, err := fileOutput(strOutput, os.Stdout)
	if err != nil {
		return err
	}

	statuser := multigitter.Statuser{
		VersionController: vc,

		Output: output,

		FeatureBranch: branchName,
	}

	err = statuser.Statuses(context.Background())
	if err != nil {
		return err
	}

	return nil
}

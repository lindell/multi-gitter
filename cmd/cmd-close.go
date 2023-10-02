package cmd

import (
	"context"

	"github.com/lindell/multi-gitter/cmd/namedflag"
	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/spf13/cobra"
)

// CloseCmd closes pull requests
func CloseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "close",
		Short:   "Close pull requests.",
		Long:    "Close pull requests with a specified branch name in an organization and with specified conditions.",
		Args:    cobra.NoArgs,
		PreRunE: logFlagInit,
		RunE:    closeCMD,
	}

	fss := namedflag.New(cmd)
	flags := fss.FlagSet("Close")

	flags.StringP("branch", "B", "multi-gitter-branch", "The name of the branch where changes are committed.")
	configurePlatform(cmd, fss.FlagSet("Platform"))
	configureRunPlatform(fss.FlagSet("Platform"), false)
	configureLogging(fss.FlagSet("Logging"), "-")
	configureConfig(fss.FlagSet("Config"))

	namedflag.SetUsageAndHelpFunc(cmd, fss)

	return cmd
}

func closeCMD(cmd *cobra.Command, _ []string) error {
	flag := cmd.Flags()

	branchName, _ := flag.GetString("branch")

	vc, err := getVersionController(flag, true, false)
	if err != nil {
		return err
	}

	statuser := multigitter.Closer{
		VersionController: vc,

		FeatureBranch: branchName,
	}

	err = statuser.Close(context.Background())
	if err != nil {
		return err
	}

	return nil
}

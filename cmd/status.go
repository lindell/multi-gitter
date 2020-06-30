package cmd

import (
	"context"

	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/spf13/cobra"
)

// StatusCmd gets statuses of pull requests
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of pull requests.",
	Long:  "Get the status of all pull requests with a specified branch name in an organization.",
	Args:  cobra.NoArgs,
	RunE:  status,
}

func init() {
	StatusCmd.Flags().StringP("branch", "B", "multi-gitter-branch", "The name of the branch where changes are committed.")
}

func status(cmd *cobra.Command, args []string) error {
	flag := cmd.Flags()

	branchName, _ := flag.GetString("branch")

	vc, err := getVersionController(flag)
	if err != nil {
		return err
	}

	statuser := multigitter.Statuser{
		VersionController: vc,

		FeatureBranch: branchName,
	}

	err = statuser.Statuses(context.Background())
	if err != nil {
		return err
	}

	return nil
}

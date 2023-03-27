package cmd

import (
	"context"
	"os"
	"runtime"

	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/spf13/cobra"
)

// ReviewCmd reviews pull requests
func ReviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "review",
		Short:   "Review pull requests.",
		Long:    "Review pull requests with a specified branch name in an organization and with specified conditions.",
		Args:    cobra.NoArgs,
		PreRunE: logFlagInit,
		RunE:    review,
	}

	cmd.Flags().StringP("branch", "B", "multi-gitter-branch", "The name of the branch where changes are committed.")
	cmd.Flags().StringP("comment", "c", "", "Leave a review comment.")
	cmd.Flags().BoolP("all", "a", false, "Review all pull requests in one go instead each individually.")
	cmd.Flags().StringP("batch", "", "", "Review all pull requests without confirmation (--all is not required). Batch accepts one of [approve, reject, comment].")
	_ = cmd.RegisterFlagCompletionFunc("batch", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"approve", "reject", "comment"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Flags().BoolP("no-pager", "", false, "Do not use a pager for reviewing pull request diffs.")
	cmd.Flags().BoolP("include-approved", "", false, "Include pull requests already approved by you.")
	configurePlatform(cmd)
	configureRunPlatform(cmd, false)
	configureLogging(cmd, "-")
	configureConfig(cmd)

	return cmd
}

func review(cmd *cobra.Command, _ []string) error {
	flag := cmd.Flags()

	branchName, _ := flag.GetString("branch")
	comment, _ := flag.GetString("comment")
	all, _ := flag.GetBool("all")
	batch, _ := flag.GetString("batch")
	disablePaging, _ := flag.GetBool("no-pager")
	includeApproved, _ := flag.GetBool("include-approved")

	var batchOperation *multigitter.BatchOperation
	if batch != "" {
		b, err := multigitter.ParseBatchOperation(batch)
		if err != nil {
			return err
		}

		batchOperation = &b
	}

	vc, err := getVersionController(flag, true, false)
	if err != nil {
		return err
	}

	statuser := multigitter.Reviewer{
		VersionController: vc,
		FeatureBranch:     branchName,
		Comment:           comment,
		All:               all,
		BatchOperation:    batchOperation,
		Pager:             getPager(),
		DisablePaging:     disablePaging,
		IncludeApproved:   includeApproved,
	}

	err = statuser.Review(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func getPager() string {
	pager := os.Getenv("PAGER")
	if pager != "" {
		return pager
	}

	switch runtime.GOOS {
	case "windows":
		return "more"
	default:
		return "less"
	}
}

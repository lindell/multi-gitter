package cmd

import (
	"context"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/spf13/cobra"
)

const runHelp = `
This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. If the script finished with a zero exit code, and the script resulted in file changes, a pull request will be created with.

The environment variable REPOSITORY_NAME will be set to the name of the repository currently being executed by the script.
`

// RunCmd is the main command that runs a script for multiple repositories and creates PRs with the changes made
var RunCmd = &cobra.Command{
	Use:   "run [script path]",
	Short: "Clones multiple repostories, run a script in that directory, and creates a PR with those changes.",
	Long:  runHelp,
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

func init() {
	RunCmd.Flags().StringP("branch", "B", "multi-gitter-branch", "The name of the branch where changes are committed.")
	RunCmd.Flags().StringP("pr-title", "t", "", "The title of the PR. Will default to the first line of the commit message if none is set.")
	RunCmd.Flags().StringP("pr-body", "b", "", "The body of the commit message. Will default to everything but the first line of the commit message if none is set.")
	RunCmd.Flags().StringP("commit-message", "m", "", "The commit message. Will default to title + body if none is set.")
	RunCmd.Flags().StringSliceP("reviewers", "r", nil, "The username of the reviewers to be added on the pull request.")
	RunCmd.Flags().IntP("max-reviewers", "R", 0, "If this value is set, reviewers will be randomized")
	RunCmd.Flags().BoolP("dry-run", "d", false, "Run without pushing changes or creating pull requests")
}

func run(cmd *cobra.Command, args []string) error {
	flag := cmd.Flags()

	branchName, _ := flag.GetString("branch")
	prTitle, _ := flag.GetString("pr-title")
	prBody, _ := flag.GetString("pr-body")
	commitMessage, _ := flag.GetString("commit-message")
	reviewers, _ := flag.GetStringSlice("reviewers")
	maxReviewers, _ := flag.GetInt("max-reviewers")
	dryRun, _ := flag.GetBool("dry-run")

	token, err := getToken(flag)
	if err != nil {
		return err
	}

	programPath := flag.Arg(0)

	// Set commit message based on pr title and body or the reverse
	if commitMessage == "" && prTitle == "" {
		return errors.New("pull request title or commit message must be set")
	} else if commitMessage == "" {
		commitMessage = prTitle
		if prBody != "" {
			commitMessage += "\n" + prBody
		}
	} else if prTitle == "" {
		split := strings.SplitN(commitMessage, "\n", 2)
		prTitle = split[0]
		if prBody == "" && len(split) == 2 {
			prBody = split[2]
		}
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return errors.New("could not get the working directory")
	}

	vc, err := getVersionController(flag)
	if err != nil {
		return err
	}

	runner := multigitter.Runner{
		ScriptPath:    path.Join(workingDir, programPath),
		FeatureBranch: branchName,
		Token:         token,

		VersionController: vc,

		CommitMessage:    commitMessage,
		PullRequestTitle: prTitle,
		PullRequestBody:  prBody,
		Reviewers:        reviewers,
		MaxReviewers:     maxReviewers,
		DryRun:           dryRun,
	}

	err = runner.Run(context.Background())
	if err != nil {
		return err
	}

	return nil
}

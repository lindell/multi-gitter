package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lindell/multi-gitter/internal/git"

	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/spf13/cobra"
)

const runHelp = `
This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. If the script finished with a zero exit code, and the script resulted in file changes, a pull request will be created with.

When the script is invoked, these environment variables are set:
- REPOSITORY will be set to the name of the repository currently being executed
- DRY_RUN will be set =true, when running in with the --dry-run flag, otherwise it's absent
`

// RunCmd is the main command that runs a script for multiple repositories and creates PRs with the changes made
func RunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run [script path]",
		Short:   "Clones multiple repositories, run a script in that directory, and creates a PR with those changes.",
		Long:    runHelp,
		Args:    cobra.ExactArgs(1),
		PreRunE: logFlagInit,
		RunE:    run,
	}

	cmd.Flags().StringP("branch", "B", "multi-gitter-branch", "The name of the branch where changes are committed.")
	cmd.Flags().StringP("base-branch", "", "", "The branch which the changes will be based on.")
	cmd.Flags().StringP("pr-title", "t", "", "The title of the PR. Will default to the first line of the commit message if none is set.")
	cmd.Flags().StringP("pr-body", "b", "", "The body of the commit message. Will default to everything but the first line of the commit message if none is set.")
	cmd.Flags().StringP("commit-message", "m", "", "The commit message. Will default to title + body if none is set.")
	cmd.Flags().StringSliceP("reviewers", "r", nil, "The username of the reviewers to be added on the pull request.")
	cmd.Flags().StringSliceP("assignees", "a", nil, "The username of the assignees to be added on the pull request.")
	cmd.Flags().IntP("max-reviewers", "M", 0, "If this value is set, reviewers will be randomized.")
	cmd.Flags().IntP("concurrent", "C", 1, "The maximum number of concurrent runs.")
	cmd.Flags().BoolP("skip-pr", "", false, "Skip pull request and directly push to the branch.")
	cmd.Flags().StringSliceP("skip-repo", "s", nil, "Skip changes on specified repositories, the name is including the owner of repository in the format \"ownerName/repoName\".")
	cmd.Flags().BoolP("interactive", "i", false, "Take manual decision before committing any change. Requires git to be installed.")
	cmd.Flags().BoolP("dry-run", "d", false, "Run without pushing changes or creating pull requests.")
	cmd.Flags().StringP("conflict-strategy", "", "skip", `What should happen if the branch already exist.
Available values:
  skip: Skip making any changes to the existing branch and do not create a new pull request.
  replace: Replace the existing content of the branch by force pushing any new changes, then reuse any existing pull request, or create a new one if none exist.
`)
	cmd.Flags().BoolP("draft", "", false, "Create pull request(s) as draft.")
	_ = cmd.RegisterFlagCompletionFunc("conflict-strategy", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"skip", "replace"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringSliceP("labels", "", nil, "Labels to be added to any created pull request.")
	cmd.Flags().StringP("author-name", "", "", "Name of the committer. If not set, the global git config setting will be used.")
	cmd.Flags().StringP("author-email", "", "", "Email of the committer. If not set, the global git config setting will be used.")
	configureGit(cmd)
	configurePlatform(cmd)
	configureRunPlatform(cmd, true)
	configureLogging(cmd, "-")
	configureConfig(cmd)
	cmd.Flags().AddFlagSet(outputFlag())

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	flag := cmd.Flags()

	branchName, _ := flag.GetString("branch")
	baseBranchName, _ := flag.GetString("base-branch")
	prTitle, _ := flag.GetString("pr-title")
	prBody, _ := flag.GetString("pr-body")
	commitMessage, _ := flag.GetString("commit-message")
	reviewers, _ := flag.GetStringSlice("reviewers")
	maxReviewers, _ := flag.GetInt("max-reviewers")
	concurrent, _ := flag.GetInt("concurrent")
	skipPullRequest, _ := flag.GetBool("skip-pr")
	skipRepository, _ := flag.GetStringSlice("skip-repo")
	interactive, _ := flag.GetBool("interactive")
	dryRun, _ := flag.GetBool("dry-run")
	forkMode, _ := flag.GetBool("fork")
	forkOwner, _ := flag.GetString("fork-owner")
	conflictStrategyStr, _ := flag.GetString("conflict-strategy")
	authorName, _ := flag.GetString("author-name")
	authorEmail, _ := flag.GetString("author-email")
	strOutput, _ := flag.GetString("output")
	assignees, _ := flag.GetStringSlice("assignees")
	draft, _ := flag.GetBool("draft")
	labels, _ := flag.GetStringSlice("labels")

	if concurrent < 1 {
		return errors.New("concurrent runs can't be less than one")
	}

	output, err := fileOutput(strOutput, os.Stdout)
	if err != nil {
		return err
	}

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
			prBody = split[1]
		}
	}

	if skipPullRequest && forkMode {
		return errors.New("--fork and --skip-pr can't be used at the same time")
	}

	if concurrent > 1 && interactive {
		return errors.New("--concurrent and --interactive can't be used at the same time")
	}

	// Parse commit author data
	var commitAuthor *git.CommitAuthor
	if authorName != "" || authorEmail != "" {
		if authorName == "" || authorEmail == "" {
			return errors.New("both author-name and author-email has to be set if the other is set")
		}
		commitAuthor = &git.CommitAuthor{
			Name:  authorName,
			Email: authorEmail,
		}
	}

	vc, err := getVersionController(flag, true, false)
	if err != nil {
		return err
	}

	gitCreator, err := getGitCreator(flag)
	if err != nil {
		return err
	}

	executablePath, arguments, err := parseCommand(flag.Arg(0))
	if err != nil {
		return err
	}

	conflictStrategy, err := multigitter.ParseConflictStrategy(conflictStrategyStr)
	if err != nil {
		return err
	}

	// Set up signal listening to cancel the context and let started runs finish gracefully
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Finishing up ongoing runs. Press CTRL+C again to abort now.")
		cancel()
		<-c
		os.Exit(1)
	}()

	runner := &multigitter.Runner{
		ScriptPath:    executablePath,
		Arguments:     arguments,
		FeatureBranch: branchName,

		Output: output,

		VersionController: vc,

		CommitMessage:    commitMessage,
		PullRequestTitle: prTitle,
		PullRequestBody:  prBody,
		Reviewers:        reviewers,
		MaxReviewers:     maxReviewers,
		Interactive:      interactive,
		DryRun:           dryRun,
		Fork:             forkMode,
		ForkOwner:        forkOwner,
		SkipPullRequest:  skipPullRequest,
		SkipRepository:   skipRepository,
		CommitAuthor:     commitAuthor,
		BaseBranch:       baseBranchName,
		Assignees:        assignees,
		ConflictStrategy: conflictStrategy,
		Draft:            draft,
		Labels:           labels,

		Concurrent: concurrent,

		CreateGit: gitCreator,
	}

	err = runner.Run(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return nil
}

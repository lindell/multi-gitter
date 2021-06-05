package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/lindell/multi-gitter/internal/domain"

	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/spf13/cobra"
)

const runHelp = `
This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. If the script finished with a zero exit code, and the script resulted in file changes, a pull request will be created with.

The environment variable REPOSITORY will be set to the name of the repository currently being executed by the script.
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
	cmd.Flags().IntP("max-reviewers", "M", 0, "If this value is set, reviewers will be randomized")
	cmd.Flags().IntP("concurrent", "C", 1, "The maximum number of concurrent runs")
	cmd.Flags().IntP("fetch-depth", "f", 1, "Limit fetching to the specified number of commits. Set to 0 for no limit")
	cmd.Flags().BoolP("skip-pr", "", false, "Skip pull request and directly push to the branch")
	cmd.Flags().BoolP("dry-run", "d", false, "Run without pushing changes or creating pull requests")
	cmd.Flags().StringP("author-name", "", "", "Name of the committer. If not set, the global git config setting will be used.")
	cmd.Flags().StringP("author-email", "", "", "Email of the committer. If not set, the global git config setting will be used.")
	configurePlatform(cmd)
	configureLogging(cmd, "-")
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
	fetchDepth, _ := flag.GetInt("fetch-depth")
	concurrent, _ := flag.GetInt("concurrent")
	skipPullRequest, _ := flag.GetBool("skip-pr")
	dryRun, _ := flag.GetBool("dry-run")
	authorName, _ := flag.GetString("author-name")
	authorEmail, _ := flag.GetString("author-email")
	strOutput, _ := flag.GetString("output")

	token, err := getToken(flag)
	if err != nil {
		return err
	}

	command := flag.Arg(0)

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
			prBody = split[2]
		}
	}

	// Parse commit author data
	var commitAuthor *domain.CommitAuthor
	if authorName != "" || authorEmail != "" {
		if authorName == "" || authorEmail == "" {
			return errors.New("both author-name and author-email has to be set if the other is set")
		}
		commitAuthor = &domain.CommitAuthor{
			Name:  authorName,
			Email: authorEmail,
		}
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return errors.New("could not get the working directory")
	}

	vc, err := getVersionController(flag, true)
	if err != nil {
		return err
	}

	parsedCommand, err := parseCommandLine(command)
	if err != nil {
		return fmt.Errorf("could not parse command: %s", err)
	}
	executablePath, err := exec.LookPath(parsedCommand[0])
	if err != nil {
		return fmt.Errorf("could not find executable %s", parsedCommand[0])
	}
	// Executable needs to be defined with an absolute path since it will be run within the context of repositories
	if !filepath.IsAbs(executablePath) {
		executablePath = filepath.Join(workingDir, executablePath)
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

	runner := multigitter.Runner{
		ScriptPath:    executablePath,
		Arguments:     parsedCommand[1:],
		FeatureBranch: branchName,
		Token:         token,

		Output: output,

		VersionController: vc,

		CommitMessage:    commitMessage,
		PullRequestTitle: prTitle,
		PullRequestBody:  prBody,
		Reviewers:        reviewers,
		MaxReviewers:     maxReviewers,
		DryRun:           dryRun,
		SkipPullRequest:  skipPullRequest,
		CommitAuthor:     commitAuthor,
		BaseBranch:       baseBranchName,

		FetchDepth: fetchDepth,
		Concurrent: concurrent,
	}

	err = runner.Run(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return nil
}

// https://stackoverflow.com/a/46973603
func parseCommandLine(command string) ([]string, error) {
	type state int

	const (
		stateStart state = iota
		stateQuotes
		stateArg
	)

	var args []string
	currentState := stateStart
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if currentState == stateQuotes {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				currentState = stateStart
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			currentState = stateQuotes
			quote = string(c)
			continue
		}

		if currentState == stateArg {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				currentState = stateStart
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			currentState = stateArg
			current += string(c)
		}
	}

	if currentState == stateQuotes {
		return []string{}, fmt.Errorf("unclosed quote in command line: %s", command)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

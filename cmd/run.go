package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/lindell/multi-gitter/internal/domain"

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
	RunCmd.Flags().IntP("max-reviewers", "M", 0, "If this value is set, reviewers will be randomized")
	RunCmd.Flags().BoolP("dry-run", "d", false, "Run without pushing changes or creating pull requests")
	RunCmd.Flags().StringP("author-name", "", "", "If set, this fields will be used as the name of the committer")
	RunCmd.Flags().StringP("author-email", "", "", "If set, this fields will be used as the email of the committer")
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
	authorName, _ := flag.GetString("author-name")
	authorEmail, _ := flag.GetString("author-email")

	token, err := getToken(flag)
	if err != nil {
		return err
	}

	command := flag.Arg(0)

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

	vc, err := getVersionController(flag)
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
	if !path.IsAbs(executablePath) {
		executablePath = path.Join(workingDir, executablePath)
	}

	runner := multigitter.Runner{
		ScriptPath:    executablePath,
		Arguments:     parsedCommand[1:],
		FeatureBranch: branchName,
		Token:         token,

		VersionController: vc,

		CommitMessage:    commitMessage,
		PullRequestTitle: prTitle,
		PullRequestBody:  prBody,
		Reviewers:        reviewers,
		MaxReviewers:     maxReviewers,
		DryRun:           dryRun,
		CommitAuthor:     commitAuthor,
	}

	err = runner.Run(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return nil
}

// https://stackoverflow.com/a/46973603
func parseCommandLine(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
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
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, fmt.Errorf("unclosed quote in command line: %s", command)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/lindell/multi-gitter/internal/github"
	"github.com/lindell/multi-gitter/internal/multigitter"
)

func main() {
	ghBaseUrl := flag.StringP("gh-base-url", "g", "https://api.github.com/", "Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.")
	token := flag.StringP("token", "T", "", "The GitHub personal access token. Can also be set using the GITHUB_TOKEN enviroment variable.")
	branchName := flag.StringP("branch", "B", "multi-gitter-branch", "The name of the branch where changes are committed.")
	org := flag.StringP("org", "o", "", "The name of the  GitHub organization.")
	prTitle := flag.StringP("pr-title", "t", "", "The title of the PR. Will default to the first line of the commit message if none is set.")
	prBody := flag.StringP("pr-body", "b", "", "The body of the commit message. Will default to everything but the first line of the commit message if none is set.")
	commitMessage := flag.StringP("commit-message", "m", "", "The commit message. Will default to title + body if none is set.")
	reviewers := flag.StringSliceP("reviewers", "r", nil, "The username of the reviewers to be added on the pull request.")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	programPath := flag.Arg(0)

	ghConfig := github.DefaultConfig
	ghConfig.BaseURL = *ghBaseUrl
	if *token != "" {
		ghConfig.Token = *token
	} else if ght := os.Getenv("GITHUB_TOKEN"); ght != "" {
		ghConfig.Token = ght
	} else {
		fmt.Println("Either the --token flag or the GITHUB_TOKEN enviroment variable has to be set.")
		flag.Usage()
		os.Exit(1)
	}

	if *org == "" {
		fmt.Println("No organisation set.")
		flag.Usage()
		os.Exit(1)
	}

	// Set commit message based on pr title and body or the reverse
	if *commitMessage == "" && *prTitle == "" {
		fmt.Println("Pull request title or commit message must be set.")
		flag.Usage()
		os.Exit(1)
	} else if *commitMessage == "" {
		*commitMessage = *prTitle
		if *prBody != "" {
			*commitMessage += "\n" + *prBody
		}
	} else if *prTitle == "" {
		split := strings.SplitN(*commitMessage, "\n", 2)
		*prTitle = split[0]
		if *prBody == "" && len(split) == 2 {
			*prBody = split[2]
		}
	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(workingDir)
	}

	runner := multigitter.Runner{
		ScriptPath:    path.Join(workingDir, programPath),
		FeatureBranch: *branchName,

		RepoGetter: github.OrgRepoGetter{
			Config:       ghConfig,
			Organization: *org,
		},
		PullRequestCreator: github.PullRequestCreator{
			Config: ghConfig,
		},

		CommitMessage:    *commitMessage,
		PullRequestTitle: *prTitle,
		PullRequestBody:  *prBody,
		Reviewers:        *reviewers,
	}

	err = runner.Run()
	if err != nil {
		log.Fatal(err)
	}
}

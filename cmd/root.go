package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/lindell/multi-gitter/internal/github"
	"github.com/lindell/multi-gitter/internal/gitlab"
	"github.com/lindell/multi-gitter/internal/multigitter"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// RootCmd is the root command containing all subcommands
var RootCmd = &cobra.Command{
	Use:               "multi-gitter",
	Short:             "Multi gitter is a tool for making changes into multiple git repositories",
	PersistentPreRunE: persistentPreRun,
}

func init() {
	RootCmd.PersistentFlags().StringP("gh-base-url", "g", "", "Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.")
	RootCmd.PersistentFlags().StringP("token", "T", "", "The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.")
	RootCmd.PersistentFlags().StringP("log-level", "L", "info", "The level of logging that should be made. Available values: trace, debug, info, error")
	RootCmd.PersistentFlags().StringSliceP("org", "o", nil, "The name of a GitHub organization. All repositories in that organization will be used.")
	RootCmd.PersistentFlags().StringSliceP("group", "G", nil, "The name of a GitLab organization. All repositories in that group will be used.")
	RootCmd.PersistentFlags().StringSliceP("user", "u", nil, "The name of a user. All repositories owned by that user will be used.")
	RootCmd.PersistentFlags().StringSliceP("repo", "R", nil, "The name, including owner of a repository in the format \"ownerName/repoName\"")
	RootCmd.PersistentFlags().StringP("platform", "P", "github", "The platform that is used. Available values: github, gitlab")

	RootCmd.AddCommand(RunCmd)
	RootCmd.AddCommand(StatusCmd)
	RootCmd.AddCommand(MergeCmd)
	RootCmd.AddCommand(VersionCmd)

	rand.Seed(time.Now().UTC().UnixNano())
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	strLevel, _ := cmd.Flags().GetString("log-level")
	logLevel, err := log.ParseLevel(strLevel)
	if err != nil {
		return fmt.Errorf("invalid log-level: %s", strLevel)
	}
	log.SetLevel(logLevel)

	return nil
}

func getVersionController(flag *flag.FlagSet) (multigitter.VersionController, error) {
	platform, _ := flag.GetString("platform")
	switch platform {
	default:
		return nil, fmt.Errorf("unknown platform: %s", platform)
	case "github":
		return createGithubClient(flag)
	case "gitlab":
		return createGitlabClient(flag)
	}
}

func createGithubClient(flag *flag.FlagSet) (multigitter.VersionController, error) {
	ghBaseURL, _ := flag.GetString("gh-base-url")
	orgs, _ := flag.GetStringSlice("org")
	users, _ := flag.GetStringSlice("user")
	repos, _ := flag.GetStringSlice("repo")

	if len(orgs) == 0 && len(users) == 0 && len(repos) == 0 {
		return nil, errors.New("no organization or user set")
	}

	token, err := getToken(flag)
	if err != nil {
		return nil, err
	}

	repoRefs := make([]github.RepositoryReference, len(repos))
	for i := range repos {
		repoRefs[i], err = github.ParseRepositoryReference(repos[i])
		if err != nil {
			return nil, err
		}
	}

	vc, err := github.New(token, ghBaseURL, github.RepositoryListing{
		Organizations: orgs,
		Users:         users,
		Repositories:  repoRefs,
	})
	if err != nil {
		return nil, err
	}

	return vc, nil
}

func createGitlabClient(flag *flag.FlagSet) (multigitter.VersionController, error) {
	groups, _ := flag.GetStringSlice("group")

	token, err := getToken(flag)
	if err != nil {
		return nil, err
	}

	vc, err := gitlab.New(token, "", gitlab.RepositoryListing{
		Groups: groups,
	})
	if err != nil {
		return nil, err
	}

	return vc, nil
}

func getToken(flag *flag.FlagSet) (string, error) {
	token, _ := flag.GetString("token")

	if token == "" {
		if ght := os.Getenv("GITHUB_TOKEN"); ght != "" {
			token = ght
		} else if ght := os.Getenv("GITLAB_TOKEN"); ght != "" {
			token = ght
		}
	}

	if token == "" {
		return "", errors.New("either the --token flag or the GITHUB_TOKEN environment variable has to be set")
	}

	return token, nil
}

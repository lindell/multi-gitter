package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/lindell/multi-gitter/internal/github"
	"github.com/lindell/multi-gitter/internal/gitlab"
	"github.com/lindell/multi-gitter/internal/multigitter"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// RootCmd is the root command containing all subcommands
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi-gitter",
		Short: "Multi gitter is a tool for making changes into multiple git repositories",
	}

	cmd.AddCommand(RunCmd())
	cmd.AddCommand(StatusCmd())
	cmd.AddCommand(MergeCmd())
	cmd.AddCommand(CloseCmd())
	cmd.AddCommand(PrintCmd())
	cmd.AddCommand(VersionCmd())

	return cmd
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func platformFlags() *flag.FlagSet {
	flags := flag.NewFlagSet("platform", flag.ExitOnError)

	flags.StringP("gh-base-url", "g", "", "Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.")
	flags.StringP("token", "T", "", "The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.")

	flags.StringSliceP("org", "o", nil, "The name of a GitHub organization. All repositories in that organization will be used.")
	flags.StringSliceP("group", "G", nil, "The name of a GitLab organization. All repositories in that group will be used.")
	flags.StringSliceP("user", "u", nil, "The name of a user. All repositories owned by that user will be used.")
	flags.StringSliceP("repo", "R", nil, "The name, including owner of a GitHub repository in the format \"ownerName/repoName\"")
	flags.StringSliceP("project", "p", nil, "The name, including owner of a GitLab project in the format \"ownerName/repoName\"")
	flags.StringP("platform", "P", "github", "The platform that is used. Available values: github, gitlab")

	return flags
}

func logFlags(logFile string) *flag.FlagSet {
	flags := flag.NewFlagSet("log", flag.ExitOnError)

	flags.StringP("log-level", "L", "info", "The level of logging that should be made. Available values: trace, debug, info, error")
	flags.StringP("log-file", "", logFile, `The file where all logs should be printed to. "-" means stdout`)

	return flags
}

func logFlagInit(cmd *cobra.Command, args []string) error {
	// Parse and set log level
	strLevel, _ := cmd.Flags().GetString("log-level")
	logLevel, err := log.ParseLevel(strLevel)
	if err != nil {
		return fmt.Errorf("invalid log-level: %s", strLevel)
	}
	log.SetLevel(logLevel)

	// Set the output (file)
	strFile, _ := cmd.Flags().GetString("log-file")
	if strFile == "" {
		log.SetOutput(nopWriter{})
	} else if strFile != "-" {
		file, err := os.Create(strFile)
		if err != nil {
			return errors.Wrapf(err, "could not open log-file %s", strFile)
		}
		log.SetOutput(file)
	}

	return nil
}

// OverrideVersionController can be set to force a specific version controller to be used
// This is used to override the version controller with a mock, to be used during testing
var OverrideVersionController multigitter.VersionController = nil

func getVersionController(flag *flag.FlagSet) (multigitter.VersionController, error) {
	if OverrideVersionController != nil {
		return OverrideVersionController, nil
	}

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
	users, _ := flag.GetStringSlice("user")
	projects, _ := flag.GetStringSlice("project")

	token, err := getToken(flag)
	if err != nil {
		return nil, err
	}

	projRefs := make([]gitlab.ProjectReference, len(projects))
	for i := range projects {
		projRefs[i], err = gitlab.ParseProjectReference(projects[i])
		if err != nil {
			return nil, err
		}
	}

	vc, err := gitlab.New(token, "", gitlab.RepositoryListing{
		Groups:   groups,
		Users:    users,
		Projects: projRefs,
	})
	if err != nil {
		return nil, err
	}

	return vc, nil
}

func getToken(flag *flag.FlagSet) (string, error) {
	if OverrideVersionController != nil {
		return "", nil
	}

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

// nopWriter is a writer that does nothing
type nopWriter struct{}

func (nw nopWriter) Write(bb []byte) (int, error) {
	return len(bb), nil
}

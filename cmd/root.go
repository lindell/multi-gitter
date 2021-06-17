package cmd

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/lindell/multi-gitter/internal/domain"
	"github.com/lindell/multi-gitter/internal/git/cmdgit"
	"github.com/lindell/multi-gitter/internal/git/gogit"
	"github.com/lindell/multi-gitter/internal/http"
	internallog "github.com/lindell/multi-gitter/internal/log"
	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/lindell/multi-gitter/internal/scm/gitea"
	"github.com/lindell/multi-gitter/internal/scm/github"
	"github.com/lindell/multi-gitter/internal/scm/gitlab"

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
	cmd.AddCommand(CompletionCmd())
	cmd.AddCommand(VersionCmd())

	return cmd
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func configurePlatform(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringP("base-url", "g", "", "Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.")
	flags.StringP("token", "T", "", "The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN environment variable.")

	flags.StringSliceP("org", "O", nil, "The name of a GitHub organization. All repositories in that organization will be used.")
	flags.StringSliceP("group", "G", nil, "The name of a GitLab organization. All repositories in that group will be used.")
	flags.StringSliceP("user", "U", nil, "The name of a user. All repositories owned by that user will be used.")
	flags.StringSliceP("repo", "R", nil, "The name, including owner of a GitHub repository in the format \"ownerName/repoName\"")
	flags.StringSliceP("project", "P", nil, "The name, including owner of a GitLab project in the format \"ownerName/repoName\"")
	flags.BoolP("include-subgroups", "", false, "Include GitLab subgroups when using the --group flag.")

	flags.StringP("platform", "p", "github", "The platform that is used. Available values: github, gitlab, gitea")
	_ = cmd.RegisterFlagCompletionFunc("platform", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"github", "gitlab", "gitea"}, cobra.ShellCompDirectiveDefault
	})

	// Autocompletion for organizations
	_ = cmd.RegisterFlagCompletionFunc("org", func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		vc, err := getVersionController(cmd.Flags(), false)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		type getOrger interface {
			GetAutocompleteOrganizations(ctx context.Context, _ string) ([]string, error)
		}

		g, ok := vc.(getOrger)
		if !ok {
			return nil, cobra.ShellCompDirectiveError
		}

		orgs, err := g.GetAutocompleteOrganizations(cmd.Root().Context(), toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return orgs, cobra.ShellCompDirectiveDefault
	})

	// Autocompletion for users
	_ = cmd.RegisterFlagCompletionFunc("user", func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		vc, err := getVersionController(cmd.Flags(), false)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		type getUserser interface {
			GetAutocompleteUsers(ctx context.Context, _ string) ([]string, error)
		}

		g, ok := vc.(getUserser)
		if !ok {
			return nil, cobra.ShellCompDirectiveError
		}

		users, err := g.GetAutocompleteUsers(cmd.Root().Context(), toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return users, cobra.ShellCompDirectiveDefault
	})

	// Autocompletion for repositories
	_ = cmd.RegisterFlagCompletionFunc("repo", func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		vc, err := getVersionController(cmd.Flags(), false)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		type getRepositorieser interface {
			GetAutocompleteRepositories(ctx context.Context, _ string) ([]string, error)
		}

		g, ok := vc.(getRepositorieser)
		if !ok {
			return nil, cobra.ShellCompDirectiveError
		}

		users, err := g.GetAutocompleteRepositories(cmd.Root().Context(), toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return users, cobra.ShellCompDirectiveDefault
	})
}

func configureGit(cmd *cobra.Command) {
	cmd.Flags().IntP("fetch-depth", "f", 1, "Limit fetching to the specified number of commits. Set to 0 for no limit")
	cmd.Flags().StringP("git-type", "", "go", `The type of git implementation to use.
Available values:
  go: Uses go-git, a Go native implementation of git. This is compiled with the multi-gitter binary, and no extra dependencies are needed.
  cmd: Calls out to the git command. This requires git to be installed and available with by calling "git".
`)
	_ = cmd.RegisterFlagCompletionFunc("git-type", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"go", "cmd"}, cobra.ShellCompDirectiveDefault
	})
}

func getGitCreator(flag *flag.FlagSet) (func(string) multigitter.Git, error) {
	fetchDepth, _ := flag.GetInt("fetch-depth")
	gitType, _ := flag.GetString("git-type")

	switch gitType {
	case "go":
		return func(path string) multigitter.Git {
			return &gogit.Git{
				Directory:  path,
				FetchDepth: fetchDepth,
			}
		}, nil
	case "cmd":
		return func(path string) multigitter.Git {
			return &cmdgit.Git{
				Directory:  path,
				FetchDepth: fetchDepth,
			}
		}, nil
	}

	return nil, errors.Errorf(`could not parse git type "%s"`, gitType)
}

func configureLogging(cmd *cobra.Command, logFile string) {
	flags := cmd.Flags()

	flags.StringP("log-level", "L", "info", "The level of logging that should be made. Available values: trace, debug, info, error")
	_ = cmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"trace", "debug", "info", "error"}, cobra.ShellCompDirectiveDefault
	})

	flags.StringP("log-format", "", "text", `The formating of the logs. Available values: text, json, json-pretty`)
	_ = cmd.RegisterFlagCompletionFunc("log-format", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"text", "json", "json-pretty"}, cobra.ShellCompDirectiveDefault
	})

	flags.StringP("log-file", "", logFile, `The file where all logs should be printed to. "-" means stdout`)
}

func logFlagInit(cmd *cobra.Command, args []string) error {
	// Parse and set log level
	strLevel, _ := cmd.Flags().GetString("log-level")
	logLevel, err := log.ParseLevel(strLevel)
	if err != nil {
		return fmt.Errorf("invalid log-level: %s", strLevel)
	}
	log.SetLevel(logLevel)

	// Parse and set the log format
	strFormat, _ := cmd.Flags().GetString("log-format")

	var formatter log.Formatter
	switch strFormat {
	case "text":
		formatter = &log.TextFormatter{}
	case "json":
		formatter = &log.JSONFormatter{}
	case "json-pretty":
		formatter = &log.JSONFormatter{
			PrettyPrint: true,
		}
	default:
		return fmt.Errorf(`unknown log-format "%s"`, strFormat)
	}

	// Make sure sensitive data is censored before logging them
	var censorItems []internallog.CensorItem
	if token, err := getToken(cmd.Flags()); err == nil && token != "" {
		censorItems = append(censorItems, internallog.CensorItem{
			Sensitive:   token,
			Replacement: "<TOKEN>",
		})
	}

	log.SetFormatter(&internallog.CensorFormatter{
		CensorItems:         censorItems,
		UnderlyingFormatter: formatter,
	})

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

func outputFlag() *flag.FlagSet {
	flags := flag.NewFlagSet("output", flag.ExitOnError)

	flags.StringP("output", "o", "-", `The file that the output of the script should be outputted to. "-" means stdout`)

	return flags
}

// OverrideVersionController can be set to force a specific version controller to be used
// This is used to override the version controller with a mock, to be used during testing
var OverrideVersionController multigitter.VersionController = nil

// getVersionController gets the complete version controller
// the verifyFlags parameter can be set to false if a complete vc is not required (during autocompletion)
func getVersionController(flag *flag.FlagSet, verifyFlags bool) (multigitter.VersionController, error) {
	if OverrideVersionController != nil {
		return OverrideVersionController, nil
	}

	platform, _ := flag.GetString("platform")
	switch platform {
	default:
		return nil, fmt.Errorf("unknown platform: %s", platform)
	case "github":
		return createGithubClient(flag, verifyFlags)
	case "gitlab":
		return createGitlabClient(flag, verifyFlags)
	case "gitea":
		return createGiteaClient(flag, verifyFlags)
	}
}

func createGithubClient(flag *flag.FlagSet, verifyFlags bool) (multigitter.VersionController, error) {
	gitBaseURL, _ := flag.GetString("base-url")
	orgs, _ := flag.GetStringSlice("org")
	users, _ := flag.GetStringSlice("user")
	repos, _ := flag.GetStringSlice("repo")
	forkMode, _ := flag.GetBool("fork")

	if verifyFlags && len(orgs) == 0 && len(users) == 0 && len(repos) == 0 {
		return nil, errors.New("no organization, user or repo set")
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

	mergeTypes, err := getMergeTypes(flag)
	if err != nil {
		return nil, err
	}

	vc, err := github.New(token, gitBaseURL, http.NewLoggingRoundTripper, github.RepositoryListing{
		Organizations: orgs,
		Users:         users,
		Repositories:  repoRefs,
	}, mergeTypes, forkMode)
	if err != nil {
		return nil, err
	}

	return vc, nil
}

func createGitlabClient(flag *flag.FlagSet, verifyFlags bool) (multigitter.VersionController, error) {
	gitBaseURL, _ := flag.GetString("base-url")
	groups, _ := flag.GetStringSlice("group")
	users, _ := flag.GetStringSlice("user")
	projects, _ := flag.GetStringSlice("project")
	includeSubgroups, _ := flag.GetBool("include-subgroups")

	if verifyFlags && len(groups) == 0 && len(users) == 0 && len(projects) == 0 {
		return nil, errors.New("no group user or project set")
	}

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

	vc, err := gitlab.New(token, gitBaseURL, gitlab.RepositoryListing{
		Groups:   groups,
		Users:    users,
		Projects: projRefs,
	}, gitlab.Config{
		IncludeSubgroups: includeSubgroups,
	})
	if err != nil {
		return nil, err
	}

	return vc, nil
}

func createGiteaClient(flag *flag.FlagSet, verifyFlags bool) (multigitter.VersionController, error) {
	giteaBaseURL, _ := flag.GetString("base-url")
	orgs, _ := flag.GetStringSlice("org")
	users, _ := flag.GetStringSlice("user")
	repos, _ := flag.GetStringSlice("repo")

	if verifyFlags && len(orgs) == 0 && len(users) == 0 && len(repos) == 0 {
		return nil, errors.New("no organization, user or repository set")
	}

	if giteaBaseURL == "" {
		return nil, errors.New("no base-url set")
	}

	token, err := getToken(flag)
	if err != nil {
		return nil, err
	}

	repoRefs := make([]gitea.RepositoryReference, len(repos))
	for i := range repos {
		repoRefs[i], err = gitea.ParseRepositoryReference(repos[i])
		if err != nil {
			return nil, err
		}
	}

	mergeTypes, err := getMergeTypes(flag)
	if err != nil {
		return nil, err
	}

	vc, err := gitea.New(token, giteaBaseURL, gitea.RepositoryListing{
		Organizations: orgs,
		Users:         users,
		Repositories:  repoRefs,
	}, mergeTypes)
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
		} else if ght := os.Getenv("GITEA_TOKEN"); ght != "" {
			token = ght
		}
	}

	if token == "" {
		return "", errors.New("either the --token flag or the GITHUB_TOKEN environment variable has to be set")
	}

	return token, nil
}

func getMergeTypes(flag *flag.FlagSet) ([]domain.MergeType, error) {
	mergeTypeStrs, _ := flag.GetStringSlice("merge-type") // Only used for the merge command

	// Convert all defined merge types (if any)
	var err error
	mergeTypes := make([]domain.MergeType, len(mergeTypeStrs))
	for i, mt := range mergeTypeStrs {
		mergeTypes[i], err = domain.ParseMergeType(mt)
		if err != nil {
			return nil, err
		}
	}

	return mergeTypes, nil
}

// nopWriter is a writer that does nothing
type nopWriter struct{}

func (nw nopWriter) Write(bb []byte) (int, error) {
	return len(bb), nil
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func fileOutput(value string, std io.Writer) (io.WriteCloser, error) {
	if value != "-" {
		file, err := os.Create(value)
		if err != nil {
			return nil, errors.Wrapf(err, "could not open file %s", value)
		}
		return file, nil
	}
	return nopCloser{std}, nil
}

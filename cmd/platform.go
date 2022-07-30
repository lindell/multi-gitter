package cmd

import (
	"context"
	"fmt"

	"github.com/lindell/multi-gitter/internal/http"
	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/lindell/multi-gitter/internal/scm/bitbucketserver"
	"github.com/lindell/multi-gitter/internal/scm/gitea"
	"github.com/lindell/multi-gitter/internal/scm/github"
	"github.com/lindell/multi-gitter/internal/scm/gitlab"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func configurePlatform(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringP("base-url", "g", "", "Base URL of the GitHub API, needs to be changed if GitHub enterprise is used. Or the url to a self-hosted GitLab instance.")
	flags.BoolP("insecure", "", false, "Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.")
	flags.StringP("username", "u", "", "The Bitbucket server username.")
	flags.StringP("token", "T", "", "The GitHub/GitLab personal access token. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable.")

	flags.StringSliceP("org", "O", nil, "The name of a GitHub organization. All repositories in that organization will be used.")
	flags.StringSliceP("group", "G", nil, "The name of a GitLab organization. All repositories in that group will be used.")
	flags.StringSliceP("user", "U", nil, "The name of a user. All repositories owned by that user will be used.")
	flags.StringSliceP("repo", "R", nil, "The name, including owner of a GitHub repository in the format \"ownerName/repoName\".")
	flags.StringSliceP("project", "P", nil, "The name, including owner of a GitLab project in the format \"ownerName/repoName\".")
	flags.BoolP("include-subgroups", "", false, "Include GitLab subgroups when using the --group flag.")
	flags.BoolP("ssh-auth", "", false, `Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.`)

	// This is only used by PrintCmd to mark readOnly mode for version control platform
	flags.Bool("readOnly", false, "If set, This is running in readonly will be read-only.")
	_ = flags.MarkHidden("readOnly")

	flags.StringP("platform", "p", "github", "The platform that is used. Available values: github, gitlab, gitea, bitbucket_server.")
	_ = cmd.RegisterFlagCompletionFunc("platform", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"github", "gitlab", "gitea", "bitbucket_server"}, cobra.ShellCompDirectiveDefault
	})

	// Autocompletion for organizations
	versionControllerCompletion(cmd, "org", func(vc multigitter.VersionController, toComplete string) ([]string, error) {
		g, ok := vc.(interface {
			GetAutocompleteOrganizations(ctx context.Context, _ string) ([]string, error)
		})
		if !ok {
			return nil, nil
		}
		return g.GetAutocompleteOrganizations(cmd.Root().Context(), toComplete)
	})

	// Autocompletion for users
	versionControllerCompletion(cmd, "user", func(vc multigitter.VersionController, toComplete string) ([]string, error) {
		g, ok := vc.(interface {
			GetAutocompleteUsers(ctx context.Context, _ string) ([]string, error)
		})
		if !ok {
			return nil, nil
		}
		return g.GetAutocompleteUsers(cmd.Root().Context(), toComplete)
	})

	// Autocompletion for repositories
	versionControllerCompletion(cmd, "repo", func(vc multigitter.VersionController, toComplete string) ([]string, error) {
		g, ok := vc.(interface {
			GetAutocompleteRepositories(ctx context.Context, _ string) ([]string, error)
		})
		if !ok {
			return nil, nil
		}
		return g.GetAutocompleteRepositories(cmd.Root().Context(), toComplete)
	})
}

// configureRunPlatform defines platform flags that are relevant for commands that either make changes, or handling changes made
func configureRunPlatform(cmd *cobra.Command, prCreating bool) {
	flags := cmd.Flags()

	forkDesc := "Fork the repository instead of creating a new branch on the same owner."
	if !prCreating {
		forkDesc = "Use pull requests made from forks instead of from the same repository."
	}
	flags.BoolP("fork", "", false, forkDesc)

	forkOwnerDesc := "If set, make the fork to the defined value. Default behavior is for the fork to be on the logged in user."
	if !prCreating {
		forkOwnerDesc = "If set, use forks from the defined value instead of the logged in user."
	}

	flags.StringP("fork-owner", "", "", forkOwnerDesc)
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
	case "github":
		return createGithubClient(flag, verifyFlags)
	case "gitlab":
		return createGitlabClient(flag, verifyFlags)
	case "gitea":
		return createGiteaClient(flag, verifyFlags)
	case "bitbucket_server":
		return createBitbucketServerClient(flag, verifyFlags)
	default:
		return nil, fmt.Errorf("unknown platform: %s", platform)
	}
}

func createGithubClient(flag *flag.FlagSet, verifyFlags bool) (multigitter.VersionController, error) {
	gitBaseURL, _ := flag.GetString("base-url")
	orgs, _ := flag.GetStringSlice("org")
	users, _ := flag.GetStringSlice("user")
	repos, _ := flag.GetStringSlice("repo")
	forkMode, _ := flag.GetBool("fork")
	forkOwner, _ := flag.GetString("fork-owner")
	sshAuth, _ := flag.GetBool("ssh-auth")
	readOnly, _ := flag.GetBool("readOnly")

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
	}, mergeTypes, forkMode, forkOwner, sshAuth, readOnly)
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
	sshAuth, _ := flag.GetBool("ssh-auth")

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
		SSHAuth:          sshAuth,
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
	sshAuth, _ := flag.GetBool("ssh-auth")

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
	}, mergeTypes, sshAuth)
	if err != nil {
		return nil, err
	}

	return vc, nil
}

func createBitbucketServerClient(flag *flag.FlagSet, verifyFlags bool) (multigitter.VersionController, error) {
	bitbucketServerBaseURL, _ := flag.GetString("base-url")
	projects, _ := flag.GetStringSlice("org")
	users, _ := flag.GetStringSlice("user")
	repos, _ := flag.GetStringSlice("repo")
	username, _ := flag.GetString("username")
	insecure, _ := flag.GetBool("insecure")
	sshAuth, _ := flag.GetBool("ssh-auth")

	if verifyFlags && len(projects) == 0 && len(users) == 0 && len(repos) == 0 {
		return nil, errors.New("no organization, user or repository set")
	}

	if bitbucketServerBaseURL == "" {
		return nil, errors.New("no base-url set for bitbucket server")
	}

	if username == "" {
		return nil, errors.New("no username set")
	}

	token, err := getToken(flag)
	if err != nil {
		return nil, err
	}

	repoRefs := make([]bitbucketserver.RepositoryReference, len(repos))
	for i := range repos {
		repoRefs[i], err = bitbucketserver.ParseRepositoryReference(repos[i])
		if err != nil {
			return nil, err
		}
	}

	vc, err := bitbucketserver.New(username, token, bitbucketServerBaseURL, insecure, sshAuth, http.NewLoggingRoundTripper, bitbucketserver.RepositoryListing{
		Projects:     projects,
		Users:        users,
		Repositories: repoRefs,
	})
	if err != nil {
		return nil, err
	}

	return vc, nil
}

// versionControllerCompletion is a helper function to allow for easier implementation of Cobra autocompletions that depend on a version controller
func versionControllerCompletion(cmd *cobra.Command, flagName string, fn func(vc multigitter.VersionController, toComplete string) ([]string, error)) {
	_ = cmd.RegisterFlagCompletionFunc(flagName, func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Make sure config files are loaded
		_ = initializeConfig(cmd)

		vc, err := getVersionController(cmd.Flags(), false)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		strs, err := fn(vc, toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return strs, cobra.ShellCompDirectiveNoFileComp
	})
}

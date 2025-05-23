package cmd

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/lindell/multi-gitter/internal/http"
	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/lindell/multi-gitter/internal/scm/bitbucketcloud"
	"github.com/lindell/multi-gitter/internal/scm/bitbucketserver"
	"github.com/lindell/multi-gitter/internal/scm/gerrit"
	"github.com/lindell/multi-gitter/internal/scm/gitea"
	"github.com/lindell/multi-gitter/internal/scm/github"
	"github.com/lindell/multi-gitter/internal/scm/gitlab"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func configurePlatform(cmd *cobra.Command) {
	flags := cmd.Flags()

	flags.StringP("base-url", "g", "", "Base URL of the target platform, needs to be changed for GitHub enterprise, a self-hosted GitLab instance, Gitea or BitBucket, Gerrit.")
	flags.BoolP("insecure", "", false, "Insecure controls whether a client verifies the server certificate chain and host name. Used only for Bitbucket server.")
	flags.StringP("username", "u", "", "The Bitbucket server username.")
	flags.StringP("token", "T", "", "The personal access token for the targeting platform. Can also be set using the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN/BITBUCKET_CLOUD_APP_PASSWORD/BITBUCKET_CLOUD_WORKSPACE_TOKEN/GERRIT_TOKEN environment variable.")
	flags.StringP("auth-type", "", "app-password", `The authentication type. Used only for Bitbucket cloud. Available values: app-password, workspace-token.`)
	_ = cmd.RegisterFlagCompletionFunc("auth-type", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"app-password", "workspace-token"}, cobra.ShellCompDirectiveNoFileComp
	})

	flags.StringSliceP("org", "O", nil, "The name of a GitHub organization. All repositories in that organization will be used.")
	flags.StringSliceP("group", "G", nil, "The name of a GitLab organization. All repositories in that group will be used.")
	flags.StringSliceP("user", "U", nil, "The name of a user. All repositories owned by that user will be used.")
	flags.StringSliceP("repo", "R", nil, "The name, including owner of a GitHub repository in the format \"ownerName/repoName\".")
	flags.StringP("repo-search", "", "", "Use a repository search to find repositories to target (GitHub only). Forks are NOT included by default, use `fork:true` to include them. See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories.")
	flags.StringP("code-search", "", "", "Use a code search to find a set of repositories to target (GitHub only). Repeated results from a given repository will be ignored, forks are NOT included by default (use `fork:true` to include them). See the GitHub documentation for full syntax: https://docs.github.com/en/search-github/searching-on-github/searching-code.")
	flags.StringSliceP("topic", "", nil, "The topic of a GitHub/GitLab/Gitea repository. All repositories having at least one matching topic are targeted.")
	flags.StringSliceP("project", "P", nil, "The name, including owner of a GitLab project in the format \"ownerName/repoName\".")
	flags.BoolP("include-subgroups", "", false, "Include GitLab subgroups when using the --group flag.")
	flags.BoolP("ssh-auth", "", false, `Use SSH cloning URL instead of HTTPS + token. This requires that a setup with ssh keys that have access to all repos and that the server is already in known_hosts.`)
	flags.BoolP("skip-forks", "", false, `Skip repositories which are forks.`)

	flags.StringP("platform", "p", "github", "The platform that is used. Available values: github, gitlab, gitea, bitbucket_server, bitbucket_cloud, gerrit. Note: bitbucket_cloud is in Beta")
	_ = cmd.RegisterFlagCompletionFunc("platform", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"github", "gitlab", "gitea", "bitbucket_server", "bitbucket_cloud"}, cobra.ShellCompDirectiveDefault
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
var OverrideVersionController multigitter.VersionController

// getVersionController gets the complete version controller
// the verifyFlags parameter can be set to false if a complete vc is not required (during autocompletion)
func getVersionController(flag *flag.FlagSet, verifyFlags bool, readOnly bool) (multigitter.VersionController, error) {
	if OverrideVersionController != nil {
		return OverrideVersionController, nil
	}

	platform, _ := flag.GetString("platform")
	switch platform {
	case "github":
		return createGithubClient(flag, verifyFlags, readOnly)
	case "gitlab":
		return createGitlabClient(flag, verifyFlags)
	case "gitea":
		return createGiteaClient(flag, verifyFlags)
	case "bitbucket_server":
		return createBitbucketServerClient(flag, verifyFlags)
	case "bitbucket_cloud":
		return createBitbucketCloudClient(flag, verifyFlags)
	case "gerrit":
		return createGerritClient(flag, verifyFlags)
	default:
		return nil, fmt.Errorf("unknown platform: %s", platform)
	}
}

func createGithubClient(flag *flag.FlagSet, verifyFlags bool, readOnly bool) (multigitter.VersionController, error) {
	gitBaseURL, _ := flag.GetString("base-url")
	orgs, _ := flag.GetStringSlice("org")
	users, _ := flag.GetStringSlice("user")
	repos, _ := flag.GetStringSlice("repo")
	repoSearch, _ := flag.GetString("repo-search")
	codeSearch, _ := flag.GetString("code-search")
	topics, _ := flag.GetStringSlice("topic")
	forkMode, _ := flag.GetBool("fork")
	forkOwner, _ := flag.GetString("fork-owner")
	sshAuth, _ := flag.GetBool("ssh-auth")
	skipForks, _ := flag.GetBool("skip-forks")

	if verifyFlags && len(orgs) == 0 && len(users) == 0 && len(repos) == 0 && repoSearch == "" && codeSearch == "" {
		return nil, errors.New("no organization, user, repo, repo-search or code-search set")
	}

	token, err := getToken(flag)
	if err != nil {
		return nil, err
	}

	// Permissions returned from GitHub does not represent reality for some token types,
	// see https://github.com/lindell/multi-gitter/issues/224 for more information.
	// In those cases, we don't check permissions, and let errors occur if
	// repositories are inaccessible.
	checkPermissions := true
	if strings.HasPrefix(token, "ghs_") {
		checkPermissions = false
	}

	repoRefs := make([]github.RepositoryReference, len(repos))
	for i := range repos {
		repoRefs[i], err = github.ParseRepositoryReference(repos[i])
		if err != nil {
			return nil, err
		}
		if slices.Contains(orgs, repoRefs[i].OwnerName) {
			log.Warnf("Repository %s and organization %s are both set. This is likely a mistake", repoRefs[i].String(), repoRefs[i].OwnerName)
		}
		if slices.Contains(users, repoRefs[i].OwnerName) {
			log.Warnf("Repository %s and user %s are both set. This is likely a mistake", repoRefs[i].String(), repoRefs[i].OwnerName)
		}
	}

	mergeTypes, err := getMergeTypes(flag)
	if err != nil {
		return nil, err
	}

	vc, err := github.New(github.Config{
		Token:               token,
		BaseURL:             gitBaseURL,
		TransportMiddleware: http.NewLoggingRoundTripper,
		RepoListing: github.RepositoryListing{
			CodeSearch:       codeSearch,
			Organizations:    orgs,
			Users:            users,
			Repositories:     repoRefs,
			RepositorySearch: repoSearch,
			Topics:           topics,
			SkipForks:        skipForks,
		},
		MergeTypes:       mergeTypes,
		ForkMode:         forkMode,
		ForkOwner:        forkOwner,
		SSHAuth:          sshAuth,
		ReadOnly:         readOnly,
		CheckPermissions: checkPermissions,
	})
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
	topics, _ := flag.GetStringSlice("topic")
	includeSubgroups, _ := flag.GetBool("include-subgroups")
	sshAuth, _ := flag.GetBool("ssh-auth")
	skipForks, _ := flag.GetBool("skip-forks")

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
		if slices.Contains(groups, projRefs[i].OwnerName) {
			log.Warnf("Repository %s and group %s are both set. This is likely a mistake", projRefs[i].String(), projRefs[i].OwnerName)
		}
		if slices.Contains(users, projRefs[i].OwnerName) {
			log.Warnf("Repository %s and user %s are both set. This is likely a mistake", projRefs[i].String(), projRefs[i].OwnerName)
		}
	}

	vc, err := gitlab.New(token, gitBaseURL, gitlab.RepositoryListing{
		Groups:    groups,
		Users:     users,
		Projects:  projRefs,
		Topics:    topics,
		SkipForks: skipForks,
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
	topics, _ := flag.GetStringSlice("topic")
	sshAuth, _ := flag.GetBool("ssh-auth")
	skipForks, _ := flag.GetBool("skip-forks")

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
		if slices.Contains(orgs, repoRefs[i].OwnerName) {
			log.Warnf("Repository %s and organization %s are both set. This is likely a mistake", repoRefs[i].String(), repoRefs[i].OwnerName)
		}
		if slices.Contains(users, repoRefs[i].OwnerName) {
			log.Warnf("Repository %s and user %s are both set. This is likely a mistake", repoRefs[i].String(), repoRefs[i].OwnerName)
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
		Topics:        topics,
		SkipForks:     skipForks,
	}, mergeTypes, sshAuth)
	if err != nil {
		return nil, err
	}

	return vc, nil
}

func createBitbucketCloudClient(flag *flag.FlagSet, verifyFlags bool) (multigitter.VersionController, error) {
	workspaces, _ := flag.GetStringSlice("org")
	users, _ := flag.GetStringSlice("user")
	repos, _ := flag.GetStringSlice("repo")
	username, _ := flag.GetString("username")
	sshAuth, _ := flag.GetBool("ssh-auth")
	fork, _ := flag.GetBool("fork")
	newOwner, _ := flag.GetString("fork-owner")
	authTypeStr, _ := flag.GetString("auth-type")

	if verifyFlags && len(workspaces) == 0 && len(users) == 0 && len(repos) == 0 {
		return nil, errors.New("no workspace, user or repository set")
	}

	if username == "" {
		return nil, errors.New("no username set")
	}

	token, err := getToken(flag)
	if err != nil {
		return nil, err
	}

	authType, err := bitbucketcloud.ParseAuthType(authTypeStr)
	if err != nil {
		return nil, err
	}

	vc, err := bitbucketcloud.New(username, token, repos, workspaces, users, fork, sshAuth, newOwner, authType)
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

func createGerritClient(flag *flag.FlagSet, verifyFlags bool) (multigitter.VersionController, error) {
	username, _ := flag.GetString("username")
	if username == "" {
		return nil, errors.New("no username set")
	}

	gerritBaseURL, _ := flag.GetString("base-url")
	if gerritBaseURL == "" {
		return nil, errors.New("no base-url set")
	}

	repoSearch, _ := flag.GetString("repo-search")

	token, err := getToken(flag)
	if err != nil {
		return nil, err
	}

	vc, err := gerrit.New(username, token, gerritBaseURL, repoSearch)
	return vc, err
}

// versionControllerCompletion is a helper function to allow for easier implementation of Cobra autocompletions that depend on a version controller
func versionControllerCompletion(cmd *cobra.Command, flagName string, fn func(vc multigitter.VersionController, toComplete string) ([]string, error)) {
	_ = cmd.RegisterFlagCompletionFunc(flagName, func(cmd *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Make sure config files are loaded
		_ = initializeConfig(cmd)

		vc, err := getVersionController(cmd.Flags(), false, false)
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

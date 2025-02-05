package cmd

import (
	"github.com/lindell/multi-gitter/internal/git"
	"github.com/lindell/multi-gitter/internal/git/cmdgit"
	"github.com/lindell/multi-gitter/internal/git/gogit"
	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func configureGit(cmd *cobra.Command) {
	cmd.Flags().IntP("fetch-depth", "f", 1, "Limit fetching to the specified number of commits. Set to 0 for no limit.")
	cmd.Flags().StringP("git-type", "", "go", `The type of git implementation to use.
Available values:
  go: Uses go-git, a Go native implementation of git. This is compiled with the multi-gitter binary, and no extra dependencies are needed.
  cmd: Calls out to the git command. This requires git to be installed and available with by calling "git".
`)
	_ = cmd.RegisterFlagCompletionFunc("git-type", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"go", "cmd"}, cobra.ShellCompDirectiveDefault
	})
}

func getGitCreator(flag *flag.FlagSet) (func(string, scm.Repository) multigitter.Git, error) {
	fetchDepth, _ := flag.GetInt("fetch-depth")
	gitType, _ := flag.GetString("git-type")

	switch gitType {
	case "go":
		return func(path string, repo scm.Repository) multigitter.Git {
			return &gogit.Git{
				Directory:   path,
				FetchDepth:  fetchDepth,
				Credentials: getRepoCredentials(repo),
			}
		}, nil
	case "cmd":
		return func(path string, repo scm.Repository) multigitter.Git {
			return &cmdgit.Git{
				Directory:   path,
				FetchDepth:  fetchDepth,
				Credentials: getRepoCredentials(repo),
			}
		}, nil
	}

	return nil, errors.Errorf(`could not parse git type "%s"`, gitType)
}

func getRepoCredentials(repo scm.Repository) *git.Credentials {
	type credentialser interface {
		Credentials() *git.Credentials
	}

	if repo, ok := repo.(credentialser); ok {
		return repo.Credentials()
	}

	return nil
}

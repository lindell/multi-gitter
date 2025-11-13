package cmd

import (
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/lindell/multi-gitter/internal/multigitter"
)

// configureRepoFilters adds the repository filtering flags to a command
func configureRepoFilters(cmd *cobra.Command) {
	cmd.Flags().StringP("repo-include", "", "", "Include repositories that match with a given Regular Expression")
	cmd.Flags().StringP("repo-exclude", "", "", "Exclude repositories that match with a given Regular Expression")
	cmd.Flags().StringSliceP("skip-repo", "s", nil, "Skip specified repositories, the name is including the owner of repository in the format \"ownerName/repoName\".")
}

// parseRepoFilters parses flags into a multigitter.RepoFilters struct
func parseRepoFilters(flag *pflag.FlagSet) (multigitter.RepoFilters, error) {
	repoInclude, _ := flag.GetString("repo-include")
	repoExclude, _ := flag.GetString("repo-exclude")
	skipRepository, _ := flag.GetStringSlice("skip-repo")

	var regExIncludeRepository *regexp.Regexp
	var regExExcludeRepository *regexp.Regexp
	if repoInclude != "" {
		compiled, err := regexp.Compile(repoInclude)
		if err != nil {
			return multigitter.RepoFilters{}, errors.WithMessage(err, "could not parse repo-include")
		}
		regExIncludeRepository = compiled
	}
	if repoExclude != "" {
		compiled, err := regexp.Compile(repoExclude)
		if err != nil {
			return multigitter.RepoFilters{}, errors.WithMessage(err, "could not parse repo-exclude")
		}
		regExExcludeRepository = compiled
	}

	return multigitter.RepoFilters{
		SkipRepository:         skipRepository,
		RegExIncludeRepository: regExIncludeRepository,
		RegExExcludeRepository: regExExcludeRepository,
	}, nil
}

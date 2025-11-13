package multigitter

import (
	"regexp"

	"github.com/lindell/multi-gitter/internal/scm"
	log "github.com/sirupsen/logrus"
)

// RepoFilters contains repository filtering options shared across commands
type RepoFilters struct {
	// SkipRepository is a list of repositories (owner/repo) that should be skipped
	SkipRepository []string

	// RegExIncludeRepository, when set, only repositories matching it will be included
	RegExIncludeRepository *regexp.Regexp

	// RegExExcludeRepository, when set, repositories matching it will be excluded
	RegExExcludeRepository *regexp.Regexp
}

// Determines if Repository should be excluded based on provided Regular Expression
func excludeRepositoryFilter(repoName string, regExp *regexp.Regexp) bool {
	if regExp == nil {
		return false
	}
	return regExp.MatchString(repoName)
}

// Determines if Repository should be included based on provided Regular Expression
func matchesRepositoryFilter(repoName string, regExp *regexp.Regexp) bool {
	if regExp == nil {
		return true
	}
	return regExp.MatchString(repoName)
}

func filterRepositories(repos []scm.Repository, filters RepoFilters) []scm.Repository {
	skipReposMap := map[string]struct{}{}
	for _, skipRepo := range filters.SkipRepository {
		skipReposMap[skipRepo] = struct{}{}
	}

	filteredRepos := make([]scm.Repository, 0, len(repos))
	for _, r := range repos {
		if _, shouldSkip := skipReposMap[r.FullName()]; shouldSkip {
			log.Infof("Skipping %s since it is in exclusion list", r.FullName())
		} else if !matchesRepositoryFilter(r.FullName(), filters.RegExIncludeRepository) {
			log.Infof("Skipping %s since it does not match the inclusion regexp", r.FullName())
		} else if excludeRepositoryFilter(r.FullName(), filters.RegExExcludeRepository) {
			log.Infof("Skipping %s since it match the exclusion regexp", r.FullName())
		} else {
			filteredRepos = append(filteredRepos, r)
		}
	}
	return filteredRepos
}

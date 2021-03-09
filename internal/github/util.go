package github

import (
	"github.com/google/go-github/v33/github"
	"github.com/lindell/multi-gitter/internal/domain"
)

// maps merge types to what they are called in the github api
var mergeTypeGhName = map[domain.MergeType]string{
	domain.MergeTypeMerge:  "merge",
	domain.MergeTypeRebase: "rebase",
	domain.MergeTypeSquash: "squash",
}

// repoMergeTypes returns a list of all allowed merge types
func repoMergeTypes(repo *github.Repository) []domain.MergeType {
	ret := []domain.MergeType{}
	if repo.GetAllowMergeCommit() {
		ret = append(ret, domain.MergeTypeMerge)
	}
	if repo.GetAllowRebaseMerge() {
		ret = append(ret, domain.MergeTypeRebase)
	}
	if repo.GetAllowSquashMerge() {
		ret = append(ret, domain.MergeTypeSquash)
	}
	return ret
}

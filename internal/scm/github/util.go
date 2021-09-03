package github

import (
	"github.com/google/go-github/v38/github"
	"github.com/lindell/multi-gitter/internal/pullrequest"
)

// maps merge types to what they are called in the github api
var mergeTypeGhName = map[pullrequest.MergeType]string{
	pullrequest.MergeTypeMerge:  "merge",
	pullrequest.MergeTypeRebase: "rebase",
	pullrequest.MergeTypeSquash: "squash",
}

// repoMergeTypes returns a list of all allowed merge types
func repoMergeTypes(repo *github.Repository) []pullrequest.MergeType {
	ret := []pullrequest.MergeType{}
	if repo.GetAllowMergeCommit() {
		ret = append(ret, pullrequest.MergeTypeMerge)
	}
	if repo.GetAllowRebaseMerge() {
		ret = append(ret, pullrequest.MergeTypeRebase)
	}
	if repo.GetAllowSquashMerge() {
		ret = append(ret, pullrequest.MergeTypeSquash)
	}
	return ret
}

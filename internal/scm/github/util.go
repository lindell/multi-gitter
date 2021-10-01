package github

import (
	"github.com/google/go-github/v39/github"
	"github.com/lindell/multi-gitter/internal/git"
)

// maps merge types to what they are called in the github api
var mergeTypeGhName = map[git.MergeType]string{
	git.MergeTypeMerge:  "merge",
	git.MergeTypeRebase: "rebase",
	git.MergeTypeSquash: "squash",
}

// repoMergeTypes returns a list of all allowed merge types
func repoMergeTypes(repo *github.Repository) []git.MergeType {
	ret := []git.MergeType{}
	if repo.GetAllowMergeCommit() {
		ret = append(ret, git.MergeTypeMerge)
	}
	if repo.GetAllowRebaseMerge() {
		ret = append(ret, git.MergeTypeRebase)
	}
	if repo.GetAllowSquashMerge() {
		ret = append(ret, git.MergeTypeSquash)
	}
	return ret
}

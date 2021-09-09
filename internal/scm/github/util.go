package github

import (
	"github.com/google/go-github/v38/github"
	"github.com/lindell/multi-gitter/internal/git"
)

// maps merge types to what they are called in the github api
var mergeTypeGhName = map[git.MergeType]string{
	git.PullRequestMergeTypeMerge:  "merge",
	git.PullRequestMergeTypeRebase: "rebase",
	git.PullRequestMergeTypeSquash: "squash",
}

// repoMergeTypes returns a list of all allowed merge types
func repoMergeTypes(repo *github.Repository) []git.MergeType {
	ret := []git.MergeType{}
	if repo.GetAllowMergeCommit() {
		ret = append(ret, git.PullRequestMergeTypeMerge)
	}
	if repo.GetAllowRebaseMerge() {
		ret = append(ret, git.PullRequestMergeTypeRebase)
	}
	if repo.GetAllowSquashMerge() {
		ret = append(ret, git.PullRequestMergeTypeSquash)
	}
	return ret
}

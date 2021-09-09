package gitea

import (
	"code.gitea.io/sdk/gitea"
	"github.com/lindell/multi-gitter/internal/git"
)

// maps merge types to what they are called in the gitea api
var mergeTypeGiteaName = map[git.MergeType]gitea.MergeStyle{
	git.PullRequestMergeTypeMerge:  gitea.MergeStyleMerge,
	git.PullRequestMergeTypeRebase: gitea.MergeStyleRebase,
	git.PullRequestMergeTypeSquash: gitea.MergeStyleSquash,
}

// repoMergeTypes returns a list of all allowed merge types
func repoMergeTypes(repo *gitea.Repository) []git.MergeType {
	ret := []git.MergeType{}
	if repo.AllowMerge {
		ret = append(ret, git.PullRequestMergeTypeMerge)
	}
	if repo.AllowMerge {
		ret = append(ret, git.PullRequestMergeTypeRebase)
	}
	if repo.AllowSquash {
		ret = append(ret, git.PullRequestMergeTypeSquash)
	}
	return ret
}

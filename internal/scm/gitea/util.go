package gitea

import (
	"code.gitea.io/sdk/gitea"
	"github.com/lindell/multi-gitter/internal/git"
)

// maps merge types to what they are called in the gitea api
var mergeTypeGiteaName = map[git.MergeType]gitea.MergeStyle{
	git.MergeTypeMerge:  gitea.MergeStyleMerge,
	git.MergeTypeRebase: gitea.MergeStyleRebase,
	git.MergeTypeSquash: gitea.MergeStyleSquash,
}

// repoMergeTypes returns a list of all allowed merge types
func repoMergeTypes(repo *gitea.Repository) []git.MergeType {
	ret := []git.MergeType{}
	if repo.AllowMerge {
		ret = append(ret, git.MergeTypeMerge)
	}
	if repo.AllowMerge {
		ret = append(ret, git.MergeTypeRebase)
	}
	if repo.AllowSquash {
		ret = append(ret, git.MergeTypeSquash)
	}
	return ret
}

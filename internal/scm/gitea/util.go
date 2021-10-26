package gitea

import (
	"code.gitea.io/sdk/gitea"
	"github.com/lindell/multi-gitter/internal/scm"
)

// maps merge types to what they are called in the gitea api
var mergeTypeGiteaName = map[scm.MergeType]gitea.MergeStyle{
	scm.MergeTypeMerge:  gitea.MergeStyleMerge,
	scm.MergeTypeRebase: gitea.MergeStyleRebase,
	scm.MergeTypeSquash: gitea.MergeStyleSquash,
}

// repoMergeTypes returns a list of all allowed merge types
func repoMergeTypes(repo *gitea.Repository) []scm.MergeType {
	ret := []scm.MergeType{}
	if repo.AllowMerge {
		ret = append(ret, scm.MergeTypeMerge)
	}
	if repo.AllowMerge {
		ret = append(ret, scm.MergeTypeRebase)
	}
	if repo.AllowSquash {
		ret = append(ret, scm.MergeTypeSquash)
	}
	return ret
}

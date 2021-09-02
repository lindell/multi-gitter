package gitea

import (
	"code.gitea.io/sdk/gitea"
	"github.com/lindell/multi-gitter/internal/pullrequest"
)

// maps merge types to what they are called in the gitea api
var mergeTypeGiteaName = map[pullrequest.MergeType]gitea.MergeStyle{
	pullrequest.MergeTypeMerge:  gitea.MergeStyleMerge,
	pullrequest.MergeTypeRebase: gitea.MergeStyleRebase,
	pullrequest.MergeTypeSquash: gitea.MergeStyleSquash,
}

// repoMergeTypes returns a list of all allowed merge types
func repoMergeTypes(repo *gitea.Repository) []pullrequest.MergeType {
	ret := []pullrequest.MergeType{}
	if repo.AllowMerge {
		ret = append(ret, pullrequest.MergeTypeMerge)
	}
	if repo.AllowMerge {
		ret = append(ret, pullrequest.MergeTypeRebase)
	}
	if repo.AllowSquash {
		ret = append(ret, pullrequest.MergeTypeSquash)
	}
	return ret
}

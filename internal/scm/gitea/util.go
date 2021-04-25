package gitea

import (
	"code.gitea.io/sdk/gitea"
	"github.com/lindell/multi-gitter/internal/domain"
)

// maps merge types to what they are called in the gitea api
var mergeTypeGiteaName = map[domain.MergeType]gitea.MergeStyle{
	domain.MergeTypeMerge:  gitea.MergeStyleMerge,
	domain.MergeTypeRebase: gitea.MergeStyleRebase,
	domain.MergeTypeSquash: gitea.MergeStyleSquash,
}

// repoMergeTypes returns a list of all allowed merge types
func repoMergeTypes(repo *gitea.Repository) []domain.MergeType {
	ret := []domain.MergeType{}
	if repo.AllowMerge {
		ret = append(ret, domain.MergeTypeMerge)
	}
	if repo.AllowMerge {
		ret = append(ret, domain.MergeTypeRebase)
	}
	if repo.AllowSquash {
		ret = append(ret, domain.MergeTypeSquash)
	}
	return ret
}

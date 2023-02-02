package github

import (
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/lindell/multi-gitter/internal/scm"
)

// maps merge types to what they are called in the github api
var mergeTypeGhName = map[scm.MergeType]string{
	scm.MergeTypeMerge:  "merge",
	scm.MergeTypeRebase: "rebase",
	scm.MergeTypeSquash: "squash",
}

// repoMergeTypes returns a list of all allowed merge types
func repoMergeTypes(repo *github.Repository) []scm.MergeType {
	ret := []scm.MergeType{}
	if repo.GetAllowMergeCommit() {
		ret = append(ret, scm.MergeTypeMerge)
	}
	if repo.GetAllowRebaseMerge() {
		ret = append(ret, scm.MergeTypeRebase)
	}
	if repo.GetAllowSquashMerge() {
		ret = append(ret, scm.MergeTypeSquash)
	}
	return ret
}

func stripSuffixIfExist(str string, suffix string) string {
	if strings.HasSuffix(str, suffix) {
		return str[:len(str)-len(suffix)]
	}
	return str
}

func compareStringSlice(slice1 []string, slice2 []string) bool {
	for _, v1 := range slice1 {
		for _, v2 := range slice2 {
			if v1 == v2 {
				return true
			}
		}
	}

	return false
}

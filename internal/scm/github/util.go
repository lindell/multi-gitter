package github

import (
	"strings"

	"github.com/google/go-github/v58/github"
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

func chunkSlice[T any](stack []T, chunkSize int) [][]T {
	var chunks = make([][]T, 0, (len(stack)/chunkSize)+1)
	for chunkSize < len(stack) {
		stack, chunks = stack[chunkSize:], append(chunks, stack[0:chunkSize:chunkSize])
	}

	return append(chunks, stack)
}

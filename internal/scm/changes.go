package scm

import (
	"context"

	"github.com/lindell/multi-gitter/internal/git"
)

// ChangePusher makes a commit through the API
type ChangePusher interface {
	Push(
		ctx context.Context,
		repo Repository,
		commitMessage string,
		change git.Changes,
		featureBranch string,
		branchExist bool,
		forcePush bool,
	) error
}

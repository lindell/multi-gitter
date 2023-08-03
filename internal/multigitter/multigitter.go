package multigitter

import (
	"context"

	"github.com/lindell/multi-gitter/internal/scm"
)

// VersionController fetches repositories
type VersionController interface {
	GetRepositories(ctx context.Context) ([]scm.Repository, error)
	CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error)
	GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error)
	GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error)
	MergePullRequest(ctx context.Context, pr scm.PullRequest) error
	IsPullRequestApprovedByMe(ctx context.Context, pullReq scm.PullRequest) (bool, error)
	ApprovePullRequest(ctx context.Context, pr scm.PullRequest, comment string) error
	RejectPullRequest(ctx context.Context, pr scm.PullRequest, comment string) error
	CommentPullRequest(ctx context.Context, pr scm.PullRequest, comment string) error
	DiffPullRequest(ctx context.Context, pr scm.PullRequest) (string, error)
	ClosePullRequest(ctx context.Context, pr scm.PullRequest) error
	ForkRepository(ctx context.Context, repo scm.Repository, newOwner string) (scm.Repository, error)
}

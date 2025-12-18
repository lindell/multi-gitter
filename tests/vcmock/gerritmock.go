package vcmock

import (
	"github.com/lindell/multi-gitter/internal/scm"
	"golang.org/x/net/context"
)

type GerritVersionController struct {
	VC VersionController
}

// Here we are mocking functions only used for Gerrit VersionController

func (gm *GerritVersionController) EnhanceCommit(_ context.Context, _ scm.Repository, branchName string, commitMessage string) (string, error) {
	message := commitMessage
	message = message + "\n\n" + "Mocked-Footer" + ": " + branchName
	return message, nil
}

func (gm *GerritVersionController) RemoteReference(_ string, featureBranch string, _ bool, _ bool) string {
	return "refs/heads/mocked-" + featureBranch
}

func (gm *GerritVersionController) FeatureBranchExist(ctx context.Context, repo scm.Repository, branchName string) (bool, error) {
	pr, err := gm.VC.GetOpenPullRequest(ctx, repo, branchName)
	return pr != nil, err
}

// Below we are just calling the mock VersionController

func (gm *GerritVersionController) ForkRepository(ctx context.Context, repo scm.Repository, newOwner string) (scm.Repository, error) {
	return gm.VC.ForkRepository(ctx, repo, newOwner)
}

func (gm *GerritVersionController) GetRepositories(ctx context.Context) ([]scm.Repository, error) {
	return gm.VC.GetRepositories(ctx)
}

func (gm *GerritVersionController) CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
	return gm.VC.CreatePullRequest(ctx, repo, prRepo, newPR)
}

func (gm *GerritVersionController) UpdatePullRequest(ctx context.Context, repo scm.Repository, pullReq scm.PullRequest, updatedPR scm.NewPullRequest) (scm.PullRequest, error) {
	return gm.VC.UpdatePullRequest(ctx, repo, pullReq, updatedPR)
}

func (gm *GerritVersionController) GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error) {
	return gm.VC.GetPullRequests(ctx, branchName)
}

func (gm *GerritVersionController) GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
	return gm.VC.GetOpenPullRequest(ctx, repo, branchName)
}

func (gm *GerritVersionController) MergePullRequest(ctx context.Context, pr scm.PullRequest) error {
	return gm.VC.MergePullRequest(ctx, pr)
}

func (gm *GerritVersionController) ClosePullRequest(ctx context.Context, pr scm.PullRequest) error {
	return gm.VC.ClosePullRequest(ctx, pr)
}

func (gm *GerritVersionController) Clean() {
	gm.VC.Clean()
}

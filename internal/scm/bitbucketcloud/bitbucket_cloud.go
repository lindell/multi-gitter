package bitbucketcloud

import (
    "context"
    "github.com/lindell/multi-gitter/internal/scm"
)

type BitbucketCloud struct{}

func (bc *BitbucketCloud) GetRepositories(ctx context.Context) ([]scm.Repository, error) {
    //TODO implement me
    panic("implement me")
}

func (bc *BitbucketCloud) CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
    //TODO implement me
    panic("implement me")
}

func (bc *BitbucketCloud) UpdatePullRequest(ctx context.Context, repo scm.Repository, pullReq scm.PullRequest, updatedPR scm.NewPullRequest) (scm.PullRequest, error) {
    //TODO implement me
    panic("implement me")
}

func (bc *BitbucketCloud) GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error) {
    //TODO implement me
    panic("implement me")
}

func (bc *BitbucketCloud) GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
    //TODO implement me
    panic("implement me")
}

func (bc *BitbucketCloud) MergePullRequest(ctx context.Context, pr scm.PullRequest) error {
    //TODO implement me
    panic("implement me")
}

func (bc *BitbucketCloud) ClosePullRequest(ctx context.Context, pr scm.PullRequest) error {
    //TODO implement me
    panic("implement me")
}

func (bc *BitbucketCloud) ForkRepository(ctx context.Context, repo scm.Repository, newOwner string) (scm.Repository, error) {
    //TODO implement me
    panic("implement me")
}

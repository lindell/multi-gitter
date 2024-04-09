package bitbucketcloud

import (
    "context"
    "github.com/ktrysmt/go-bitbucket"
    "github.com/lindell/multi-gitter/internal/scm"
    "net/http"
)

type BitbucketCloud struct{
    repositories *bitbucket.Repositories
    workspaces *bitbucket.Workspace
    user bitbucket.User
    bearerToken string
    sshAuth         bool
    httpClient      *http.Client
    bbClient *bitbucket.Client
}

func New(username string, bearerToken string, repositories *bitbucket.Repositories, workspaces *bitbucket.Workspace, sshAuth bool)(BitbucketCloud, error){
    //TODO add logic to create client here and populate it with the values present here
    return BitbucketCloud{}, nil
}

func (bbc *BitbucketCloud) CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
    //TODO update to use configured settings and actually handle responses

    opt := &bitbucket.PullRequestsOptions{
        Owner:             "your-team",
        RepoSlug:          "awesome-project",
        SourceBranch:      "develop",
        DestinationBranch: "master",
        Title:             "fix bug. #9999",
        CloseSourceBranch: true,
    }

    _, _ = bbc.bbClient.Repositories.PullRequests.Create(opt)

    panic("implement me")
}

func (bbc *BitbucketCloud) UpdatePullRequest(ctx context.Context, repo scm.Repository, pullReq scm.PullRequest, updatedPR scm.NewPullRequest) (scm.PullRequest, error) {
    //TODO implement me
    panic("implement me")
}

func (bbc *BitbucketCloud) GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error) {
    //TODO implement me
    panic("implement me")
}

func (bbc *BitbucketCloud) GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
    //TODO implement me
    panic("implement me")
}

func (bbc *BitbucketCloud) MergePullRequest(ctx context.Context, pr scm.PullRequest) error {
    //TODO implement me
    panic("implement me")
}

func (bbc *BitbucketCloud) ClosePullRequest(ctx context.Context, pr scm.PullRequest) error {
    //TODO implement me
    panic("implement me")
}

func (bbc *BitbucketCloud) GetRepositories(ctx context.Context) ([]scm.Repository, error) {
    //TODO implement me
    panic("implement me")
}

func (bbc *BitbucketCloud) ForkRepository(ctx context.Context, repo scm.Repository, newOwner string) (scm.Repository, error) {
    //TODO implement me
    panic("implement me")
}

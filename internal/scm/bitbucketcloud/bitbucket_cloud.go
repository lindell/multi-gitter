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
    //type Client struct {
    //    Auth         *auth
    //    Users        users
    //    User         user
    //    Teams        teams
    //    Repositories *Repositories
    //    Workspaces   *Workspace
    //    Pagelen      int
    //    MaxDepth     int
    //    // LimitPages limits the number of pages for a request
    //    //	default value as 0 -- disable limits
    //    LimitPages int
    //    // DisableAutoPaging allows you to disable the default behavior of automatically requesting
    //    // all the pages for a paginated response.
    //    DisableAutoPaging bool
    //
    //    HttpClient *http.Client
    //    // contains filtered or unexported fields
    //}
}

func (bc *BitbucketCloud) CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
    //TODO implement me
    //ctoken := bitbucket
    c := bitbucket.NewBasicAuth("username", "password")

    opt := &bitbucket.PullRequestsOptions{
        Owner:             "your-team",
        RepoSlug:          "awesome-project",
        SourceBranch:      "develop",
        DestinationBranch: "master",
        Title:             "fix bug. #9999",
        CloseSourceBranch: true,
    }

    _, _ = c.Repositories.PullRequests.Create(opt)

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

func (bc *BitbucketCloud) GetRepositories(ctx context.Context) ([]scm.Repository, error) {
    //TODO implement me
    panic("implement me")
}

func (bc *BitbucketCloud) ForkRepository(ctx context.Context, repo scm.Repository, newOwner string) (scm.Repository, error) {
    //TODO implement me
    panic("implement me")
}

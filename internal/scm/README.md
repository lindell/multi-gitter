# Source Control Managers

This folder contains all Source Control Managers. They do all implement the `VersionController` interface described below.

```go
type VersionController interface {
    // Should get repositories based on the configuration
	GetRepositories(ctx context.Context) ([]domain.Repository, error)
    // Creates a pull request. The repo parameter can be expected to originate from the GetRepositories function in the same package
	CreatePullRequest(ctx context.Context, repo domain.Repository, newPR domain.NewPullRequest) (domain.PullRequest, error)

	GetPullRequestStatuses(ctx context.Context, branchName string) ([]domain.PullRequest, error)

	MergePullRequest(ctx context.Context, pr domain.PullRequest) error

	ClosePullRequest(ctx context.Context, pr domain.PullRequest) error
}
```

The version controller can also implement additional functions to support features such as shell-autocompletion.

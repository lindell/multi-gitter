# Source Control Managers

This folder contains all Source Control Managers. They do all implement the `VersionController` interface described below.

```go
type VersionController interface {
	// Should get repositories based on the scm configuration
	GetRepositories(ctx context.Context) ([]domain.Repository, error)
	// Creates a pull request. The repo parameter will always originate from the same package
	CreatePullRequest(ctx context.Context, repo domain.Repository, newPR domain.NewPullRequest) (domain.PullRequest, error)
	// Gets the latest pull requests from repositories based on the scm configuration
	GetPullRequests(ctx context.Context, branchName string) ([]domain.PullRequest, error)
	// Merges a pull request, the pr parameter will always originate from the same package
	MergePullRequest(ctx context.Context, pr domain.PullRequest) error
	// Close a pull request, the pr parameter will always originate from the same package
	ClosePullRequest(ctx context.Context, pr domain.PullRequest) error
}
```

The version controller can also implement additional functions to support features such as shell-autocompletion.

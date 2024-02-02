# Source Control Managers

This folder contains all Source Control Managers. They do all implement the `VersionController` interface described below.

```go
type VersionController interface {
	// Should get repositories based on the scm configuration
	GetRepositories(ctx context.Context) ([]scm.Repository, error)
	// Creates a pull request. The repo parameter will always originate from the same package
	CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error)
	// Gets the latest pull requests from repositories based on the scm configuration
	GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error)
	// Merges a pull request, the pr parameter will always originate from the same package
	MergePullRequest(ctx context.Context, pr scm.PullRequest) error
	// IsPullRequestApprovedByMe returns true if the pr is approved by the current user
	IsPullRequestApprovedByMe(ctx context.Context, pullReq scm.PullRequest) (bool, error)
	// ApprovePullRequest approves a pull request
	ApprovePullRequest(ctx context.Context, pr scm.PullRequest, comment string) error
	// RejectPullRequest requests changes (Note for gitlab this just leaves a comment)
	RejectPullRequest(ctx context.Context, pr scm.PullRequest, comment string) error
	// CommentPullRequest leaves a comment
	CommentPullRequest(ctx context.Context, pr scm.PullRequest, comment string) error
	// DiffPullRequest returns a diff of the pull request
	DiffPullRequest(ctx context.Context, pr scm.PullRequest) (string, error)
	// Close a pull request, the pr parameter will always originate from the same package
	ClosePullRequest(ctx context.Context, pr scm.PullRequest) error
	// ForkRepository forks a repository. If newOwner is set, use it, otherwise fork to the current user
	ForkRepository(ctx context.Context, repo scm.Repository, newOwner string) (scm.Repository, error)
}
```


## Autocompletion

The version controller can also implement additional functions to support features such as shell-autocompletion. The following functions can be implemented independently and will automatically be used for tab completions when the user has activated it.

```go
func GetAutocompleteOrganizations(ctx context.Context, search string) ([]string, error)
```
```go
func GetAutocompleteUsers(ctx context.Context, search string) ([]string, error)
```
```go
func GetAutocompleteRepositories(ctx context.Context, search string) ([]string, error)
```

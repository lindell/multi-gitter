// This package contains a a mock version controller (github/gitlab etc.)

package vcmock

import (
	"context"
	"errors"

	"github.com/lindell/multi-gitter/internal/domain"
)

// VersionController is a mock of an version controller (Github/Gitlab/etc.)
type VersionController struct {
	Repositories []Repository
	PullRequests []PullRequest
}

// GetRepositories returns mock repositories
func (vc *VersionController) GetRepositories(ctx context.Context) ([]domain.Repository, error) {
	ret := make([]domain.Repository, len(vc.Repositories))
	for i := range vc.Repositories {
		ret[i] = vc.Repositories[i]
	}
	return ret, nil
}

// CreatePullRequest stores a mock pull request
func (vc *VersionController) CreatePullRequest(ctx context.Context, repo domain.Repository, newPR domain.NewPullRequest) error {
	repository := repo.(Repository)

	vc.PullRequests = append(vc.PullRequests, PullRequest{
		PRStatus:       domain.PullRequestStatusPending,
		Repository:     repository,
		NewPullRequest: newPR,
	})

	return nil
}

// GetPullRequestStatuses gets mock pull request statuses
func (vc *VersionController) GetPullRequestStatuses(ctx context.Context, branchName string) ([]domain.PullRequest, error) {
	ret := make([]domain.PullRequest, 0, len(vc.PullRequests))
	for _, pr := range vc.PullRequests {
		if pr.NewPullRequest.Head == branchName {
			ret = append(ret, pr)
		}
	}
	return ret, nil
}

// MergePullRequest sets the status of a mock pull requests to merged
func (vc *VersionController) MergePullRequest(ctx context.Context, pr domain.PullRequest) error {
	pullRequest := pr.(PullRequest)
	for i := range vc.PullRequests {
		if vc.PullRequests[i].Repository.Name == pullRequest.Repository.Name {
			vc.PullRequests[i].PRStatus = domain.PullRequestStatusMerged
			return nil
		}
	}
	return errors.New("could not find pull request")
}

// ClosePullRequest sets the status of a mock pull requests to closed
func (vc *VersionController) ClosePullRequest(ctx context.Context, pr domain.PullRequest) error {
	pullRequest := pr.(PullRequest)
	for i := range vc.PullRequests {
		if vc.PullRequests[i].Repository.Name == pullRequest.Repository.Name {
			vc.PullRequests[i].PRStatus = domain.PullRequestStatusClosed
			return nil
		}
	}
	return errors.New("could not find pull request")
}

// AddRepository adds a repository to the mock
func (vc *VersionController) AddRepository(repo ...Repository) {
	vc.Repositories = append(vc.Repositories, repo...)
}

// SetPRStatus sets the status of a pull request
func (vc *VersionController) SetPRStatus(repoName string, branchName string, newStatus domain.PullRequestStatus) {
	for i := range vc.PullRequests {
		if vc.PullRequests[i].Repository.Name == repoName && vc.PullRequests[i].Head == branchName {
			vc.PullRequests[i].PRStatus = newStatus
		}
	}
}

// PullRequest is a mock pr
type PullRequest struct {
	PRStatus domain.PullRequestStatus
	Merged   bool

	Repository
	domain.NewPullRequest
}

// Status returns the pr status
func (pr PullRequest) Status() domain.PullRequestStatus {
	return pr.PRStatus
}

// String return a description of the pr
func (pr PullRequest) String() string {
	return pr.Repository.Name + "/XX"
}

// Repository is a mock repository
type Repository struct {
	Name string
	Path string
}

// URL return the URL (filepath) of the repository on disk
func (r Repository) URL(token string) string {
	return "file://" + r.Path
}

// DefaultBranch returns "master"
func (r Repository) DefaultBranch() string {
	return "master"
}

// FullName returns the name of the mock repo
func (r Repository) FullName() string {
	return r.Name
}

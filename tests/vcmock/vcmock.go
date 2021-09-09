// This package contains a a mock version controller (github/gitlab etc.)

package vcmock

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	git2 "github.com/lindell/multi-gitter/internal/git"
)

// VersionController is a mock of an version controller (Github/Gitlab/etc.)
type VersionController struct {
	PRNumber     int
	Repositories []Repository
	PullRequests []PullRequest
}

// GetRepositories returns mock repositories
func (vc *VersionController) GetRepositories(ctx context.Context) ([]git2.Repository, error) {
	ret := make([]git2.Repository, len(vc.Repositories))
	for i := range vc.Repositories {
		ret[i] = vc.Repositories[i]
	}
	return ret, nil
}

// CreatePullRequest stores a mock pull request
func (vc *VersionController) CreatePullRequest(ctx context.Context, repo git2.Repository, prRepo git2.Repository, newPR git2.NewPullRequest) (git2.PullRequest, error) {
	repository := repo.(Repository)

	vc.PRNumber++
	pr := PullRequest{
		PRStatus:       git2.StatusPending,
		PRNumber:       vc.PRNumber,
		Repository:     repository,
		NewPullRequest: newPR,
	}
	vc.PullRequests = append(vc.PullRequests, pr)

	return pr, nil
}

// GetPullRequests gets mock pull request statuses
func (vc *VersionController) GetPullRequests(ctx context.Context, branchName string) ([]git2.PullRequest, error) {
	ret := make([]git2.PullRequest, 0, len(vc.PullRequests))
	for _, pr := range vc.PullRequests {
		if pr.NewPullRequest.Head == branchName {
			ret = append(ret, pr)
		}
	}
	return ret, nil
}

// MergePullRequest sets the status of a mock pull requests to merged
func (vc *VersionController) MergePullRequest(ctx context.Context, pr git2.PullRequest) error {
	pullRequest := pr.(PullRequest)
	for i := range vc.PullRequests {
		if vc.PullRequests[i].Repository.FullName() == pullRequest.Repository.FullName() {
			vc.PullRequests[i].PRStatus = git2.StatusMerged
			return nil
		}
	}
	return errors.New("could not find pull request")
}

// ClosePullRequest sets the status of a mock pull requests to closed
func (vc *VersionController) ClosePullRequest(ctx context.Context, pr git2.PullRequest) error {
	pullRequest := pr.(PullRequest)
	for i := range vc.PullRequests {
		if vc.PullRequests[i].Repository.FullName() == pullRequest.Repository.FullName() {
			vc.PullRequests[i].PRStatus = git2.StatusClosed
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
func (vc *VersionController) SetPRStatus(repoName string, branchName string, newStatus git2.PullRequestStatus) {
	for i := range vc.PullRequests {
		if vc.PullRequests[i].Repository.RepoName == repoName && vc.PullRequests[i].Head == branchName {
			vc.PullRequests[i].PRStatus = newStatus
		}
	}
}

// GetAutocompleteOrganizations gets organizations for autocompletion
func (vc *VersionController) GetAutocompleteOrganizations(ctx context.Context, str string) ([]string, error) {
	return []string{"static-org", str}, nil
}

// GetAutocompleteUsers gets users for autocompletion
func (vc *VersionController) GetAutocompleteUsers(ctx context.Context, str string) ([]string, error) {
	return []string{"static-user", str}, nil
}

// GetAutocompleteRepositories gets repositories for autocompletion
func (vc *VersionController) GetAutocompleteRepositories(ctx context.Context, str string) ([]string, error) {
	return []string{"static-repo", str}, nil
}

// ForkRepository forks a repository
func (vc *VersionController) ForkRepository(ctx context.Context, repo git2.Repository, newOwner string) (git2.Repository, error) {
	r := repo.(Repository)

	if newOwner == "" {
		newOwner = "default-owner"
	}

	newPath := fmt.Sprintf("%s-forked-%s", r.Path, newOwner)

	_, err := git.PlainCloneContext(ctx, newPath, false, &git.CloneOptions{
		URL: fmt.Sprintf(`file://"%s"`, filepath.ToSlash(r.Path)),
	})
	if err != nil {
		return nil, err
	}

	return Repository{
		OwnerName: newOwner,
		RepoName:  r.RepoName,
		Path:      newPath,
	}, nil
}

// Clean cleans up the data on disk that exist within the version controller mock
func (vc *VersionController) Clean() {
	for _, repo := range vc.Repositories {
		repo.Delete()
	}
}

// PullRequest is a mock pr
type PullRequest struct {
	PRStatus git2.PullRequestStatus
	PRNumber int
	Merged   bool

	Repository
	git2.NewPullRequest
}

// Status returns the pr status
func (pr PullRequest) Status() git2.PullRequestStatus {
	return pr.PRStatus
}

// String return a description of the pr
func (pr PullRequest) String() string {
	return fmt.Sprintf("%s #%d", pr.Repository.FullName(), pr.PRNumber)
}

// Repository is a mock repository
type Repository struct {
	OwnerName string
	RepoName  string
	Path      string
}

// CloneURL return the URL (filepath) of the repository on disk
func (r Repository) CloneURL() string {
	return fmt.Sprintf(`file://%s`, filepath.ToSlash(r.Path))
}

// DefaultBranch returns "master"
func (r Repository) DefaultBranch() string {
	return "master"
}

// FullName returns the name of the mock repo
func (r Repository) FullName() string {
	return fmt.Sprintf("%s/%s", r.OwnerName, r.RepoName)
}

// Owner returns the owner of a repo
func (r Repository) Owner() string {
	return r.OwnerName
}

// Delete deletes data on disk
func (r Repository) Delete() {
	os.RemoveAll(r.Path)
}

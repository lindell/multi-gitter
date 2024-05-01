// This package contains a a mock version controller (github/gitlab etc.)

package vcmock

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	git "github.com/go-git/go-git/v5"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/pkg/errors"
)

// VersionController is a mock of an version controller (Github/Gitlab/etc.)
type VersionController struct {
	PRNumber     int
	Repositories []Repository
	PullRequests []PullRequest

	prLock sync.RWMutex
}

// GetRepositories returns mock repositories
func (vc *VersionController) GetRepositories(_ context.Context) ([]scm.Repository, error) {
	ret := make([]scm.Repository, len(vc.Repositories))
	for i := range vc.Repositories {
		ret[i] = vc.Repositories[i]
	}
	return ret, nil
}

// CreatePullRequest stores a mock pull request
func (vc *VersionController) CreatePullRequest(_ context.Context, repo scm.Repository, _ scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
	repository := repo.(Repository)

	vc.prLock.Lock()
	defer vc.prLock.Unlock()

	vc.PRNumber++
	pr := PullRequest{
		PRStatus:       scm.PullRequestStatusPending,
		PRNumber:       vc.PRNumber,
		Repository:     repository,
		NewPullRequest: newPR,
	}
	vc.PullRequests = append(vc.PullRequests, pr)

	return pr, nil
}

// UpdatePullRequest updates an existing mock pull request
func (vc *VersionController) UpdatePullRequest(_ context.Context, _ scm.Repository, pullReq scm.PullRequest, updatedPR scm.NewPullRequest) (scm.PullRequest, error) {
	pullRequest := pullReq.(PullRequest)

	vc.prLock.Lock()
	defer vc.prLock.Unlock()

	for i := range vc.PullRequests {
		if vc.PullRequests[i].PRNumber == pullRequest.PRNumber && vc.PullRequests[i].Repository.FullName() == pullRequest.Repository.FullName() {
			vc.PullRequests[i].Title = updatedPR.Title
			vc.PullRequests[i].Body = updatedPR.Body
			vc.PullRequests[i].Reviewers = updatedPR.Reviewers
			vc.PullRequests[i].TeamReviewers = updatedPR.TeamReviewers
			vc.PullRequests[i].Assignees = updatedPR.Assignees
			vc.PullRequests[i].Labels = updatedPR.Labels
			return vc.PullRequests[i], nil
		}
	}
	return nil, errors.New("could not find pull request")
}

// GetPullRequests gets mock pull request statuses
func (vc *VersionController) GetPullRequests(_ context.Context, branchName string) ([]scm.PullRequest, error) {
	vc.prLock.RLock()
	defer vc.prLock.RUnlock()

	ret := make([]scm.PullRequest, 0, len(vc.PullRequests))
	for _, pr := range vc.PullRequests {
		if pr.NewPullRequest.Head == branchName {
			ret = append(ret, pr)
		}
	}
	return ret, nil
}

// GetOpenPullRequest gets mock open pull request
func (vc *VersionController) GetOpenPullRequest(_ context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
	vc.prLock.RLock()
	defer vc.prLock.RUnlock()

	r := repo.(Repository)

	for _, pr := range vc.PullRequests {
		if r.OwnerName == pr.OwnerName && r.RepoName == pr.RepoName && pr.NewPullRequest.Head == branchName && openPullRequest(pr) {
			return pr, nil
		}
	}
	return nil, nil
}

func openPullRequest(pr PullRequest) bool {
	return pr.PRStatus == scm.PullRequestStatusSuccess || pr.PRStatus == scm.PullRequestStatusPending
}

// MergePullRequest sets the status of a mock pull requests to merged
func (vc *VersionController) MergePullRequest(_ context.Context, pr scm.PullRequest) error {
	vc.prLock.Lock()
	defer vc.prLock.Unlock()

	pullRequest := pr.(PullRequest)
	for i := range vc.PullRequests {
		if vc.PullRequests[i].Repository.FullName() == pullRequest.Repository.FullName() {
			vc.PullRequests[i].PRStatus = scm.PullRequestStatusMerged
			return nil
		}
	}
	return errors.New("could not find pull request")
}

// ClosePullRequest sets the status of a mock pull requests to closed
func (vc *VersionController) ClosePullRequest(_ context.Context, pr scm.PullRequest) error {
	vc.prLock.Lock()
	defer vc.prLock.Unlock()

	pullRequest := pr.(PullRequest)
	for i := range vc.PullRequests {
		if vc.PullRequests[i].Repository.FullName() == pullRequest.Repository.FullName() {
			vc.PullRequests[i].PRStatus = scm.PullRequestStatusClosed
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
func (vc *VersionController) SetPRStatus(repoName string, branchName string, newStatus scm.PullRequestStatus) {
	vc.prLock.Lock()
	defer vc.prLock.Unlock()

	for i := range vc.PullRequests {
		if vc.PullRequests[i].Repository.RepoName == repoName && vc.PullRequests[i].Head == branchName {
			vc.PullRequests[i].PRStatus = newStatus
		}
	}
}

// GetAutocompleteOrganizations gets organizations for autocompletion
func (vc *VersionController) GetAutocompleteOrganizations(_ context.Context, str string) ([]string, error) {
	return []string{"static-org", str}, nil
}

// GetAutocompleteUsers gets users for autocompletion
func (vc *VersionController) GetAutocompleteUsers(_ context.Context, str string) ([]string, error) {
	return []string{"static-user", str}, nil
}

// GetAutocompleteRepositories gets repositories for autocompletion
func (vc *VersionController) GetAutocompleteRepositories(_ context.Context, str string) ([]string, error) {
	return []string{"static-repo", str}, nil
}

// ForkRepository forks a repository
func (vc *VersionController) ForkRepository(ctx context.Context, repo scm.Repository, newOwner string) (scm.Repository, error) {
	r := repo.(Repository)

	if newOwner == "" {
		newOwner = "default-owner"
	}

	newPath := fmt.Sprintf("%s-forked-%s", r.Path, newOwner)

	_, err := git.PlainCloneContext(ctx, newPath, false, &git.CloneOptions{
		URL: fmt.Sprintf(`file://%s`, filepath.ToSlash(r.Path)),
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
	PRStatus scm.PullRequestStatus
	PRNumber int
	Merged   bool

	Repository
	scm.NewPullRequest
}

// Status returns the pr status
func (pr PullRequest) Status() scm.PullRequestStatus {
	return pr.PRStatus
}

// String return a description of the pr
func (pr PullRequest) String() string {
	return fmt.Sprintf("%s #%d", pr.Repository.FullName(), pr.PRNumber)
}

func (pr PullRequest) URL() string {
	if pr.Repository.RepoName == "has-url" {
		return "https://github.com/owner/has-url/pull/1"
	}

	return ""
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

package repocounter

import (
	"fmt"
	"strings"
	"sync"

	"github.com/lindell/multi-gitter/internal/multigitter/terminal"
	"github.com/lindell/multi-gitter/internal/scm"
)

// Counter keeps track of succeeded and failed repositories
type Counter struct {
	successRepositories []repoInfo
	errors              map[string][]repoInfo
	lock                sync.RWMutex
}

type repoInfo struct {
	repository  scm.Repository
	pullRequest scm.PullRequest
}

// NewCounter create a new repo counter
func NewCounter() *Counter {
	return &Counter{
		errors: map[string][]repoInfo{},
	}
}

// AddError add a failing repository together with the error that caused it
func (r *Counter) AddError(err error, repo scm.Repository, pr scm.PullRequest) {
	defer r.lock.Unlock()
	r.lock.Lock()

	msg := err.Error()
	r.errors[msg] = append(r.errors[msg], repoInfo{
		repository:  repo,
		pullRequest: pr,
	})
}

// AddSuccessRepositories adds a repository that succeeded
func (r *Counter) AddSuccessRepositories(repo scm.Repository) {
	defer r.lock.Unlock()
	r.lock.Lock()

	r.successRepositories = append(r.successRepositories, repoInfo{
		repository: repo,
	})
}

// AddSuccessPullRequest adds a pullrequest that succeeded
func (r *Counter) AddSuccessPullRequest(repo scm.Repository, pr scm.PullRequest) {
	defer r.lock.Unlock()
	r.lock.Lock()

	r.successRepositories = append(r.successRepositories, repoInfo{
		repository:  repo,
		pullRequest: pr,
	})
}

// Info returns a formatted string about all repositories
func (r *Counter) Info() string {
	defer r.lock.RUnlock()
	r.lock.RLock()

	var exitInfo string

	for errMsg := range r.errors {
		exitInfo += fmt.Sprintf("%s:\n", strings.ToUpper(errMsg[0:1])+errMsg[1:])
		for _, errInfo := range r.errors[errMsg] {
			if errInfo.pullRequest == nil {
				exitInfo += fmt.Sprintf("  %s\n", errInfo.repository.FullName())
			} else {
				if urler, hasURL := errInfo.pullRequest.(urler); hasURL && urler.URL() != "" {
					exitInfo += fmt.Sprintf("  %s\n", terminal.Link(errInfo.pullRequest.String(), urler.URL()))
				} else {
					exitInfo += fmt.Sprintf("  %s\n", errInfo.pullRequest.String())
				}
			}
		}
	}

	if len(r.successRepositories) > 0 {
		exitInfo += "Repositories with a successful run:\n"
		for _, repo := range r.successRepositories {
			if repo.pullRequest != nil {
				if urler, hasURL := repo.pullRequest.(urler); hasURL && urler.URL() != "" {
					exitInfo += fmt.Sprintf("  %s\n", terminal.Link(repo.pullRequest.String(), urler.URL()))
				} else {
					exitInfo += fmt.Sprintf("  %s\n", repo.pullRequest.String())
				}
			} else {
				exitInfo += fmt.Sprintf("  %s\n", repo.repository.FullName())
			}
		}
	}

	return exitInfo
}

type urler interface {
	URL() string
}

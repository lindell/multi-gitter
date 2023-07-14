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
	successPullRequests []scm.PullRequest
	successRepositories []scm.Repository
	errors              map[string][]errorInfo
	lock                sync.RWMutex
}

type errorInfo struct {
	repository  scm.Repository
	pullRequest scm.PullRequest
}

// NewCounter create a new repo counter
func NewCounter() *Counter {
	return &Counter{
		errors: map[string][]errorInfo{},
	}
}

// AddError add a failing repository together with the error that caused it
func (r *Counter) AddError(err error, repo scm.Repository, pr scm.PullRequest) {
	defer r.lock.Unlock()
	r.lock.Lock()

	msg := err.Error()
	r.errors[msg] = append(r.errors[msg], errorInfo{
		repository:  repo,
		pullRequest: pr,
	})
}

// AddSuccessRepositories adds a repository that succeeded
func (r *Counter) AddSuccessRepositories(repo scm.Repository) {
	defer r.lock.Unlock()
	r.lock.Lock()

	r.successRepositories = append(r.successRepositories, repo)
}

// AddSuccessPullRequest adds a pullrequest that succeeded
func (r *Counter) AddSuccessPullRequest(pr scm.PullRequest) {
	defer r.lock.Unlock()
	r.lock.Lock()

	r.successPullRequests = append(r.successPullRequests, pr)
}

// Info returns a formatted string about all repositories
func (r *Counter) Info() string {
	defer r.lock.RUnlock()
	r.lock.RLock()

	var exitInfo string

	for errMsg := range r.errors {
		exitInfo += fmt.Sprintf("%s:\n", strings.ToUpper(errMsg[0:1])+errMsg[1:])
		for _, err := range r.errors[errMsg] {
			if err.pullRequest == nil {
				exitInfo += fmt.Sprintf("  %s\n", err.repository.FullName())
			} else {
				if urler, ok := err.pullRequest.(urler); ok {
					exitInfo += fmt.Sprintf("  %s\n", terminal.Link(err.pullRequest.String(), urler.URL()))
				} else {
					exitInfo += fmt.Sprintf("  %s\n", err.pullRequest.String())
				}
			}
		}
	}

	if len(r.successPullRequests) > 0 {
		exitInfo += "Repositories with a successful run:\n"
		for _, pr := range r.successPullRequests {
			if urler, ok := pr.(urler); ok {
				exitInfo += fmt.Sprintf("  %s\n", terminal.Link(pr.String(), urler.URL()))
			} else {
				exitInfo += fmt.Sprintf("  %s\n", pr.String())
			}
		}
	}

	if len(r.successRepositories) > 0 {
		exitInfo += "Repositories with a successful run:\n"
		for _, repo := range r.successRepositories {
			exitInfo += fmt.Sprintf("  %s\n", repo.FullName())
		}
	}

	return exitInfo
}

type urler interface {
	URL() string
}

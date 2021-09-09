package repocounter

import (
	"fmt"
	"strings"
	"sync"

	"github.com/lindell/multi-gitter/internal/git"
	"github.com/lindell/multi-gitter/internal/multigitter/terminal"
)

// Counter keeps track of succeeded and failed repositories
type Counter struct {
	successPullRequests []git.PullRequest
	successRepositories []git.Repository
	errorRepositories   map[string][]git.Repository
	lock                sync.RWMutex
}

// NewCounter create a new repo counter
func NewCounter() *Counter {
	return &Counter{
		errorRepositories: map[string][]git.Repository{},
	}
}

// AddError add a failing repository together with the error that caused it
func (r *Counter) AddError(err error, repo git.Repository) {
	defer r.lock.Unlock()
	r.lock.Lock()

	msg := err.Error()
	r.errorRepositories[msg] = append(r.errorRepositories[msg], repo)
}

// AddSuccessRepositories adds a repository that succeeded
func (r *Counter) AddSuccessRepositories(repo git.Repository) {
	defer r.lock.Unlock()
	r.lock.Lock()

	r.successRepositories = append(r.successRepositories, repo)
}

// AddSuccessPullRequest adds a pullrequest that succeeded
func (r *Counter) AddSuccessPullRequest(repo git.PullRequest) {
	defer r.lock.Unlock()
	r.lock.Lock()

	r.successPullRequests = append(r.successPullRequests, repo)
}

// Info returns a formated string about all repositories
func (r *Counter) Info() string {
	defer r.lock.RUnlock()
	r.lock.RLock()

	var exitInfo string

	for errMsg := range r.errorRepositories {
		exitInfo += fmt.Sprintf("%s:\n", strings.ToUpper(errMsg[0:1])+errMsg[1:])
		for _, repo := range r.errorRepositories[errMsg] {
			exitInfo += fmt.Sprintf("  %s\n", repo.FullName())
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

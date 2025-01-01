package repocounter

import (
	"fmt"
	"maps"
	"slices"
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

func (ri repoInfo) String() string {
	if ri.pullRequest == nil {
		return ri.repository.FullName()
	} else {
		if urler, hasURL := ri.pullRequest.(urler); hasURL && urler.URL() != "" {
			return terminal.Link(ri.pullRequest.String(), urler.URL())
		} else {
			return ri.pullRequest.String()
		}
	}
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

	errors := slices.Collect(maps.Keys(r.errors))
	slices.Sort(errors)

	for _, errMsg := range errors {
		exitInfo += fmt.Sprintf("%s:\n", strings.ToUpper(errMsg[0:1])+errMsg[1:])
		for _, errInfo := range r.errors[errMsg] {
			exitInfo += fmt.Sprintf("  %s\n", errInfo.String())
		}
	}

	if len(r.successRepositories) > 0 {
		exitInfo += "Repositories with a successful run:\n"
		for _, repoInfo := range r.successRepositories {
			exitInfo += fmt.Sprintf("  %s\n", repoInfo.String())
		}
	}

	return exitInfo
}

type urler interface {
	URL() string
}

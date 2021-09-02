package gitea

import (
	"fmt"

	"github.com/lindell/multi-gitter/internal/pullrequest"
)

type PullRequest struct {
	ownerName   string
	repoName    string
	branchName  string
	prOwnerName string
	prRepoName  string
	index       int64 // The id of the PR
	webURL      string
	status      pullrequest.Status
}

func (pr PullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.ownerName, pr.repoName, pr.index)
}

func (pr PullRequest) Status() pullrequest.Status {
	return pr.status
}

func (pr PullRequest) URL() string {
	return pr.webURL
}

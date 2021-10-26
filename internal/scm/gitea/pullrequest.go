package gitea

import (
	"fmt"

	"github.com/lindell/multi-gitter/internal/scm"
)

type pullRequest struct {
	ownerName   string
	repoName    string
	branchName  string
	prOwnerName string
	prRepoName  string
	index       int64 // The id of the PR
	webURL      string
	status      scm.PullRequestStatus
}

func (pr pullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.ownerName, pr.repoName, pr.index)
}

func (pr pullRequest) Status() scm.PullRequestStatus {
	return pr.status
}

func (pr pullRequest) URL() string {
	return pr.webURL
}

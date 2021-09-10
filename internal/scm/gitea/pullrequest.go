package gitea

import (
	"fmt"

	"github.com/lindell/multi-gitter/internal/git"
)

type pullRequest struct {
	ownerName   string
	repoName    string
	branchName  string
	prOwnerName string
	prRepoName  string
	index       int64 // The id of the PR
	webURL      string
	status      git.PullRequestStatus
}

func (pr pullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.ownerName, pr.repoName, pr.index)
}

func (pr pullRequest) Status() git.PullRequestStatus {
	return pr.status
}

func (pr pullRequest) URL() string {
	return pr.webURL
}

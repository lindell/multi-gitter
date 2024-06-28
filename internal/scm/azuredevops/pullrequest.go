package azuredevops

import (
	"fmt"

	"github.com/lindell/multi-gitter/internal/scm"
)

type pullRequest struct {
	ownerName               string
	repoName                string
	repoID                  string
	branchName              string
	id                      int
	status                  scm.PullRequestStatus
	lastMergeSourceCommitID string
}

func (pr pullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.ownerName, pr.repoName, pr.id)
}

func (pr pullRequest) Status() scm.PullRequestStatus {
	return pr.status
}

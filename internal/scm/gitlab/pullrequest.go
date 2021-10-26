package gitlab

import (
	"fmt"

	"github.com/lindell/multi-gitter/internal/scm"
)

type pullRequest struct {
	ownerName  string
	repoName   string
	targetPID  int
	sourcePID  int
	branchName string
	iid        int
	webURL     string
	status     scm.PullRequestStatus
}

func (pr pullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.ownerName, pr.repoName, pr.iid)
}

func (pr pullRequest) Status() scm.PullRequestStatus {
	return pr.status
}

func (pr pullRequest) URL() string {
	return pr.webURL
}

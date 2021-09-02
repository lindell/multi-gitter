package gitlab

import (
	"fmt"

	"github.com/lindell/multi-gitter/internal/pullrequest"
)

type PullRequest struct {
	ownerName  string
	repoName   string
	targetPID  int
	sourcePID  int
	branchName string
	iid        int
	webURL     string
	status     pullrequest.Status
}

func (pr PullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.ownerName, pr.repoName, pr.iid)
}

func (pr PullRequest) Status() pullrequest.Status {
	return pr.status
}

func (pr PullRequest) URL() string {
	return pr.webURL
}

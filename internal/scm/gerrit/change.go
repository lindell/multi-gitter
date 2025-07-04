package gerrit

import (
	"fmt"
	"github.com/lindell/multi-gitter/internal/scm"

	gogerrit "github.com/andygrunwald/go-gerrit"
)

type change struct {
	id       string
	project  string
	branch   string
	number   int
	changeID string
	status   scm.PullRequestStatus
	webURL   string
}

func (r change) String() string {
	return fmt.Sprintf("%d: %s", r.number, r.project)
}

func (r change) Status() scm.PullRequestStatus {
	return r.status
}

func (r change) URL() string {
	return r.webURL
}

func convertChange(changeInfo gogerrit.ChangeInfo, baseURL string) scm.PullRequest {
	status := scm.PullRequestStatusUnknown

	if changeInfo.Submittable {
		status = scm.PullRequestStatusSuccess
	} else {
		switch changeInfo.Status {
		case "NEW":
			status = scm.PullRequestStatusPending
		case "MERGED":
			status = scm.PullRequestStatusMerged
		case "ABANDONED":
			status = scm.PullRequestStatusClosed
		}
	}

	return change{
		id:       changeInfo.ID,
		project:  changeInfo.Project,
		branch:   changeInfo.Branch,
		number:   changeInfo.Number,
		changeID: changeInfo.ChangeID,
		status:   status,
		webURL:   fmt.Sprintf("%s/c/%s/+/%d", baseURL, changeInfo.Project, changeInfo.Number),
	}
}

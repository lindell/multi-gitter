package gerrit

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lindell/multi-gitter/internal/scm"

	gogerrit "github.com/andygrunwald/go-gerrit"
)

type change struct {
	id       string
	project  string
	branch   string
	number   int
	changeId string
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

func convertChange(changeInfo gogerrit.ChangeInfo, baseUrl string) scm.PullRequest {
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
		changeId: changeInfo.ChangeID,
		status:   status,
		webURL:   strings.Join([]string{baseUrl, "c", changeInfo.Project, "+", strconv.Itoa(changeInfo.Number)}, "/"),
	}
}

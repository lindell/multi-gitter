package github

import (
	"fmt"

	"github.com/google/go-github/v38/github"
	"github.com/lindell/multi-gitter/internal/pullrequest"
)

func convertPullRequest(pr *github.PullRequest) PullRequest {
	return PullRequest{
		ownerName:   pr.GetBase().GetUser().GetLogin(),
		repoName:    pr.GetBase().GetRepo().GetName(),
		branchName:  pr.GetHead().GetRef(),
		prOwnerName: pr.GetHead().GetUser().GetLogin(),
		prRepoName:  pr.GetHead().GetRepo().GetName(),
		number:      pr.GetNumber(),
		guiURL:      pr.GetHTMLURL(),
	}
}

type PullRequest struct {
	ownerName   string
	repoName    string
	branchName  string
	prOwnerName string
	prRepoName  string
	number      int
	guiURL      string
	status      pullrequest.Status
}

func (pr PullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.ownerName, pr.repoName, pr.number)
}

func (pr PullRequest) Status() pullrequest.Status {
	return pr.status
}

func (pr PullRequest) URL() string {
	return pr.guiURL
}

package github

import (
	"fmt"

	"github.com/google/go-github/v39/github"
	"github.com/lindell/multi-gitter/internal/git"
)

func convertPullRequest(pr *github.PullRequest) pullRequest {
	return pullRequest{
		ownerName:   pr.GetBase().GetUser().GetLogin(),
		repoName:    pr.GetBase().GetRepo().GetName(),
		branchName:  pr.GetHead().GetRef(),
		prOwnerName: pr.GetHead().GetUser().GetLogin(),
		prRepoName:  pr.GetHead().GetRepo().GetName(),
		number:      pr.GetNumber(),
		guiURL:      pr.GetHTMLURL(),
	}
}

type pullRequest struct {
	ownerName   string
	repoName    string
	branchName  string
	prOwnerName string
	prRepoName  string
	number      int
	guiURL      string
	status      git.PullRequestStatus
}

func (pr pullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.ownerName, pr.repoName, pr.number)
}

func (pr pullRequest) Status() git.PullRequestStatus {
	return pr.status
}

func (pr pullRequest) URL() string {
	return pr.guiURL
}

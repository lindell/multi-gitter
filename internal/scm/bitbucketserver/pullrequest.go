package bitbucketserver

import (
	"fmt"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/lindell/multi-gitter/internal/pullrequest"
)

func newPullRequest(pr bitbucketv1.PullRequest) PullRequest {
	return PullRequest{
		project:    pr.ToRef.Repository.Project.Key,
		repoName:   pr.ToRef.Repository.Slug,
		branchName: pr.FromRef.DisplayID,
		prProject:  pr.FromRef.Repository.Project.Key,
		prRepoName: pr.FromRef.Repository.Slug,
		number:     pr.ID,
		guiURL:     pr.Links.Self[0].Href,
	}
}

type PullRequest struct {
	project    string
	repoName   string
	branchName string
	prProject  string
	prRepoName string
	number     int
	guiURL     string
	status     pullrequest.Status
}

func (pr PullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.project, pr.repoName, pr.number)
}

func (pr PullRequest) Status() pullrequest.Status {
	return pr.status
}

func (pr PullRequest) URL() string {
	return pr.guiURL
}

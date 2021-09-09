package bitbucketserver

import (
	"fmt"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/lindell/multi-gitter/internal/git"
)

func newPullRequest(pr bitbucketv1.PullRequest) pullRequest {
	return pullRequest{
		project:    pr.ToRef.Repository.Project.Key,
		repoName:   pr.ToRef.Repository.Slug,
		branchName: pr.FromRef.DisplayID,
		prProject:  pr.FromRef.Repository.Project.Key,
		prRepoName: pr.FromRef.Repository.Slug,
		number:     pr.ID,
		guiURL:     pr.Links.Self[0].Href,
	}
}

type pullRequest struct {
	project    string
	repoName   string
	branchName string
	prProject  string
	prRepoName string
	number     int
	guiURL     string
	status     git.PullRequestStatus
}

func (pr pullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.project, pr.repoName, pr.number)
}

func (pr pullRequest) Status() git.PullRequestStatus {
	return pr.status
}

func (pr pullRequest) URL() string {
	return pr.guiURL
}

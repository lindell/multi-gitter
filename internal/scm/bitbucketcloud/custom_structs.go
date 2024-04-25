package bitbucketcloud

import (
	"fmt"
	
	"github.com/ktrysmt/go-bitbucket"
	"github.com/lindell/multi-gitter/internal/scm"
)

const (
	cloneHTTPType = "https"
	cloneSSHType  = "ssh"
	stateMerged   = "MERGED"
	stateDeclined = "DECLINED"
)

type bitbucketPullRequests struct {
	Next		string			`json:"next"`
	Page		int				`json:"page"`
	PageLen 	int				`json:"pagelen"`
	Previous 	string			`json:"previous"`
	Size		int				`json:"size"`
	Values 		[]bbPullRequest `json:"values"`
}

type bbPullRequest struct {
	State		string			`json:"state"`
	Source 		pullRequestRef	`json:"source"`
	Destination pullRequestRef	`json:"destination"`
	Links 		links 			`json:"links"`
	Title 		string			`json:"title"`
	Type 		string 			`json:"type"`
	ID 			int				`json:"id"`
}

type pullRequestRef struct {
	Branch      branch     				`json:"branch"`
	Commit 		commit 					`json:"Commit"`
	Repository  bitbucket.Repository 	`json:"repository"`
}

type branch struct {
	Name	string 	`json:"name"`
}

type commit struct {
	Hash 	string	`json:"hash"`
	Type 	string 	`json:"type"`
	Links 	links	`json:"links"`
}

type links struct {
	Self hrefLink `json:"self,omitempty"`
	Html hrefLink `json:"html,omitempty"`
}

type repoLinks struct {
	Clone []hrefLink `json:"clone,omitempty"`
	Self []hrefLink `json:"self,omitempty"`
	Html []hrefLink `json:"html,omitempty"`
}

type hrefLink struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

type pullRequest struct {
	project    string
	repoName   string
	branchName string
	prProject  string
	prRepoName string
	number     int
	guiURL     string
	status     scm.PullRequestStatus
}

func (pr pullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.project, pr.repoName, pr.number)
}

func (pr pullRequest) Status() scm.PullRequestStatus {
	return pr.status
}

func (pr pullRequest) URL() string {
	return pr.guiURL
}

// repository contains information about a bitbucket repository
type repository struct {
	name          string
	project       string
	defaultBranch string
	cloneURL      string
}

func (r repository) CloneURL() string {
	return r.cloneURL
}

func (r repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r repository) FullName() string {
	return r.project + "/" + r.name
}
package bitbucketcloud

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/pkg/errors"
)

type BitbucketCloud struct {
	repositories []string
	workspaces   []string
	users        []string
	username     string
	token        string
	sshAuth      bool
	httpClient   *http.Client
	bbClient     *bitbucket.Client
}

func New(username string, token string, repositories []string, workspaces []string, users []string, sshAuth bool, transportMiddleware func(http.RoundTripper) http.RoundTripper) (*BitbucketCloud, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("bearer token is empty")
	}

	bitbucketCloud := &BitbucketCloud{}
	bitbucketCloud.repositories = repositories
	bitbucketCloud.workspaces = workspaces
	bitbucketCloud.users = users
	bitbucketCloud.username = username
	bitbucketCloud.token = token
	bitbucketCloud.sshAuth = sshAuth
	bitbucketCloud.httpClient = &http.Client{
		Transport: transportMiddleware(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false}, // nolint: gosec
		}),
	}
	bitbucketCloud.bbClient = bitbucket.NewBasicAuth(username, token)

	return bitbucketCloud, nil
}

func (bbc *BitbucketCloud) CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {

	splitRepoFullName := strings.Split(prRepo.FullName(), "/")

	repoOptions := &bitbucket.RepositoryOptions{
		Owner:    bbc.workspaces[0],
		RepoSlug: splitRepoFullName[1],
	}
	defaultReviewers, _ := bbc.bbClient.Repositories.Repository.ListDefaultReviewers(repoOptions)
	for _, reviewer := range defaultReviewers.DefaultReviewers {
		newPR.Reviewers = append(newPR.Reviewers, reviewer.Uuid)
	}

	prOptions := &bitbucket.PullRequestsOptions{
		Owner:             bbc.workspaces[0],
		RepoSlug:          splitRepoFullName[1],
		SourceBranch:      newPR.Head,
		DestinationBranch: newPR.Base,
		Title:             newPR.Title,
		CloseSourceBranch: true,
		Reviewers:         newPR.Reviewers,
	}

	_, err := bbc.bbClient.Repositories.PullRequests.Create(prOptions)
	if err != nil {
		return nil, err
	}

	return &pullRequest{
		project:    splitRepoFullName[0],
		repoName:   splitRepoFullName[1],
		branchName: newPR.Head,
		prProject:  splitRepoFullName[0],
		prRepoName: splitRepoFullName[1],
		status:     scm.PullRequestStatusSuccess,
	}, nil
}

func (bbc *BitbucketCloud) UpdatePullRequest(ctx context.Context, repo scm.Repository, pullReq scm.PullRequest, updatedPR scm.NewPullRequest) (scm.PullRequest, error) {
	bbcPR := pullReq.(pullRequest)
	repoSlug := strings.Split(bbcPR.guiURL, "/")

	// Note the specs of the bitbucket client here, reviewers field must be UUID of the reviewers, not their usernames
	prOptions := &bitbucket.PullRequestsOptions{
		ID:                fmt.Sprintf("%d", bbcPR.number),
		Owner:             bbc.workspaces[0],
		RepoSlug:          repoSlug[4],
		Title:             updatedPR.Title,
		Description:       updatedPR.Body,
		CloseSourceBranch: true,
		SourceBranch:      updatedPR.Head,
		DestinationBranch: updatedPR.Base,
		Reviewers:         updatedPR.Reviewers,
	}
	_, err := bbc.bbClient.Repositories.PullRequests.Update(prOptions)
	if err != nil {
		return nil, err
	}

	return &pullRequest{
		project:    bbcPR.project,
		repoName:   bbcPR.repoName,
		branchName: updatedPR.Head,
		prProject:  bbcPR.prProject,
		prRepoName: bbcPR.prRepoName,
		status:     scm.PullRequestStatusSuccess,
	}, nil
}

func (bbc *BitbucketCloud) GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error) {
	var responsePRs []scm.PullRequest
	for _, repoName := range bbc.repositories {
		prs, _ := bbc.bbClient.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{Owner: bbc.workspaces[0], RepoSlug: repoName})
		prBytes, _ := json.Marshal(prs)
		bbPullRequests := &bitbucketPullRequests{}
		err := json.Unmarshal(prBytes, bbPullRequests)
		if err != nil {
			return nil, err
		}
		for _, pr := range bbPullRequests.Values {
			convertedPr, err := bbc.convertPullRequest(bbc.workspaces[0], repoName, &pr)
			if err != nil {
				return nil, err
			}
			responsePRs = append(responsePRs, convertedPr)
		}
	}
	fmt.Println(responsePRs)
	return responsePRs, nil
}

func (bbc *BitbucketCloud) getPullRequests(ctx context.Context, repoName string) ([]pullRequest, error) {

	var repoPRs []pullRequest
	prs, _ := bbc.bbClient.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{Owner: bbc.workspaces[0], RepoSlug: repoName})
	prBytes, _ := json.Marshal(prs)
	bbPullRequests := &bitbucketPullRequests{}
	err := json.Unmarshal(prBytes, bbPullRequests)
	if err != nil {
		return nil, err
	}
	for _, pr := range bbPullRequests.Values {
		convertedPr, err := bbc.convertPullRequest(bbc.workspaces[0], repoName, &pr)
		if err != nil {
			return nil, err
		}
		repoPRs = append(repoPRs, convertedPr)
	}
	return repoPRs, nil
}

func (bbc *BitbucketCloud) convertPullRequest(project, repoName string, pr *bbPullRequest) (pullRequest, error) {
	status, err := bbc.pullRequestStatus(pr)
	if err != nil {
		return pullRequest{}, err
	}

	return pullRequest{
		repoName:   repoName,
		project:    project,
		branchName: pr.Source.Branch.Name,

		prProject:  pr.Source.Repository.Project.Key,
		prRepoName: pr.Source.Repository.Slug,
		number:     pr.ID,

		guiURL: pr.Links.Html.Href,
		status: status,
	}, nil
}

func (bbc *BitbucketCloud) pullRequestStatus(pr *bbPullRequest) (scm.PullRequestStatus, error) {
	switch pr.State {
	case stateMerged:
		return scm.PullRequestStatusMerged, nil
	case stateDeclined:
		return scm.PullRequestStatusClosed, nil
	}

	return scm.PullRequestStatusSuccess, nil
}

func (bbc *BitbucketCloud) GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
	repoFN := repo.FullName()
	repoSlug := strings.Split(repoFN, "/")
	repoPRs, _ := bbc.getPullRequests(ctx, repoSlug[1])
	for _, repoPR := range repoPRs {
		pr := pullRequest(repoPR)
		if pr.branchName == branchName && pr.status == scm.PullRequestStatusSuccess {
			return repoPR, nil
		}
	}
	return nil, nil
}

func (bbc *BitbucketCloud) MergePullRequest(ctx context.Context, pr scm.PullRequest) error {
	bbcPR := pr.(pullRequest)
	repoSlug := strings.Split(bbcPR.guiURL, "/")
	prOptions := &bitbucket.PullRequestsOptions{
		ID:           fmt.Sprintf("%d", bbcPR.number),
		SourceBranch: bbcPR.branchName,
		RepoSlug:     repoSlug[4],
		Owner:        bbc.workspaces[0],
	}
	_, err := bbc.bbClient.Repositories.PullRequests.Merge(prOptions)
	return err
}

func (bbc *BitbucketCloud) ClosePullRequest(ctx context.Context, pr scm.PullRequest) error {
	bbcPR := pr.(pullRequest)
	repoSlug := strings.Split(bbcPR.guiURL, "/")
	prOptions := &bitbucket.PullRequestsOptions{
		ID:           fmt.Sprintf("%d", bbcPR.number),
		SourceBranch: bbcPR.branchName,
		RepoSlug:     repoSlug[4],
		Owner:        bbc.workspaces[0],
	}
	_, err := bbc.bbClient.Repositories.PullRequests.Decline(prOptions)
	return err
}

func (bbc *BitbucketCloud) GetRepositories(ctx context.Context) ([]scm.Repository, error) {
	repoOptions := &bitbucket.RepositoriesOptions{
		Role:  "member",
		Owner: bbc.workspaces[0],
	}

	repos, err := bbc.bbClient.Repositories.ListForAccount(repoOptions)

	if err != nil {
		return nil, err
	}

	repositories := make([]scm.Repository, 0, len(repos.Items))
	for _, repo := range repos.Items {
		if slices.Contains(bbc.repositories, repo.Name) {
			converted, err := bbc.convertRepository(repo)
			if err != nil {
				return nil, err
			}
			repositories = append(repositories, *converted)
		}
	}

	return repositories, nil
}

func (bbc *BitbucketCloud) ForkRepository(ctx context.Context, repo scm.Repository, newOwner string) (scm.Repository, error) {
	splitRepoFullName := strings.Split(repo.FullName(), "/")

	if newOwner == "" {
		newOwner = bbc.username
	}
	options := &bitbucket.RepositoryForkOptions{
		FromOwner: bbc.workspaces[0],
		FromSlug:  splitRepoFullName[1],
		Owner:     newOwner,
		Name:      splitRepoFullName[1],
	}
	// TODO: Support for selecting Bitbucket project to fork into
	resp, err := bbc.bbClient.Repositories.Repository.Fork(options)
	if err != nil {
		return nil, err
	}
	res, err := bbc.convertRepository(*resp)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

func (bbc *BitbucketCloud) convertRepository(repo bitbucket.Repository) (*repository, error) {
	var cloneURL string

	rLinks := &repoLinks{}
	linkBytes, _ := json.Marshal(repo.Links)
	_ = json.Unmarshal(linkBytes, rLinks)

	if bbc.sshAuth {
		cloneURL = findLinkType(rLinks.Clone, cloneSSHType)
		if cloneURL == "" {
			return nil, errors.Errorf("unable to find clone url for repository %s using clone type %s", repo.Name, cloneSSHType)
		}
	} else {
		httpURL := findLinkType(rLinks.Clone, cloneHTTPType)
		if httpURL == "" {
			return nil, errors.Errorf("unable to find clone url for repository %s using clone type %s", repo.Name, cloneHTTPType)
		}
		parsedURL, err := url.Parse(httpURL)
		if err != nil {
			return nil, err
		}

		parsedURL.User = url.UserPassword(bbc.username, bbc.token)
		cloneURL = parsedURL.String()
	}

	return &repository{
		name:          repo.Slug,
		project:       repo.Project.Name,
		defaultBranch: repo.Mainbranch.Name,
		cloneURL:      cloneURL,
	}, nil
}

func findLinkType(cloneLinks []hrefLink, cloneType string) string {
	for _, clone := range cloneLinks {
		if strings.EqualFold(clone.Name, cloneType) {
			return clone.Href
		}
	}

	return ""
}

package bitbucketcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	internalHTTP "github.com/lindell/multi-gitter/internal/http"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/pkg/errors"
)

type BitbucketCloud struct {
	repositories []string
	workspaces   []string
	users        []string
	fork         bool
	username     string
	token        string
	sshAuth      bool
	newOwner     string
	httpClient   *http.Client
	bbClient     *bitbucket.Client
}

func New(username string, token string, repositories []string, workspaces []string, users []string, fork bool, sshAuth bool,
	newOwner string) (*BitbucketCloud, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("bearer token is empty")
	}

	bitbucketCloud := &BitbucketCloud{}
	bitbucketCloud.repositories = repositories
	bitbucketCloud.workspaces = workspaces
	bitbucketCloud.users = users
	bitbucketCloud.fork = fork
	bitbucketCloud.username = username
	bitbucketCloud.token = token
	bitbucketCloud.sshAuth = sshAuth
	bitbucketCloud.newOwner = newOwner
	bitbucketCloud.httpClient = &http.Client{
		Transport: internalHTTP.LoggingRoundTripper{},
	}
	bitbucketCloud.bbClient = bitbucket.NewBasicAuth(username, token)

	return bitbucketCloud, nil
}

func (bbc *BitbucketCloud) CreatePullRequest(_ context.Context, _ scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
	bbcRepo := prRepo.(repository)

	repoOptions := &bitbucket.RepositoryOptions{
		Owner:    bbc.workspaces[0],
		RepoSlug: bbcRepo.name,
	}
	currentUser, err := bbc.bbClient.User.Profile()
	if err != nil {
		return nil, err
	}
	defaultReviewers, err := bbc.bbClient.Repositories.Repository.ListDefaultReviewers(repoOptions)
	if err != nil {
		return nil, err
	}
	for _, reviewer := range defaultReviewers.DefaultReviewers {
		if currentUser.Uuid != reviewer.Uuid {
			newPR.Reviewers = append(newPR.Reviewers, reviewer.Uuid)
		}
	}

	if bbc.newOwner == "" {
		bbc.newOwner = bbc.username
	}

	prOptions := &bitbucket.PullRequestsOptions{
		Owner:             bbc.workspaces[0],
		RepoSlug:          bbcRepo.name,
		SourceBranch:      newPR.Head,
		DestinationBranch: newPR.Base,
		Title:             newPR.Title,
		CloseSourceBranch: true,
		Reviewers:         newPR.Reviewers,
	}

	// If we are performing a fork, set the source repository.
	// We do not if we are just creating a PR within the same repo, as it will cause issues
	if bbc.fork {
		prOptions.SourceRepository = fmt.Sprintf("%s/%s", bbc.newOwner, repoOptions.RepoSlug)
	}

	resp, err := bbc.bbClient.Repositories.PullRequests.Create(prOptions)
	if err != nil {
		return nil, err
	}
	createBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	r := newPrResponse{}
	err = json.Unmarshal(createBytes, &r)
	if err != nil {
		return nil, err
	}
	// Were currently using scm.PullRequestStatusSuccess here for simplicity
	// We could eventually look it up using bbc.pullRequestStatus but we will need to refactor it to support passing in the needed variable
	return &pullRequest{
		number:     r.ID,
		guiURL:     r.Links.HTML.Href,
		project:    bbcRepo.project,
		repoName:   bbcRepo.name,
		branchName: newPR.Head,
		prProject:  bbcRepo.project,
		prRepoName: bbcRepo.name,
		status:     scm.PullRequestStatusSuccess,
	}, nil
}

func (bbc *BitbucketCloud) UpdatePullRequest(_ context.Context, _ scm.Repository, pullReq scm.PullRequest, updatedPR scm.NewPullRequest) (scm.PullRequest, error) {
	bbcPR := pullReq.(pullRequest)

	// Note the specs of the bitbucket client here, reviewers field must be UUID of the reviewers, not their usernames
	prOptions := &bitbucket.PullRequestsOptions{
		ID:                fmt.Sprintf("%d", bbcPR.number),
		Owner:             bbc.workspaces[0],
		RepoSlug:          bbcPR.repoName,
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
		number:     bbcPR.number,
		guiURL:     bbcPR.guiURL,
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
	repositories, err := bbc.GetRepositories(ctx)
	if err != nil {
		return nil, err
	}
	for _, repo := range repositories {
		bbcRepo := repo.(repository)
		prs, err := bbc.bbClient.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{Owner: bbc.workspaces[0], RepoSlug: bbcRepo.name, SourceBranch: branchName})
		if err != nil {
			return nil, err
		}
		prBytes, err := json.Marshal(prs)
		if err != nil {
			return nil, err
		}
		bbPullRequests := &bitbucketPullRequests{}
		err = json.Unmarshal(prBytes, bbPullRequests)
		if err != nil {
			return nil, err
		}
		for _, pr := range bbPullRequests.Values {
			convertedPr := bbc.convertPullRequest(bbc.workspaces[0], bbcRepo.name, &pr)
			responsePRs = append(responsePRs, convertedPr)
		}
	}
	return responsePRs, nil
}

func (bbc *BitbucketCloud) getPullRequests(_ context.Context, repoName string) ([]pullRequest, error) {
	var repoPRs []pullRequest
	prs, err := bbc.bbClient.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{Owner: bbc.workspaces[0], RepoSlug: repoName})
	if err != nil {
		return nil, err
	}
	prBytes, err := json.Marshal(prs)
	if err != nil {
		return nil, err
	}
	bbPullRequests := &bitbucketPullRequests{}
	err = json.Unmarshal(prBytes, bbPullRequests)
	if err != nil {
		return nil, err
	}
	for _, pr := range bbPullRequests.Values {
		convertedPr := bbc.convertPullRequest(bbc.workspaces[0], repoName, &pr)
		repoPRs = append(repoPRs, convertedPr)
	}
	return repoPRs, nil
}

func (bbc *BitbucketCloud) convertPullRequest(project, repoName string, pr *bbPullRequest) pullRequest {
	status := bbc.pullRequestStatus(pr)

	return pullRequest{
		repoName:   repoName,
		project:    project,
		branchName: pr.Source.Branch.Name,

		prProject:  pr.Source.Repository.Project.Key,
		prRepoName: pr.Source.Repository.Slug,
		number:     pr.ID,

		guiURL: pr.Links.HTML.Href,
		status: status,
	}
}

func (bbc *BitbucketCloud) pullRequestStatus(pr *bbPullRequest) scm.PullRequestStatus {
	switch pr.State {
	case stateMerged:
		return scm.PullRequestStatusMerged
	case stateDeclined:
		return scm.PullRequestStatusClosed
	}

	return scm.PullRequestStatusSuccess
}

func extractRepoSlug(bbcPR pullRequest) string {
	repoSlug := strings.Split(bbcPR.guiURL, "/")[4]
	return repoSlug
}

func (bbc *BitbucketCloud) GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
	bbcRepo := repo.(repository)
	repoPRs, err := bbc.getPullRequests(ctx, bbcRepo.name)
	if err != nil {
		return nil, err
	}
	for _, repoPR := range repoPRs {
		pr := repoPR
		if pr.branchName == branchName && pr.status == scm.PullRequestStatusSuccess {
			return repoPR, nil
		}
	}
	return nil, nil
}

func (bbc *BitbucketCloud) MergePullRequest(_ context.Context, pr scm.PullRequest) error {
	bbcPR := pr.(pullRequest)
	repoSlug := extractRepoSlug(bbcPR)
	prOptions := &bitbucket.PullRequestsOptions{
		ID:           fmt.Sprintf("%d", bbcPR.number),
		SourceBranch: bbcPR.branchName,
		RepoSlug:     repoSlug,
		Owner:        bbc.workspaces[0],
	}
	_, err := bbc.bbClient.Repositories.PullRequests.Merge(prOptions)
	return err
}

func (bbc *BitbucketCloud) ClosePullRequest(_ context.Context, pr scm.PullRequest) error {
	bbcPR := pr.(pullRequest)
	repoSlug := extractRepoSlug(bbcPR)
	prOptions := &bitbucket.PullRequestsOptions{
		ID:           fmt.Sprintf("%d", bbcPR.number),
		SourceBranch: bbcPR.branchName,
		RepoSlug:     repoSlug,
		Owner:        bbc.workspaces[0],
	}
	_, err := bbc.bbClient.Repositories.PullRequests.Decline(prOptions)
	return err
}

func (bbc *BitbucketCloud) GetRepositories(_ context.Context) ([]scm.Repository, error) {
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

func (bbc *BitbucketCloud) ForkRepository(_ context.Context, repo scm.Repository, newOwner string) (scm.Repository, error) {
	bbcRepo := repo.(repository)
	if newOwner == "" {
		newOwner = bbc.username
	}
	options := &bitbucket.RepositoryForkOptions{
		FromOwner: bbc.workspaces[0],
		FromSlug:  bbcRepo.name,
		Owner:     newOwner,
		Name:      bbcRepo.name,
	}

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
	linkBytes, err := json.Marshal(repo.Links)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(linkBytes, rLinks)

	if bbc.sshAuth {
		cloneURL, err = findLinkType(rLinks.Clone, cloneSSHType, repo.Name)
		if err != nil {
			return nil, err
		}
	} else {
		httpURL, err := findLinkType(rLinks.Clone, cloneHTTPType, repo.Name)
		if err != nil {
			return nil, err
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

func findLinkType(cloneLinks []hrefLink, cloneType string, repoName string) (string, error) {
	for _, clone := range cloneLinks {
		if strings.EqualFold(clone.Name, cloneType) {
			return clone.Href, nil
		}
	}

	return "", errors.Errorf("unable to find clone url for repository %s using clone type %s", repoName, cloneType)
}

package bitbucketcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	authType     AuthType
	httpClient   *http.Client
	bbClient     *bitbucket.Client
}

func New(username string, token string, repositories []string, workspaces []string, users []string, fork bool, sshAuth bool,
	newOwner string, authType AuthType) (*BitbucketCloud, error) {
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
	bitbucketCloud.authType = authType
	bitbucketCloud.httpClient = &http.Client{
		Transport: internalHTTP.LoggingRoundTripper{},
	}

	if authType == AuthTypeAppPassword {
		// Authenticate using app password
		bitbucketCloud.bbClient = bitbucket.NewBasicAuth(username, token)
	} else if authType == AuthTypeWorkspaceToken {
		// Authenticate using workspace token
		bitbucketCloud.bbClient = bitbucket.NewOAuthbearerToken(token)
	}

	bitbucketCloud.bbClient.HttpClient.Transport = internalHTTP.LoggingRoundTripper{
		Next: bitbucketCloud.bbClient.HttpClient.Transport,
	}

	return bitbucketCloud, nil
}

func (bbc *BitbucketCloud) CreatePullRequest(_ context.Context, _ scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
	bbcRepo := prRepo.(repository)

	repoOptions := &bitbucket.RepositoryOptions{
		Owner:    bbc.workspaces[0],
		RepoSlug: bbcRepo.name,
	}
	var currentUserUUID string
	if bbc.authType == AuthTypeAppPassword {
		currentUser, err := bbc.bbClient.User.Profile()
		if err != nil {
			return nil, err
		}
		currentUserUUID = currentUser.Uuid
	}

	defaultReviewers, err := bbc.bbClient.Repositories.Repository.ListEffectiveDefaultReviewers(repoOptions)
	if err != nil {
		return nil, err
	}
	for _, reviewer := range defaultReviewers.EffectiveDefaultReviewers {
		if currentUserUUID != reviewer.User.Uuid {
			newPR.Reviewers = append(newPR.Reviewers, reviewer.User.Uuid)
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

	// Query for open pull requests on the specific branch
	// The bitbucket API uses `source.branch.name="<branchName>"` in the q parameter for filtering.
	// The go-bitbucket library uses the States field for filtering by PR state.
	queryString := fmt.Sprintf("source.branch.name = \"%s\"", branchName)
	prs, err := bbc.bbClient.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
		Owner:    bbc.workspaces[0],
		RepoSlug: bbcRepo.name,
		Query:    queryString,      // Using the query parameter for branch filtering.
		States:   []string{"OPEN"}, // Use States field for PR state filtering.
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get pull requests for branch %s in repo %s", branchName, bbcRepo.name)
	}

	// The Gets method returns a struct that contains a list of pull requests.
	// We need to iterate through these, though with the query, it should ideally be 0 or 1.
	prBytes, err := json.Marshal(prs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal pull requests response")
	}
	bbPullRequests := &bitbucketPullRequests{}
	err = json.Unmarshal(prBytes, bbPullRequests)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal pull requests response")
	}

	for _, pr := range bbPullRequests.Values {
		// The API query should be precise, but we can double-check here if necessary.
		// if pr.Source.Branch.Name == branchName && pr.State == "OPEN" {
		convertedPR := bbc.convertPullRequest(bbc.workspaces[0], bbcRepo.name, &pr)
		return convertedPR, nil
		// }
	}

	return nil, nil // No open pull request found for the branch
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
	if len(bbc.repositories) > 0 {
		// If specific repository names are provided, fetch them directly
		repositories := make([]scm.Repository, 0, len(bbc.repositories))
		for _, repoName := range bbc.repositories {
			repoDetails, err := bbc.bbClient.Repositories.Repository.Get(&bitbucket.RepositoryOptions{
				Owner:    bbc.workspaces[0],
				RepoSlug: repoName,
			})
			if err != nil {
				// It might be desirable to collect errors or decide if one failure should stop all
				return nil, errors.Wrapf(err, "failed to get repository %s/%s", bbc.workspaces[0], repoName)
			}
			converted, err := bbc.convertRepository(*repoDetails)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to convert repository %s", repoName)
			}
			repositories = append(repositories, *converted)
		}
		return repositories, nil
	}

	// Otherwise, fetch all repositories in the workspace
	repoOptions := &bitbucket.RepositoriesOptions{
		Role:  "member", // Role filter is important for accessibility
		Owner: bbc.workspaces[0],
	}

	allReposResponse, err := bbc.bbClient.Repositories.ListForAccount(repoOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list repositories for account")
	}

	repositories := make([]scm.Repository, 0, len(allReposResponse.Items))
	for _, repo := range allReposResponse.Items {
		// No need to filter by bbc.repositories here, as this branch means it was empty
		converted, err := bbc.convertRepository(repo)
		if err != nil {
			// It might be desirable to log and continue, or collect errors
			return nil, errors.Wrapf(err, "failed to convert repository %s", repo.Name)
		}
		repositories = append(repositories, *converted)
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

		if bbc.authType == AuthTypeAppPassword {
			parsedURL.User = url.UserPassword(bbc.username, bbc.token)
		} else if bbc.authType == AuthTypeWorkspaceToken {
			parsedURL.User = url.UserPassword("x-token-auth", bbc.token)
		}

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

// AuthType defines the authentication method for Bitbucket Cloud
type AuthType int

const (
	// AuthTypeAppPassword will use app password authentication
	AuthTypeAppPassword AuthType = iota + 1
	// AuthTypeWorkspaceToken will use workspace token authentication
	AuthTypeWorkspaceToken
)

// ParseAuthType parses an auth type from a string
func ParseAuthType(str string) (AuthType, error) {
	switch str {
	default:
		return AuthType(0), fmt.Errorf("could not parse \"%s\" as auth type", str)
	case "app-password":
		return AuthTypeAppPassword, nil
	case "workspace-token":
		return AuthTypeWorkspaceToken, nil
	}
}

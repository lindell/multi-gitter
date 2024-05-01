package bitbucketserver

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const (
	cloneHTTPType = "http"
	cloneSSHType  = "ssh"
	stateMerged   = "MERGED"
	stateDeclined = "DECLINED"
)

// New create a new BitbucketServer client
func New(username, token, baseURL string, insecure bool, sshAuth bool, transportMiddleware func(http.RoundTripper) http.RoundTripper, repoListing RepositoryListing) (*BitbucketServer, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("token is empty")
	}

	if strings.TrimSpace(baseURL) == "" {
		return nil, errors.New("base url is empty")
	}

	bitbucketBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(bitbucketBaseURL.Path, "/rest") {
		bitbucketBaseURL.Path = path.Join(bitbucketBaseURL.Path, "/rest")
	}

	bitbucketServer := &BitbucketServer{}
	bitbucketServer.RepositoryListing = repoListing
	bitbucketServer.baseURL = bitbucketBaseURL
	bitbucketServer.username = username
	bitbucketServer.token = token
	bitbucketServer.sshAuth = sshAuth
	bitbucketServer.httpClient = &http.Client{
		Transport: transportMiddleware(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}, // nolint: gosec
		}),
	}
	bitbucketServer.config = bitbucketv1.NewConfiguration(bitbucketBaseURL.String(), func(config *bitbucketv1.Configuration) {
		config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
		config.HTTPClient = &http.Client{
			Transport: transportMiddleware(&http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}, // nolint: gosec
			}),
		}
	})

	return bitbucketServer, nil
}

func newClient(ctx context.Context, config *bitbucketv1.Configuration) *bitbucketv1.APIClient {
	return bitbucketv1.NewAPIClient(
		ctx,
		config,
	)
}

// BitbucketServer is a SCM instance of Bitbucket Server the on-prem version of Bitbucket
type BitbucketServer struct {
	RepositoryListing
	baseURL         *url.URL
	username, token string
	sshAuth         bool
	config          *bitbucketv1.Configuration
	httpClient      *http.Client
}

// RepositoryListing contains information about which repositories that should be fetched
type RepositoryListing struct {
	Projects     []string
	Users        []string
	Repositories []RepositoryReference
}

// RepositoryReference contains information to be able to reference a repository
type RepositoryReference struct {
	ProjectKey string
	Name       string
}

// ParseRepositoryReference parses a GiteaRepository reference from the format "projectKey/repoName"
func ParseRepositoryReference(val string) (RepositoryReference, error) {
	split := strings.Split(val, "/")
	if len(split) != 2 {
		return RepositoryReference{}, fmt.Errorf("could not parse repository reference: %s", val)
	}
	return RepositoryReference{
		ProjectKey: split[0],
		Name:       split[1],
	}, nil
}

// String returns the string representation of a repo reference
func (rr RepositoryReference) String() string {
	return fmt.Sprintf("%s/%s", rr.ProjectKey, rr.Name)
}

// GetRepositories Should get repositories based on the scm configuration
func (b *BitbucketServer) GetRepositories(ctx context.Context) ([]scm.Repository, error) {
	client := newClient(ctx, b.config)

	bitbucketRepositories, err := b.getRepositories(client)
	if err != nil {
		return nil, err
	}

	repositories := make([]scm.Repository, 0, len(bitbucketRepositories))

	// Get default branches and create repo interfaces
	for _, bitbucketRepository := range bitbucketRepositories {
		response, getDefaultBranchErr := client.DefaultApi.GetDefaultBranch(bitbucketRepository.Project.Key, bitbucketRepository.Slug)
		if getDefaultBranchErr != nil {
			return nil, getDefaultBranchErr
		}

		var defaultBranch bitbucketv1.Branch
		err = mapstructure.Decode(response.Values, &defaultBranch)
		if err != nil {
			return nil, err
		}

		repo, repoErr := b.convertRepository(bitbucketRepository, defaultBranch)
		if repoErr != nil {
			return nil, repoErr
		}

		repositories = append(repositories, *repo)
	}

	return repositories, nil
}

func (b *BitbucketServer) getRepositories(client *bitbucketv1.APIClient) ([]*bitbucketv1.Repository, error) {
	var bitbucketRepositories []*bitbucketv1.Repository

	for _, project := range b.Projects {
		repos, err := b.getProjectRepositories(client, project)
		if err != nil {
			return nil, err
		}

		bitbucketRepositories = append(bitbucketRepositories, repos...)
	}

	for _, user := range b.Users {
		repos, err := b.getProjectRepositories(client, user)
		if err != nil {
			return nil, err
		}

		bitbucketRepositories = append(bitbucketRepositories, repos...)
	}

	for _, repositoryRef := range b.Repositories {
		repo, err := b.getRepository(client, repositoryRef.ProjectKey, repositoryRef.Name)
		if err != nil {
			return nil, err
		}

		bitbucketRepositories = append(bitbucketRepositories, repo)
	}

	// Remove duplicate repos
	repositoryMap := make(map[int]*bitbucketv1.Repository, len(bitbucketRepositories))
	for _, bitbucketRepository := range bitbucketRepositories {
		repositoryMap[bitbucketRepository.ID] = bitbucketRepository
	}
	bitbucketRepositories = make([]*bitbucketv1.Repository, 0, len(repositoryMap))
	for _, repo := range repositoryMap {
		bitbucketRepositories = append(bitbucketRepositories, repo)
	}
	sort.Slice(bitbucketRepositories, func(i, j int) bool {
		return bitbucketRepositories[i].ID < bitbucketRepositories[j].ID
	})

	return bitbucketRepositories, nil
}

func (b *BitbucketServer) getRepository(client *bitbucketv1.APIClient, projectKey, repositorySlug string) (*bitbucketv1.Repository, error) {
	response, err := client.DefaultApi.GetRepository(projectKey, repositorySlug)
	if err != nil {
		return nil, err
	}

	var bitbucketRepository bitbucketv1.Repository
	err = mapstructure.Decode(response.Values, &bitbucketRepository)
	if err != nil {
		return nil, err
	}

	return &bitbucketRepository, nil
}

func (b *BitbucketServer) getProjectRepositories(client *bitbucketv1.APIClient, projectKey string) ([]*bitbucketv1.Repository, error) {
	params := map[string]interface{}{"start": 0, "limit": 25}

	var repositories []*bitbucketv1.Repository
	for {
		response, err := client.DefaultApi.GetRepositoriesWithOptions(projectKey, params)
		if err != nil {
			return nil, err
		}

		var pager bitbucketRepositoryPager
		err = mapstructure.Decode(response.Values, &pager)
		if err != nil {
			return nil, err
		}

		for _, repo := range pager.Values {
			r := repo
			repositories = append(repositories, &r)
		}

		if pager.IsLastPage {
			break
		}

		params["start"] = pager.NextPageStart
	}

	return repositories, nil
}

// CreatePullRequest Creates a pull request. The repo parameter will always originate from the same package
func (b *BitbucketServer) CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
	r := repo.(repository)
	prR := prRepo.(repository)

	client := newClient(ctx, b.config)

	reviewers, err := b.getUsersWithLinks(newPR.Reviewers, client)
	if err != nil {
		return nil, err
	}

	response, err := client.DefaultApi.CreatePullRequest(r.project, r.name, bitbucketv1.PullRequest{
		Title:       newPR.Title,
		Description: newPR.Body,
		Reviewers:   reviewers,
		FromRef: bitbucketv1.PullRequestRef{
			ID: fmt.Sprintf("refs/heads/%s", newPR.Head),
			Repository: bitbucketv1.Repository{
				Slug: prR.name,
				Project: &bitbucketv1.Project{
					Key: prR.project,
				},
			},
		},
		ToRef: bitbucketv1.PullRequestRef{
			ID: fmt.Sprintf("refs/heads/%s", newPR.Base),
			Repository: bitbucketv1.Repository{
				Slug: r.name,
				Project: &bitbucketv1.Project{
					Key: r.project,
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create pull request for repository %s: %s", r.name, err)
	}

	pullRequestResp, err := bitbucketv1.GetPullRequestResponse(response)
	if err != nil {
		return nil, fmt.Errorf("unable to create pull request for repository %s: %s", r.name, err)
	}

	return newPullRequest(pullRequestResp), nil
}

func (b *BitbucketServer) getUsersWithLinks(usernames []string, client *bitbucketv1.APIClient) ([]bitbucketv1.UserWithMetadata, error) {
	var usersWithMetadata []bitbucketv1.UserWithMetadata

	for _, username := range usernames {
		response, err := client.DefaultApi.GetUser(username)
		if err != nil {
			return nil, err
		}

		var userWithLinks bitbucketv1.UserWithLinks
		err = mapstructure.Decode(response.Values, &userWithLinks)
		if err != nil {
			return nil, err
		}

		usersWithMetadata = append(usersWithMetadata, bitbucketv1.UserWithMetadata{User: userWithLinks})
	}

	return usersWithMetadata, nil
}

// UpdatePullRequest updates an existing pull request
func (b *BitbucketServer) UpdatePullRequest(_ context.Context, _ scm.Repository, pullReq scm.PullRequest, _ scm.NewPullRequest) (scm.PullRequest, error) {
	// currently unsupported by gfleury/go-bitbucket-v1 see https://github.com/gfleury/go-bitbucket-v1/issues/66
	// for now just ignore the request rather than returning an error
	return pullReq, nil
}

// GetPullRequests Gets the latest pull requests from repositories based on the scm configuration
func (b *BitbucketServer) GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error) {
	client := newClient(ctx, b.config)

	repositories, err := b.getRepositories(client)
	if err != nil {
		return nil, err
	}

	var prs []scm.PullRequest
	for _, repo := range repositories {
		pr, getPullRequestErr := b.getPullRequest(client, branchName, repo.Project.Key, repo.Slug)
		if getPullRequestErr != nil {
			return nil, getPullRequestErr
		}
		if pr == nil {
			continue
		}

		convertedPR, err := b.convertPullRequest(client, repo.Project.Key, repo.Slug, branchName, pr)
		if err != nil {
			return nil, err
		}

		prs = append(prs, convertedPR)
	}

	return prs, nil
}

func (b *BitbucketServer) convertPullRequest(client *bitbucketv1.APIClient, project, repoName, branchName string, pr *bitbucketv1.PullRequest) (pullRequest, error) {
	status, err := b.pullRequestStatus(client, project, repoName, pr)
	if err != nil {
		return pullRequest{}, err
	}

	return pullRequest{
		repoName:   repoName,
		project:    project,
		branchName: branchName,
		prProject:  pr.FromRef.Repository.Project.Key,
		prRepoName: pr.FromRef.Repository.Slug,
		number:     pr.ID,
		version:    pr.Version,
		guiURL:     pr.Links.Self[0].Href,
		status:     status,
	}, nil
}

func (b *BitbucketServer) pullRequestStatus(client *bitbucketv1.APIClient, project, repoName string, pr *bitbucketv1.PullRequest) (scm.PullRequestStatus, error) {
	switch pr.State {
	case stateMerged:
		return scm.PullRequestStatusMerged, nil
	case stateDeclined:
		return scm.PullRequestStatusClosed, nil
	}

	response, err := client.DefaultApi.CanMerge(project, repoName, pr.ID)
	if err != nil {
		return scm.PullRequestStatusUnknown, err
	}

	var merge bitbucketv1.MergeGetResponse
	err = mapstructure.Decode(response.Values, &merge)
	if err != nil {
		return scm.PullRequestStatusUnknown, err
	}

	if !merge.CanMerge {
		return scm.PullRequestStatusPending, nil
	}

	if merge.Conflicted {
		return scm.PullRequestStatusError, nil
	}

	return scm.PullRequestStatusSuccess, nil
}

func (b *BitbucketServer) getPullRequest(client *bitbucketv1.APIClient, branchName, project, repoName string) (*bitbucketv1.PullRequest, error) {
	params := map[string]interface{}{"start": 0, "limit": 25}

	var pullRequests []bitbucketv1.PullRequest
	for {
		response, err := client.DefaultApi.GetPullRequestsPage(project, repoName, params)
		if err != nil {
			return nil, err
		}

		var pager bitbucketPullRequestPager
		err = mapstructure.Decode(response.Values, &pager)
		if err != nil {
			return nil, err
		}

		pullRequests = append(pullRequests, pager.Values...)

		if pager.IsLastPage {
			break
		}

		params["start"] = pager.NextPageStart
	}

	for _, pr := range pullRequests {
		if pr.FromRef.DisplayID == branchName {
			return &pr, nil
		}
	}

	return nil, nil
}

// GetOpenPullRequest gets a pull request for one specific repository
func (b *BitbucketServer) GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
	r := repo.(repository)

	client := newClient(ctx, b.config)

	pr, err := b.getPullRequest(client, branchName, r.project, r.name)
	if err != nil {
		return nil, err
	}

	if pr == nil {
		return nil, nil
	}

	convertedPR, err := b.convertPullRequest(client, r.project, r.name, branchName, pr)
	if err != nil {
		return nil, err
	}

	return convertedPR, nil
}

// MergePullRequest Merges a pull request, the pr parameter will always originate from the same package
func (b *BitbucketServer) MergePullRequest(ctx context.Context, pr scm.PullRequest) error {
	bitbucketPR := pr.(pullRequest)

	client := newClient(ctx, b.config)

	response, err := client.DefaultApi.GetPullRequest(bitbucketPR.project, bitbucketPR.repoName, bitbucketPR.number)
	if err != nil {
		if strings.Contains(err.Error(), "com.atlassian.bitbucket.pull.NoSuchPullRequestException") {
			return nil
		}
		return err
	}

	pullRequestResponse, err := bitbucketv1.GetPullRequestResponse(response)
	if err != nil {
		return err
	}

	if !pullRequestResponse.Open {
		return nil
	}

	mergeMap := make(map[string]interface{})
	mergeMap["version"] = pullRequestResponse.Version

	_, err = client.DefaultApi.Merge(bitbucketPR.project, bitbucketPR.repoName, bitbucketPR.number, mergeMap, nil, []string{"application/json"})
	if err != nil {
		return err
	}

	return b.deleteBranch(ctx, bitbucketPR)
}

// ClosePullRequest Close a pull request, the pr parameter will always originate from the same package
func (b *BitbucketServer) ClosePullRequest(ctx context.Context, pr scm.PullRequest) error {
	bitbucketPR := pr.(pullRequest)

	client := newClient(ctx, b.config)

	_, err := client.DefaultApi.DeleteWithVersion(bitbucketPR.project, bitbucketPR.repoName, bitbucketPR.number, int(bitbucketPR.version))
	if err != nil {
		return err
	}

	return b.deleteBranch(ctx, bitbucketPR)
}

func (b *BitbucketServer) deleteBranch(ctx context.Context, pr pullRequest) error {
	urlPath := *b.baseURL
	urlPath.Path = path.Join(urlPath.Path, "branch-utils/1.0/projects", pr.project, "repos", pr.repoName, "branches")

	body := bitbucketDeleteBranch{Name: path.Join("refs", "heads", pr.branchName), DryRun: false}
	bodyBytes, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, "DELETE", urlPath.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add(
		"Authorization",
		"Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", b.username, b.token))),
	)

	response, err := b.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		buf := new(bytes.Buffer)
		_, readFromErr := buf.ReadFrom(response.Body)
		if readFromErr != nil {
			return readFromErr
		}

		return errors.Errorf("unable to delete branch: status code %d: %s", response.StatusCode, buf.String())
	}

	return nil
}

// ForkRepository forks a repository. If newOwner is set, use it, otherwise fork to the current user
func (b *BitbucketServer) ForkRepository(_ context.Context, _ scm.Repository, _ string) (scm.Repository, error) {
	return nil, errors.New("forking not implemented for bitbucket server")
}

type bitbucketRepositoryPager struct {
	Size          int                      `json:"size"`
	Limit         int                      `json:"limit"`
	Start         int                      `json:"start"`
	NextPageStart int                      `json:"nextPageStart"`
	IsLastPage    bool                     `json:"isLastPage"`
	Values        []bitbucketv1.Repository `json:"values"`
}

type bitbucketPullRequestPager struct {
	Size          int                       `json:"size"`
	Limit         int                       `json:"limit"`
	Start         int                       `json:"start"`
	NextPageStart int                       `json:"nextPageStart"`
	IsLastPage    bool                      `json:"isLastPage"`
	Values        []bitbucketv1.PullRequest `json:"values"`
}

type bitbucketDeleteBranch struct {
	Name   string `json:"name"`
	DryRun bool   `json:"dryRun"`
}

package bitbucketserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/lindell/multi-gitter/internal/git"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const (
	cloneType     = "http"
	stateMerged   = "MERGED"
	stateDeclined = "DECLINED"
)

// New create a new BitbucketServer client
func New(username, token, baseURL string, insecure bool, repoListing RepositoryListing) (*BitbucketServer, error) {
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
	bitbucketServer.username = username
	bitbucketServer.token = token
	bitbucketServer.insecure = insecure
	bitbucketServer.baseURL = bitbucketBaseURL

	return bitbucketServer, nil
}

func newClient(ctx context.Context, token, baseURL string, insecure bool) *bitbucketv1.APIClient {
	return bitbucketv1.NewAPIClient(
		ctx,
		bitbucketv1.NewConfiguration(baseURL, func(config *bitbucketv1.Configuration) {
			config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
			config.HTTPClient = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}, // nolint: gosec
				},
			}
		}),
	)
}

// BitbucketServer is a SCM instance of Bitbucket Server the on-prem version of Bitbucket
type BitbucketServer struct {
	RepositoryListing
	username, token string
	insecure        bool
	baseURL         *url.URL
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
func (b *BitbucketServer) GetRepositories(ctx context.Context) ([]git.Repository, error) {
	client := newClient(ctx, b.token, b.baseURL.String(), b.insecure)

	bitbucketRepositories, err := b.getRepositories(client)
	if err != nil {
		return nil, err
	}

	repositories := make([]git.Repository, 0, len(bitbucketRepositories))

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

		repo, repoErr := convertRepository(bitbucketRepository, defaultBranch, b.username, b.token)
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
func (b *BitbucketServer) CreatePullRequest(ctx context.Context, repo git.Repository, prRepo git.Repository, newPR git.NewPullRequest) (git.PullRequest, error) {
	r := repo.(repository)
	prR := prRepo.(repository)

	client := newClient(ctx, b.token, b.baseURL.String(), b.insecure)

	var usersWithMetadata []bitbucketv1.UserWithMetadata
	for _, reviewer := range newPR.Reviewers {
		response, err := client.DefaultApi.GetUser(reviewer)
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

	response, err := client.DefaultApi.CreatePullRequest(r.project, r.name, bitbucketv1.PullRequest{
		Title:       newPR.Title,
		Description: newPR.Body,
		Reviewers:   usersWithMetadata,
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

// GetPullRequests Gets the latest pull requests from repositories based on the scm configuration
func (b *BitbucketServer) GetPullRequests(ctx context.Context, branchName string) ([]git.PullRequest, error) {
	client := newClient(ctx, b.token, b.baseURL.String(), b.insecure)

	repositories, err := b.getRepositories(client)
	if err != nil {
		return nil, err
	}

	var prs []git.PullRequest
	for _, repo := range repositories {
		pr, getPullRequestErr := b.getPullRequest(client, branchName, repo)
		if getPullRequestErr != nil {
			return nil, getPullRequestErr
		}
		if pr == nil {
			continue
		}

		status, pullRequestStatusErr := b.pullRequestStatus(client, repo, pr)
		if pullRequestStatusErr != nil {
			return nil, pullRequestStatusErr
		}

		prs = append(prs, pullRequest{
			repoName:   repo.Slug,
			project:    repo.Project.Key,
			branchName: branchName,
			prProject:  pr.FromRef.Repository.Project.Key,
			prRepoName: pr.FromRef.Repository.Slug,
			status:     status,
			guiURL:     pr.Links.Self[0].Href,
		})
	}

	return prs, nil
}

func (b *BitbucketServer) pullRequestStatus(client *bitbucketv1.APIClient, repo *bitbucketv1.Repository, pr *bitbucketv1.PullRequest) (git.PullRequestStatus, error) {
	switch pr.State {
	case stateMerged:
		return git.PullRequestStatusMerged, nil
	case stateDeclined:
		return git.PullRequestStatusClosed, nil
	}

	response, err := client.DefaultApi.CanMerge(repo.Project.Key, repo.Slug, int64(pr.ID))
	if err != nil {
		return git.PullRequestStatusUnknown, err
	}

	var merge bitbucketv1.MergeGetResponse
	err = mapstructure.Decode(response.Values, &merge)
	if err != nil {
		return git.PullRequestStatusUnknown, err
	}

	if !merge.CanMerge {
		return git.PullRequestStatusPending, nil
	}

	if merge.Conflicted {
		return git.PullRequestStatusError, nil
	}

	return git.PullRequestStatusUnknown, nil
}

func (b *BitbucketServer) getPullRequest(client *bitbucketv1.APIClient, branchName string, repo *bitbucketv1.Repository) (*bitbucketv1.PullRequest, error) {
	params := map[string]interface{}{"start": 0, "limit": 25}

	var pullRequests []bitbucketv1.PullRequest
	for {
		response, err := client.DefaultApi.GetPullRequestsPage(repo.Project.Key, repo.Slug, params)
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

	return nil, errors.Errorf("could not find pull request in repository %s for branch %s", repo.Name, branchName)
}

// MergePullRequest Merges a pull request, the pr parameter will always originate from the same package
func (b *BitbucketServer) MergePullRequest(ctx context.Context, pr git.PullRequest) error {
	bitbucketPR := pr.(pullRequest)

	client := newClient(ctx, b.token, b.baseURL.String(), b.insecure)

	response, err := client.DefaultApi.GetPullRequest(bitbucketPR.project, bitbucketPR.repoName, bitbucketPR.number)
	if err != nil {
		if strings.Contains(err.Error(), "com.atlassian.bitbucket.pull.NoSuchPullRequestException") {
			return nil
		}
		return err
	}

	pullRequest, err := bitbucketv1.GetPullRequestResponse(response)
	if err != nil {
		return err
	}

	if !pullRequest.Open {
		return nil
	}

	mergeMap := make(map[string]interface{})
	mergeMap["version"] = pullRequest.Version

	_, err = client.DefaultApi.Merge(bitbucketPR.project, bitbucketPR.repoName, bitbucketPR.number, mergeMap, nil, []string{"application/json"})
	if err != nil {
		return err
	}

	return nil
}

// ClosePullRequest Close a pull request, the pr parameter will always originate from the same package
func (b *BitbucketServer) ClosePullRequest(ctx context.Context, pr git.PullRequest) error {
	bitbucketPR := pr.(pullRequest)

	client := newClient(ctx, b.token, b.baseURL.String(), b.insecure)

	_, err := client.DefaultApi.DeletePullRequest(bitbucketPR.project, bitbucketPR.repoName, int64(bitbucketPR.number))

	return err
}

// ForkRepository forks a repository. If newOwner is set, use it, otherwise fork to the current user
func (b *BitbucketServer) ForkRepository(_ context.Context, repo git.Repository, newOwner string) (git.Repository, error) {
	return nil, errors.New("forking not implemented")
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

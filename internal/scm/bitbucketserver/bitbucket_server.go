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
	"github.com/lindell/multi-gitter/internal/pullrequest"
	"github.com/lindell/multi-gitter/internal/repository"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const (
	cloneType     = "http"
	stateMerged   = "MERGED"
	stateOpen     = "OPEN"
	stateDeclined = "DECLINED"
)

// New create a new BitbucketServer client
func New(ctx context.Context, token, baseURL string, insecure bool, repoListing RepositoryListing) (*BitbucketServer, error) {
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
	bitbucketServer.client = bitbucketv1.NewAPIClient(
		ctx,
		bitbucketv1.NewConfiguration(bitbucketBaseURL.String(), func(config *bitbucketv1.Configuration) {
			config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
			config.HTTPClient = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}, // nolint: gosec
				},
			}
		}),
	)

	return bitbucketServer, nil
}

type BitbucketServer struct {
	RepositoryListing
	client *bitbucketv1.APIClient
}

// RepositoryListing contains information about which repositories that should be fetched
type RepositoryListing struct {
	Projects     []string
	Users        []string
	Repositories []RepositoryReference
}

// RepositoryReference contains information to be able to reference a Repository
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
func (b *BitbucketServer) GetRepositories(_ context.Context) ([]repository.Data, error) {
	bitbucketRepositories, err := b.getRepositories()
	if err != nil {
		return nil, err
	}

	repositories := make([]repository.Data, 0, len(bitbucketRepositories))

	// Get default branches and create repo interfaces
	for _, bitbucketRepository := range bitbucketRepositories {
		response, getDefaultBranchErr := b.client.DefaultApi.GetDefaultBranch(bitbucketRepository.Project.Key, bitbucketRepository.Slug)
		if getDefaultBranchErr != nil {
			return nil, getDefaultBranchErr
		}

		var defaultBranch bitbucketv1.Branch
		err = mapstructure.Decode(response.Values, &defaultBranch)
		if err != nil {
			return nil, err
		}

		repo, repoErr := newRepo(bitbucketRepository, defaultBranch)
		if repoErr != nil {
			return nil, repoErr
		}

		repositories = append(repositories, *repo)
	}

	return repositories, nil
}

func (b *BitbucketServer) getRepositories() ([]*bitbucketv1.Repository, error) {
	var bitbucketRepositories []*bitbucketv1.Repository

	for _, project := range b.Projects {
		repos, err := b.getProjectRepositories(project)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			bitbucketRepositories = append(bitbucketRepositories, repo)
		}
	}

	for _, user := range b.Users {
		repos, err := b.getProjectRepositories(user)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			bitbucketRepositories = append(bitbucketRepositories, repo)
		}
	}

	for _, repositoryRef := range b.Repositories {
		repo, err := b.getRepository(repositoryRef.ProjectKey, repositoryRef.Name)
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

func (b *BitbucketServer) getRepository(projectKey, repositorySlug string) (*bitbucketv1.Repository, error) {
	response, err := b.client.DefaultApi.GetRepository(projectKey, repositorySlug)
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

func (b *BitbucketServer) getProjectRepositories(projectKey string) ([]*bitbucketv1.Repository, error) {
	params := map[string]interface{}{"start": 0, "limit": 25}

	var repositories []*bitbucketv1.Repository
	for {
		response, err := b.client.DefaultApi.GetRepositoriesWithOptions(projectKey, params)
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
func (b *BitbucketServer) CreatePullRequest(_ context.Context, repo repository.Data, prRepo repository.Data, newPR pullrequest.NewPullRequest) (pullrequest.PullRequest, error) {
	r := repo.(Repository)
	prR := prRepo.(Repository)

	var usersWithMetadata []bitbucketv1.UserWithMetadata
	for _, reviewer := range newPR.Reviewers {
		response, err := b.client.DefaultApi.GetUser(reviewer)
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

	response, err := b.client.DefaultApi.CreatePullRequest(r.project, r.name, bitbucketv1.PullRequest{
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
		return nil, fmt.Errorf("unable to create pull request for Repository %s: %s", r.name, err)
	}

	pullRequestResp, err := bitbucketv1.GetPullRequestResponse(response)
	if err != nil {
		return nil, fmt.Errorf("unable to create pull request for Repository %s: %s", r.name, err)
	}

	return newPullRequest(pullRequestResp), nil
}

// GetPullRequests Gets the latest pull requests from repositories based on the scm configuration
func (b *BitbucketServer) GetPullRequests(_ context.Context, branchName string) ([]pullrequest.PullRequest, error) {
	repositories, err := b.getRepositories()
	if err != nil {
		return nil, err
	}

	var prs []pullrequest.PullRequest
	for _, repo := range repositories {
		pr, getPullRequestErr := b.getPullRequest(branchName, repo)
		if getPullRequestErr != nil {
			return nil, getPullRequestErr
		}
		if pr == nil {
			continue
		}

		status, pullRequestStatusErr := b.pullRequestStatus(repo, pr)
		if pullRequestStatusErr != nil {
			return nil, pullRequestStatusErr
		}

		prs = append(prs, PullRequest{
			repoName:   repo.Slug,
			project:    repo.Project.Key,
			branchName: branchName,
			prProject:  pr.FromRef.Repository.Project.Key,
			prRepoName: pr.FromRef.Repository.Slug,
			status:     status,
			guiURL:     pr.Links.Self[0].Href,
		})
	}

	return nil, nil
}

func (b *BitbucketServer) pullRequestStatus(repo *bitbucketv1.Repository, pr *bitbucketv1.PullRequest) (pullrequest.Status, error) {
	switch pr.State {
	case stateMerged:
		return pullrequest.StatusMerged, nil
	case stateDeclined:
		return pullrequest.StatusClosed, nil
	}

	response, err := b.client.DefaultApi.CanMerge(repo.Project.Key, repo.Slug, int64(pr.ID))
	if err != nil {
		return pullrequest.StatusUnknown, err
	}

	var merge bitbucketv1.MergeGetResponse
	err = mapstructure.Decode(response.Values, &merge)
	if err != nil {
		return pullrequest.StatusUnknown, err
	}

	if !merge.CanMerge {
		return pullrequest.StatusPending, nil
	}

	if merge.Conflicted {
		return pullrequest.StatusError, nil
	}

	return pullrequest.StatusUnknown, nil
}

func (b *BitbucketServer) getPullRequest(branchName string, repo *bitbucketv1.Repository) (*bitbucketv1.PullRequest, error) {
	params := map[string]interface{}{"start": 0, "limit": 25}

	var pullRequests []bitbucketv1.PullRequest
	for {
		response, err := b.client.DefaultApi.GetPullRequestsPage(repo.Project.Key, repo.Slug, params)
		if err != nil {
			return nil, err
		}

		var pager bitbucketPullRequestPager
		err = mapstructure.Decode(response.Values, &pager)
		if err != nil {
			return nil, err
		}

		for _, pr := range pager.Values {
			pullRequests = append(pullRequests, pr)
		}

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

	return nil, errors.Errorf("could not find pull request in Repository %s for branch %s", repo.Name, branchName)
}

// MergePullRequest Merges a pull request, the pr parameter will always originate from the same package
func (b *BitbucketServer) MergePullRequest(_ context.Context, pr pullrequest.PullRequest) error {
	bitbucketPR := pr.(PullRequest)

	response, err := b.client.DefaultApi.GetPullRequest(bitbucketPR.project, bitbucketPR.repoName, bitbucketPR.number)
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

	_, err = b.client.DefaultApi.Merge(bitbucketPR.project, bitbucketPR.repoName, bitbucketPR.number, mergeMap, nil, []string{"application/json"})
	if err != nil {
		return err
	}

	return nil
}

// ClosePullRequest Close a pull request, the pr parameter will always originate from the same package
func (b *BitbucketServer) ClosePullRequest(_ context.Context, pr pullrequest.PullRequest) error {
	bitbucketPR := pr.(PullRequest)

	_, err := b.client.DefaultApi.DeletePullRequest(bitbucketPR.project, bitbucketPR.repoName, int64(bitbucketPR.number))

	return err
}

// ForkRepository forks a repository. If newOwner is set, use it, otherwise fork to the current user
func (b *BitbucketServer) ForkRepository(_ context.Context, repo repository.Data, newOwner string) (repository.Data, error) {
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

func newRepo(bitbucketRepository *bitbucketv1.Repository, defaultBranch bitbucketv1.Branch) (*Repository, error) {
	var cloneURL *url.URL
	var err error
	for _, clone := range bitbucketRepository.Links.Clone {
		if strings.EqualFold(clone.Name, cloneType) {
			cloneURL, err = url.Parse(clone.Href)
			if err != nil {
				return nil, err
			}

			break
		}
	}

	if cloneURL == nil {
		return nil, errors.Errorf("unable to find clone url for repostory %s using clone type %s", bitbucketRepository.Name, cloneType)
	}

	repo := Repository{
		name:          bitbucketRepository.Slug,
		project:       bitbucketRepository.Project.Key,
		defaultBranch: defaultBranch.DisplayID,
		cloneURL:      cloneURL,
	}

	return &repo, nil
}

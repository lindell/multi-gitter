package gitea

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/lindell/multi-gitter/internal/pullrequest"
	"github.com/lindell/multi-gitter/internal/repository"
	"github.com/pkg/errors"

	internalHTTP "github.com/lindell/multi-gitter/internal/http"
)

// New create a new Gitea client
func New(token, baseURL string, repoListing RepositoryListing, mergeTypes []pullrequest.MergeType) (*Gitea, error) {
	gitea := &Gitea{
		RepositoryListing: repoListing,

		baseURL: baseURL,
		token:   token,

		MergeTypes: mergeTypes,
	}

	// Initialize the gitea client to ensure no error will occur when running a function
	_, err := gitea.giteaClientErr(context.Background())

	return gitea, err
}

func (g *Gitea) giteaClientErr(ctx context.Context) (*gitea.Client, error) {
	client, err := gitea.NewClient(
		g.baseURL,
		gitea.SetHTTPClient(&http.Client{
			Transport: internalHTTP.LoggingRoundTripper{},
		}),
		gitea.SetToken(g.token),
		gitea.SetContext(ctx),
	)
	return client, err
}

func (g *Gitea) giteaClient(ctx context.Context) *gitea.Client {
	client, _ := g.giteaClientErr(ctx)
	return client
}

// Gitea contain Gitea configuration
type Gitea struct {
	RepositoryListing

	baseURL string
	token   string

	currentUser *gitea.User

	MergeTypes []pullrequest.MergeType
}

// RepositoryListing contains information about which repositories that should be fetched
type RepositoryListing struct {
	Organizations []string
	Users         []string
	Repositories  []RepositoryReference
}

// RepositoryReference contains information to be able to reference a repository
type RepositoryReference struct {
	OwnerName string
	Name      string
}

// ParseRepositoryReference parses a repository reference from the format "ownerName/repoName"
func ParseRepositoryReference(val string) (RepositoryReference, error) {
	split := strings.Split(val, "/")
	if len(split) != 2 {
		return RepositoryReference{}, fmt.Errorf("could not parse repository reference: %s", val)
	}
	return RepositoryReference{
		OwnerName: split[0],
		Name:      split[1],
	}, nil
}

// GetRepositories fetches repositories from all sources (groups/user/specific repo)
func (g *Gitea) GetRepositories(ctx context.Context) ([]repository.Data, error) {
	allRepos, err := g.getRepositories(ctx)
	if err != nil {
		return nil, err
	}

	repos := make([]repository.Data, 0, len(allRepos))
	for _, repo := range allRepos {
		convertedRepo, err := convertRepository(repo)
		if err != nil {
			return nil, err
		}
		repos = append(repos, convertedRepo)
	}

	return repos, nil
}

func (g *Gitea) getRepositories(ctx context.Context) ([]*gitea.Repository, error) {
	allRepos := []*gitea.Repository{}

	for _, group := range g.Organizations {
		repos, err := g.getGroupRepositories(ctx, group)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
	}

	for _, user := range g.Users {
		repos, err := g.getUserRepositories(ctx, user)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
	}

	for _, repo := range g.Repositories {
		repo, err := g.getRepository(ctx, repo)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repo)
	}

	// Remove duplicate repos
	repoMap := map[int64]*gitea.Repository{}
	for _, repo := range allRepos {
		repoMap[repo.ID] = repo
	}
	allRepos = make([]*gitea.Repository, 0, len(repoMap))
	for _, repo := range repoMap {
		allRepos = append(allRepos, repo)
	}
	sort.Slice(allRepos, func(i, j int) bool {
		return allRepos[i].ID < allRepos[j].ID
	})

	return allRepos, nil
}

func (g *Gitea) getGroupRepositories(ctx context.Context, groupName string) ([]*gitea.Repository, error) {
	var allRepos []*gitea.Repository
	for i := 1; ; i++ {
		repos, _, err := g.giteaClient(ctx).ListOrgRepos(groupName, gitea.ListOrgReposOptions{
			ListOptions: gitea.ListOptions{
				Page:     i,
				PageSize: 100,
			},
		})
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if len(repos) < 100 {
			break
		}
	}
	return allRepos, nil
}

func (g *Gitea) getRepository(ctx context.Context, repoRef RepositoryReference) (*gitea.Repository, error) {
	repo, _, err := g.giteaClient(ctx).GetRepo(repoRef.OwnerName, repoRef.Name)
	if err != nil {
		return nil, err
	}
	return repo, err
}

func (g *Gitea) getUserRepositories(ctx context.Context, username string) ([]*gitea.Repository, error) {
	var allRepos []*gitea.Repository
	for i := 1; ; i++ {
		repos, _, err := g.giteaClient(ctx).ListUserRepos(username, gitea.ListReposOptions{
			ListOptions: gitea.ListOptions{
				Page:     i,
				PageSize: 100,
			},
		})
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if len(repos) < 100 {
			break
		}
	}
	return allRepos, nil
}

// CreatePullRequest creates a pull request
func (g *Gitea) CreatePullRequest(ctx context.Context, repo repository.Data, prRepo repository.Data, newPR pullrequest.NewPullRequest) (pullrequest.PullRequest, error) {
	r := repo.(Repository)
	prR := prRepo.(Repository)

	head := fmt.Sprintf("%s:%s", prR.ownerName, newPR.Head)

	pr, _, err := g.giteaClient(ctx).CreatePullRequest(r.ownerName, r.name, gitea.CreatePullRequestOption{
		Head:  head,
		Base:  newPR.Base,
		Title: newPR.Title,
		Body:  newPR.Body,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not create pull request")
	}

	_, err = g.giteaClient(ctx).CreateReviewRequests(r.ownerName, r.name, pr.Index, gitea.PullReviewRequestOptions{
		Reviewers: newPR.Reviewers,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not add reviewer to pull request")
	}

	return PullRequest{
		repoName:    r.name,
		ownerName:   r.ownerName,
		branchName:  newPR.Head,
		prOwnerName: pr.Head.Repository.Owner.UserName,
		prRepoName:  pr.Head.Repository.Name,
		index:       pr.Index,
		webURL:      pr.HTMLURL,
	}, nil
}

// GetPullRequests gets all pull requests of with a specific branch
func (g *Gitea) GetPullRequests(ctx context.Context, branchName string) ([]pullrequest.PullRequest, error) {
	repos, err := g.getRepositories(ctx)
	if err != nil {
		return nil, err
	}

	prs := []pullrequest.PullRequest{}
	for _, repo := range repos {
		pr, err := g.getPullRequest(ctx, branchName, repo)
		if err != nil {
			return nil, err
		}
		if pr == nil {
			continue
		}

		status, err := g.pullRequestStatus(ctx, repo, pr)
		if err != nil {
			return nil, err
		}

		prs = append(prs, PullRequest{
			repoName:    repo.Name,
			ownerName:   repo.Owner.UserName,
			branchName:  branchName,
			prOwnerName: pr.Head.Repository.Owner.UserName,
			prRepoName:  pr.Head.Repository.Name,
			status:      status,
			index:       pr.Index,
			webURL:      pr.HTMLURL,
		})
	}

	return prs, nil
}

func (g *Gitea) getPullRequest(ctx context.Context, branchName string, repo *gitea.Repository) (*gitea.PullRequest, error) {
	// We would like to be able to search for a pr with a specific head here, but current (2021-04-24), that option does not exist in the API
	prs, _, err := g.giteaClient(ctx).ListRepoPullRequests(repo.Owner.UserName, repo.Name, gitea.ListPullRequestsOptions{
		State: "all",
		Sort:  "recentupdate",
	})
	if err != nil {
		return nil, err
	}

	for _, pr := range prs {
		if pr.Head.Name == branchName {
			return pr, nil
		}
	}
	return nil, nil
}

func (g *Gitea) pullRequestStatus(ctx context.Context, repo *gitea.Repository, pr *gitea.PullRequest) (pullrequest.Status, error) {
	if pr.Merged != nil {
		return pullrequest.StatusMerged, nil
	}

	if pr.State == gitea.StateClosed {
		return pullrequest.StatusClosed, nil
	}

	status, _, err := g.giteaClient(ctx).GetCombinedStatus(repo.Owner.UserName, repo.Name, pr.Head.Sha)
	if err != nil {
		return pullrequest.StatusUnknown, err
	}

	if len(status.Statuses) == 0 {
		return pullrequest.StatusSuccess, nil
	}

	switch status.State {
	case gitea.StatusPending:
		return pullrequest.StatusPending, nil
	case gitea.StatusSuccess:
		return pullrequest.StatusSuccess, nil
	case gitea.StatusError, gitea.StatusFailure:
		return pullrequest.StatusError, nil
	}

	return pullrequest.StatusUnknown, nil
}

// MergePullRequest merges a pull request
func (g *Gitea) MergePullRequest(ctx context.Context, pullReq pullrequest.PullRequest) error {
	pr := pullReq.(PullRequest)

	repo, _, err := g.giteaClient(ctx).GetRepo(pr.ownerName, pr.repoName)
	if err != nil {
		return errors.Wrapf(err, "could not fetch %s/%s repository", pr.ownerName, pr.repoName)
	}

	// Filter out all merge types to only the allowed ones, but keep the order of the ones left
	mergeTypes := pullrequest.MergeTypeIntersection(g.MergeTypes, repoMergeTypes(repo))
	if len(mergeTypes) == 0 {
		return errors.New("none of the configured merge types was permitted")
	}

	merged, _, err := g.giteaClient(ctx).MergePullRequest(pr.ownerName, pr.repoName, pr.index, gitea.MergePullRequestOption{
		Style: mergeTypeGiteaName[mergeTypes[0]],
	})
	if err != nil {
		return errors.Wrapf(err, "could not merge %s/%s#%d", pr.ownerName, pr.repoName, pr.index)
	}

	if !merged {
		return errors.Errorf("could not merge %s/%s#%d", pr.ownerName, pr.repoName, pr.index)
	}

	deleted, _, err := g.giteaClient(ctx).DeleteRepoBranch(pr.prOwnerName, pr.prRepoName, pr.branchName)
	if err != nil {
		return errors.Wrapf(err, "could not delete branch after merging %s/%s", pr.ownerName, pr.repoName)
	}

	if !deleted {
		return errors.Errorf("could not delete branch after merging %s/%s", pr.ownerName, pr.repoName)
	}

	return nil
}

// ClosePullRequest closes a pull request
func (g *Gitea) ClosePullRequest(ctx context.Context, pullReq pullrequest.PullRequest) error {
	pr := pullReq.(PullRequest)

	state := gitea.StateClosed
	_, _, err := g.giteaClient(ctx).EditPullRequest(pr.ownerName, pr.repoName, pr.index, gitea.EditPullRequestOption{
		State: &state,
	})
	if err != nil {
		return errors.Wrapf(err, "could not close %s/%s#%d", pr.ownerName, pr.repoName, pr.index)
	}

	deleted, _, err := g.giteaClient(ctx).DeleteRepoBranch(pr.prOwnerName, pr.prRepoName, pr.branchName)
	if err != nil {
		return errors.Wrapf(err, "could not delete branch after merging %s/%s", pr.ownerName, pr.repoName)
	}

	if !deleted {
		return errors.Errorf("could not delete branch after merging %s/%s", pr.ownerName, pr.repoName)
	}

	return nil
}

// ForkRepository forks a GiteaRepository. If newOwner is empty, fork on the logged in user
func (g *Gitea) ForkRepository(ctx context.Context, repo repository.Data, newOwner string) (repository.Data, error) {
	r := repo.(Repository)

	forkTo := newOwner
	if forkTo == "" {
		user, err := g.getUser(ctx)
		if err != nil {
			return nil, err
		}
		forkTo = user.UserName
	}

	existingRepo, _, err := g.giteaClient(ctx).GetRepo(forkTo, r.name)
	if err == nil { // NB!
		return convertRepository(existingRepo)
	}

	forkOptions := gitea.CreateForkOption{}
	if newOwner != "" {
		forkOptions.Organization = &newOwner
	}

	createdRepo, _, err := g.giteaClient(ctx).CreateFork(r.ownerName, r.name, forkOptions)
	if err != nil {
		return nil, err
	}

	return convertRepository(createdRepo)
}

func (g *Gitea) getUser(ctx context.Context) (*gitea.User, error) {
	if g.currentUser != nil {
		return g.currentUser, nil
	}

	user, _, err := g.giteaClient(ctx).GetMyUserInfo()
	if err != nil {
		return nil, err
	}

	g.currentUser = user
	return user, nil
}

func convertRepository(repo *gitea.Repository) (Repository, error) {
	u, err := url.Parse(repo.CloneURL)
	if err != nil {
		return Repository{}, err
	}

	return Repository{
		url:           *u,
		name:          repo.Name,
		ownerName:     repo.Owner.UserName,
		defaultBranch: repo.DefaultBranch,
	}, nil
}

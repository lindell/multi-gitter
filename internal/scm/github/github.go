package github

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v58/github"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"golang.org/x/oauth2"

	"github.com/lindell/multi-gitter/internal/scm"
)

type Config struct {
	Token               string
	BaseURL             string
	TransportMiddleware func(http.RoundTripper) http.RoundTripper
	RepoListing         RepositoryListing
	MergeTypes          []scm.MergeType
	ForkMode            bool
	ForkOwner           string
	SSHAuth             bool
	ReadOnly            bool
	CheckPermissions    bool
}

// New create a new Github client
func New(
	config Config,
) (*Github, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	tc.Transport = config.TransportMiddleware(tc.Transport)

	var client *github.Client
	if config.BaseURL != "" {
		var err error
		client, err = github.NewEnterpriseClient(config.BaseURL, "", tc)
		if err != nil {
			return nil, err
		}
	} else {
		client = github.NewClient(tc)
	}

	return &Github{
		RepositoryListing: config.RepoListing,
		MergeTypes:        config.MergeTypes,
		token:             config.Token,
		baseURL:           config.BaseURL,
		Fork:              config.ForkMode,
		ForkOwner:         config.ForkOwner,
		SSHAuth:           config.SSHAuth,
		ghClient:          client,
		ReadOnly:          config.ReadOnly,
		checkPermissions:  config.CheckPermissions,
		httpClient: &http.Client{
			Transport: config.TransportMiddleware(http.DefaultTransport),
		},
	}, nil
}

// Github contain github configuration
type Github struct {
	RepositoryListing
	MergeTypes []scm.MergeType
	token      string
	baseURL    string

	// This determines if forks will be used when creating a prs.
	// In this package, it mainly determines which repos are possible to make changes on
	Fork bool

	// This determines when we are running in read only mode.
	ReadOnly bool

	// If set, the fork will happen to the ForkOwner value, and not the logged in user
	ForkOwner string

	// If set, use the SSH clone url instead of http(s)
	SSHAuth bool

	ghClient   *github.Client
	httpClient *http.Client

	// Caching of the logged in user
	user      string
	userMutex sync.Mutex

	// Used to make sure request that modifies state does not happen to often
	modMutex       sync.Mutex
	lastModRequest time.Time

	// If set, the permissions of each repository is checked before using it
	checkPermissions bool
}

// RepositoryListing contains information about which repositories that should be fetched
type RepositoryListing struct {
	Organizations    []string
	Users            []string
	Repositories     []RepositoryReference
	RepositorySearch string
	CodeSearch       string
	Topics           []string
	SkipForks        bool
}

// RepositoryReference contains information to be able to reference a repository
type RepositoryReference struct {
	OwnerName string
	Name      string
}

// String returns the string representation of a repo reference
func (rr RepositoryReference) String() string {
	return fmt.Sprintf("%s/%s", rr.OwnerName, rr.Name)
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

// GetRepositories fetches repositories from all sources (orgs/user/specific repo)
func (g *Github) GetRepositories(ctx context.Context) ([]scm.Repository, error) {
	allRepos, err := g.getRepositories(ctx)
	if err != nil {
		return nil, err
	}

	repos := make([]scm.Repository, 0, len(allRepos))
	for _, r := range allRepos {
		log := log.WithField("repo", r.GetFullName())
		permissions := r.GetPermissions()

		// Check if it's even meaningful to run on this repository or if it will just error
		// when trying to do other actions
		switch {
		case r.GetArchived():
			log.Debug("Skipping repository since it's archived")
			continue
		case r.GetDisabled():
			log.Debug("Skipping repository since it's disabled")
			continue
		case len(g.Topics) != 0 && !scm.RepoContainsTopic(r.Topics, g.Topics):
			log.Debug("Skipping repository since it does not match repository topics")
			continue
		case g.SkipForks && r.GetFork():
			log.Debug("Skipping repository since it's a fork")
			continue
		}

		if g.checkPermissions {
			switch {
			case !permissions["pull"]:
				log.Debug("Skipping repository since the token does not have pull permissions")
				continue
			case !g.Fork && !g.ReadOnly && !permissions["push"]:
				log.Debug("Skipping repository since the token does not have push permissions and the run will not fork")
				continue
			}
		}

		newRepo, err := g.convertRepo(r)
		if err != nil {
			return nil, err
		}

		repos = append(repos, newRepo)
	}

	return repos, nil
}

func (g *Github) getRepositories(ctx context.Context) ([]*github.Repository, error) {
	allRepos := []*github.Repository{}

	for _, org := range g.Organizations {
		repos, err := g.getOrganizationRepositories(ctx, org)
		if err != nil {
			return nil, errors.Wrapf(err, "could not get organization repositories for %s", org)
		}
		allRepos = append(allRepos, repos...)
	}

	for _, user := range g.Users {
		repos, err := g.getUserRepositories(ctx, user)
		if err != nil {
			return nil, errors.Wrapf(err, "could not get user repositories for %s", user)
		}
		allRepos = append(allRepos, repos...)
	}

	for _, repoRef := range g.Repositories {
		repo, err := g.getRepository(ctx, repoRef)
		if err != nil {
			return nil, errors.Wrapf(err, "could not get information about %s", repoRef.String())
		}
		allRepos = append(allRepos, repo)
	}

	if g.RepositorySearch != "" {
		repos, err := g.getSearchRepositories(ctx, g.RepositorySearch)
		if err != nil {
			return nil, errors.Wrapf(err, "could not get repository search results for '%s'", g.RepositorySearch)
		}
		allRepos = append(allRepos, repos...)
	}

	if g.CodeSearch != "" {
		repos, err := g.getCodeSearchRepositories(ctx, g.CodeSearch)
		if err != nil {
			return nil, errors.Wrapf(err, "could not get code search results for '%s'", g.CodeSearch)
		}
		allRepos = append(allRepos, repos...)
	}

	// Remove duplicate repos
	repoMap := map[string]*github.Repository{}
	for _, repo := range allRepos {
		repoMap[repo.GetFullName()] = repo
	}
	allRepos = make([]*github.Repository, 0, len(repoMap))
	for _, repo := range repoMap {
		if repo.GetArchived() || repo.GetDisabled() {
			continue
		}
		allRepos = append(allRepos, repo)
	}
	sort.Slice(allRepos, func(i, j int) bool {
		return allRepos[i].GetCreatedAt().Before(allRepos[j].GetCreatedAt().Time)
	})

	return allRepos, nil
}

func (g *Github) getOrganizationRepositories(ctx context.Context, orgName string) ([]*github.Repository, error) {
	var repos []*github.Repository
	i := 1
	for {
		rr, _, err := retry(ctx, func() ([]*github.Repository, *github.Response, error) {
			return g.ghClient.Repositories.ListByOrg(ctx, orgName, &github.RepositoryListByOrgOptions{
				ListOptions: github.ListOptions{
					Page:    i,
					PerPage: 100,
				},
			})
		})
		if err != nil {
			return nil, err
		}
		repos = append(repos, rr...)
		if len(rr) != 100 {
			break
		}
		i++
	}

	return repos, nil
}

func (g *Github) getUserRepositories(ctx context.Context, user string) ([]*github.Repository, error) {
	var repos []*github.Repository
	i := 1
	for {
		rr, _, err := retry(ctx, func() ([]*github.Repository, *github.Response, error) {
			return g.ghClient.Repositories.List(ctx, user, &github.RepositoryListOptions{
				ListOptions: github.ListOptions{
					Page:    i,
					PerPage: 100,
				},
			})
		})
		if err != nil {
			return nil, err
		}
		repos = append(repos, rr...)
		if len(rr) != 100 {
			break
		}
		i++
	}

	return repos, nil
}

func (g *Github) getSearchRepositories(ctx context.Context, search string) ([]*github.Repository, error) {
	var repos []*github.Repository
	i := 1
	for {
		rr, _, err := retry(ctx, func() ([]*github.Repository, *github.Response, error) {
			rr, resp, err := g.ghClient.Search.Repositories(ctx, search, &github.SearchOptions{
				ListOptions: github.ListOptions{
					Page:    i,
					PerPage: 100,
				},
			})
			if err != nil {
				return nil, nil, err
			}

			if rr.GetIncompleteResults() {
				// can occur when search times out on the server: for now, fail instead
				// of handling the issue
				return nil, nil, fmt.Errorf("search timed out on GitHub and was marked incomplete: try refining the search to return fewer results or be less complex")
			}

			if rr.GetTotal() > 1000 {
				return nil, nil, fmt.Errorf("%d results for this search, but only the first 1000 results will be returned: try refining your search terms", rr.GetTotal())
			}

			return rr.Repositories, resp, nil
		})
		if err != nil {
			return nil, err
		}
		repos = append(repos, rr...)
		if len(rr) != 100 {
			break
		}
		i++
	}

	return repos, nil
}

func (g *Github) getCodeSearchRepositories(ctx context.Context, search string) ([]*github.Repository, error) {
	resultRepos := make(map[string]RepositoryReference)

	i := 1
	for {
		rr, _, err := retry(ctx, func() ([]*github.CodeResult, *github.Response, error) {
			rr, resp, err := g.ghClient.Search.Code(ctx, search, &github.SearchOptions{
				ListOptions: github.ListOptions{
					Page:    i,
					PerPage: 100,
				},
			})
			if err != nil {
				return nil, nil, err
			}

			if rr.GetIncompleteResults() {
				// can occur when search times out on the server: for now, fail instead
				// of handling the issue
				return nil, nil, fmt.Errorf("search timed out on GitHub and was marked incomplete: try refining the search to return fewer results or be less complex")
			}

			if rr.GetTotal() > 1000 {
				return nil, nil, fmt.Errorf("%d results for this search, but only the first 1000 results will be returned: try refining your search terms", rr.GetTotal())
			}

			return rr.CodeResults, resp, nil
		})
		if err != nil {
			return nil, err
		}

		for _, r := range rr {
			repo := r.Repository

			resultRepos[repo.GetFullName()] = RepositoryReference{
				OwnerName: repo.GetOwner().GetLogin(),
				Name:      repo.GetName(),
			}
		}

		if len(rr) != 100 {
			break
		}
		i++
	}

	// Code search does not return full details (like permissions). So for each
	// repo discovered, we have to query it again.
	repoNames := maps.Values(resultRepos)
	return g.getAllRepositories(ctx, repoNames)
}

func (g *Github) getAllRepositories(ctx context.Context, repoRefs []RepositoryReference) ([]*github.Repository, error) {
	var repos []*github.Repository

	for _, ref := range repoRefs {
		r, err := g.getRepository(ctx, ref)
		if err != nil {
			return nil, err
		}
		repos = append(repos, r)
	}

	return repos, nil
}

func (g *Github) getRepository(ctx context.Context, repoRef RepositoryReference) (*github.Repository, error) {
	repo, _, err := retry(ctx, func() (*github.Repository, *github.Response, error) {
		return g.ghClient.Repositories.Get(ctx, repoRef.OwnerName, repoRef.Name)
	})
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// CreatePullRequest creates a pull request
func (g *Github) CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
	r := repo.(repository)
	prR := prRepo.(repository)

	g.modLock()
	defer g.modUnlock()

	pr, err := g.createPullRequest(ctx, r, prR, newPR)
	if err != nil {
		return nil, err
	}

	if err := g.setReviewers(ctx, r, newPR, pr); err != nil {
		return nil, err
	}

	if err := g.setAssignees(ctx, r, newPR, pr); err != nil {
		return nil, err
	}

	if err := g.setLabels(ctx, r, newPR, pr); err != nil {
		return nil, err
	}

	return convertPullRequest(pr), nil
}

func (g *Github) createPullRequest(ctx context.Context, repo repository, prRepo repository, newPR scm.NewPullRequest) (*github.PullRequest, error) {
	head := fmt.Sprintf("%s:%s", prRepo.ownerName, newPR.Head)

	pr, _, err := retry(ctx, func() (*github.PullRequest, *github.Response, error) {
		return g.ghClient.PullRequests.Create(ctx, repo.ownerName, repo.name, &github.NewPullRequest{
			Title: &newPR.Title,
			Body:  &newPR.Body,
			Head:  &head,
			Base:  &newPR.Base,
			Draft: &newPR.Draft,
		})
	})
	return pr, err
}

func (g *Github) setReviewers(ctx context.Context, repo repository, newPR scm.NewPullRequest, createdPR *github.PullRequest) error {
	var addedReviewers, removedReviewers []string
	if newPR.Reviewers != nil {
		existingReviewers := scm.Map(createdPR.RequestedReviewers, func(user *github.User) string {
			return user.GetLogin()
		})
		addedReviewers, removedReviewers = scm.Diff(existingReviewers, newPR.Reviewers)
	}

	var addedTeamReviewers, removedTeamReviewers []string
	if newPR.TeamReviewers != nil {
		existingTeamReviewers := scm.Map(createdPR.RequestedTeams, func(team *github.Team) string {
			return team.GetSlug()
		})
		addedTeamReviewers, removedTeamReviewers = scm.Diff(existingTeamReviewers, newPR.TeamReviewers)
	}

	if len(addedReviewers) > 0 || len(addedTeamReviewers) > 0 {
		_, _, err := retry(ctx, func() (*github.PullRequest, *github.Response, error) {
			return g.ghClient.PullRequests.RequestReviewers(ctx, repo.ownerName, repo.name, createdPR.GetNumber(), github.ReviewersRequest{
				Reviewers:     addedReviewers,
				TeamReviewers: addedTeamReviewers,
			})
		})
		if err != nil {
			return err
		}
	}

	if len(removedReviewers) > 0 || len(removedTeamReviewers) > 0 {
		_, err := retryWithoutReturn(ctx, func() (*github.Response, error) {
			return g.ghClient.PullRequests.RemoveReviewers(ctx, repo.ownerName, repo.name, createdPR.GetNumber(), github.ReviewersRequest{
				Reviewers:     removedReviewers,
				TeamReviewers: removedTeamReviewers,
			})
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Github) setAssignees(ctx context.Context, repo repository, newPR scm.NewPullRequest, createdPR *github.PullRequest) error {
	if newPR.Assignees == nil {
		return nil
	}

	existingAssignees := scm.Map(createdPR.Assignees, func(user *github.User) string {
		return user.GetLogin()
	})
	addedAssignees, removedAssignees := scm.Diff(existingAssignees, newPR.Assignees)

	if len(addedAssignees) > 0 {
		_, _, err := retry(ctx, func() (*github.Issue, *github.Response, error) {
			return g.ghClient.Issues.AddAssignees(ctx, repo.ownerName, repo.name, createdPR.GetNumber(), addedAssignees)
		})
		if err != nil {
			return err
		}
	}

	if len(removedAssignees) > 0 {
		_, _, err := retry(ctx, func() (*github.Issue, *github.Response, error) {
			return g.ghClient.Issues.RemoveAssignees(ctx, repo.ownerName, repo.name, createdPR.GetNumber(), removedAssignees)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Github) setLabels(ctx context.Context, repo repository, newPR scm.NewPullRequest, createdPR *github.PullRequest) error {
	if newPR.Labels == nil {
		return nil
	}

	existingLabels := scm.Map(createdPR.Labels, func(label *github.Label) string {
		return label.GetName()
	})
	addedLabels, removedLabels := scm.Diff(existingLabels, newPR.Labels)

	if len(addedLabels) > 0 {
		_, _, err := retry(ctx, func() ([]*github.Label, *github.Response, error) {
			return g.ghClient.Issues.AddLabelsToIssue(ctx, repo.ownerName, repo.name, createdPR.GetNumber(), addedLabels)
		})
		if err != nil {
			return err
		}
	}

	for _, label := range removedLabels {
		_, err := retryWithoutReturn(ctx, func() (*github.Response, error) {
			return g.ghClient.Issues.RemoveLabelForIssue(ctx, repo.ownerName, repo.name, createdPR.GetNumber(), label)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdatePullRequest updates an existing pull request
func (g *Github) UpdatePullRequest(ctx context.Context, repo scm.Repository, pullReq scm.PullRequest, updatedPR scm.NewPullRequest) (scm.PullRequest, error) {
	r := repo.(repository)
	pr := pullReq.(pullRequest)

	g.modLock()
	defer g.modUnlock()

	ghPR, _, err := retry(ctx, func() (*github.PullRequest, *github.Response, error) {
		return g.ghClient.PullRequests.Edit(ctx, pr.ownerName, pr.repoName, pr.number, &github.PullRequest{
			Title: &updatedPR.Title,
			Body:  &updatedPR.Body,
		})
	})
	if err != nil {
		return nil, err
	}

	if err := g.setReviewers(ctx, r, updatedPR, ghPR); err != nil {
		return nil, err
	}

	if err := g.setAssignees(ctx, r, updatedPR, ghPR); err != nil {
		return nil, err
	}

	if err := g.setLabels(ctx, r, updatedPR, ghPR); err != nil {
		return nil, err
	}

	return convertPullRequest(ghPR), nil
}

// GetPullRequests gets all pull requests of with a specific branch
func (g *Github) GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error) {
	repos, err := g.getRepositories(ctx)
	if err != nil {
		return nil, err
	}

	// github limits the amount of data which can be handled by the graphql api
	// data needs to be chunked into multiple requests
	batches := chunkSlice(repos, 50)
	var pullRequests []scm.PullRequest

	for _, repos := range batches {
		result, err := g.getPullRequests(ctx, branchName, repos)
		if err != nil {
			return pullRequests, fmt.Errorf("failed to get pull request batch: %w", err)
		}

		pullRequests = append(pullRequests, result...)
	}

	return pullRequests, nil
}

func (g *Github) getPullRequests(ctx context.Context, branchName string, repos []*github.Repository) ([]scm.PullRequest, error) {
	// The fragment is all the data needed from every repository
	const fragment = `fragment repoProperties on Repository {
		pullRequests(headRefName: $branchName, last: 1) {
			nodes {
				number
				headRefName
				closed
				url
				merged
				baseRepository {
					name
					owner {
						login
					}
				}
				headRepository {
					name
					owner {
						login
					}
				}
				commits(last: 1) {
					nodes {
						commit {
							statusCheckRollup {
								state
							}
						}
					}
				}
			}
		}
	}`

	// Prepare data for compiling the query.
	// Each repository will get its own variables ($ownerX, $repoX) and be returned
	// via and alias repoX
	repoParameters := make([]string, len(repos))
	repoQueries := make([]string, len(repos))
	queryVariables := map[string]interface{}{
		"branchName": branchName,
	}
	for i, repo := range repos {
		repoParameters[i] = fmt.Sprintf("$owner%[1]d: String!, $repo%[1]d: String!", i)
		repoQueries[i] = fmt.Sprintf("repo%[1]d: repository(owner: $owner%[1]d, name: $repo%[1]d) { ...repoProperties }", i)

		queryVariables[fmt.Sprintf("owner%d", i)] = repo.GetOwner().GetLogin()
		queryVariables[fmt.Sprintf("repo%d", i)] = repo.GetName()
	}

	// Create the final query
	query := fmt.Sprintf(`
		%s

		query ($branchName: String!, %s) {
			%s
		}`,
		fragment,
		strings.Join(repoParameters, ", "),
		strings.Join(repoQueries, "\n"),
	)

	result := map[string]graphqlRepo{}
	err := g.makeGraphQLRequest(ctx, query, queryVariables, &result)
	if err != nil {
		return nil, err
	}

	// Fetch the repo based on name instead of looping through the map since that will
	// guarantee the same ordering as the original repository list
	prs := []scm.PullRequest{}
	for i := range repos {
		repo, ok := result[fmt.Sprintf("repo%d", i)]
		if !ok {
			return nil, fmt.Errorf("could not find repo%d", i)
		}

		if len(repo.PullRequests.Nodes) != 1 {
			continue
		}
		pr := repo.PullRequests.Nodes[0]

		// The graphql API does not have a way at query time to filter out the owner of the head branch
		// of a PR. Therefore, we have to filter out any repo that does not match the head owner.
		headOwner, err := g.headOwner(ctx, pr.BaseRepository.Owner.Login)
		if err != nil {
			return nil, err
		}
		if pr.HeadRepository.Owner.Login != headOwner {
			continue
		}

		prs = append(prs, convertGraphQLPullRequest(pr))
	}

	return prs, nil
}

func (g *Github) loggedInUser(ctx context.Context) (string, error) {
	g.userMutex.Lock()
	defer g.userMutex.Unlock()

	if g.user != "" {
		return g.user, nil
	}

	user, _, err := retry(ctx, func() (*github.User, *github.Response, error) {
		return g.ghClient.Users.Get(ctx, "")
	})
	if err != nil {
		return "", err
	}

	g.user = user.GetLogin()

	return g.user, nil
}

// headOwner returns the owner of the repository from which any pullrequest will be made
// This is normally the owner of the original repository, but if a fork has been made
// it will be a different owner
func (g *Github) headOwner(ctx context.Context, repoOwner string) (string, error) {
	if !g.Fork {
		return repoOwner, nil
	}

	if g.ForkOwner != "" {
		return g.ForkOwner, nil
	}

	return g.loggedInUser(ctx)
}

// GetOpenPullRequest gets a pull request for one specific repository
func (g *Github) GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
	r := repo.(repository)

	headOwner, err := g.headOwner(ctx, r.ownerName)
	if err != nil {
		return nil, err
	}

	prs, _, err := retry(ctx, func() ([]*github.PullRequest, *github.Response, error) {
		return g.ghClient.PullRequests.List(ctx, headOwner, r.name, &github.PullRequestListOptions{
			Head:  fmt.Sprintf("%s:%s", headOwner, branchName),
			State: "open",
			ListOptions: github.ListOptions{
				PerPage: 1,
			},
		})
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get open pull requests: %w", err)
	}
	if len(prs) == 0 {
		return nil, nil
	}
	return convertPullRequest(prs[0]), nil
}

// MergePullRequest merges a pull request
func (g *Github) MergePullRequest(ctx context.Context, pullReq scm.PullRequest) error {
	pr := pullReq.(pullRequest)

	g.modLock()
	defer g.modUnlock()

	// We need to fetch the repo again since no AllowXMerge is present in listings of repositories
	repo, _, err := retry(ctx, func() (*github.Repository, *github.Response, error) {
		return g.ghClient.Repositories.Get(ctx, pr.ownerName, pr.repoName)
	})
	if err != nil {
		return err
	}

	// Filter out all merge types to only the allowed ones, but keep the order of the ones left
	mergeTypes := scm.MergeTypeIntersection(g.MergeTypes, repoMergeTypes(repo))
	if len(mergeTypes) == 0 {
		return errors.New("none of the configured merge types was permitted")
	}

	_, _, err = retry(ctx, func() (*github.PullRequestMergeResult, *github.Response, error) {
		return g.ghClient.PullRequests.Merge(ctx, pr.ownerName, pr.repoName, pr.number, "", &github.PullRequestOptions{
			MergeMethod: mergeTypeGhName[mergeTypes[0]],
		})
	})
	if err != nil {
		return err
	}

	_, err = retryWithoutReturn(ctx, func() (*github.Response, error) {
		return g.ghClient.Git.DeleteRef(ctx, pr.prOwnerName, pr.prRepoName, fmt.Sprintf("heads/%s", pr.branchName))
	})

	// Ignore errors about the reference not existing since it may be the case that GitHub has already deleted the branch
	if err != nil && !strings.Contains(err.Error(), "Reference does not exist") {
		return err
	}

	return nil
}

// ClosePullRequest closes a pull request
func (g *Github) ClosePullRequest(ctx context.Context, pullReq scm.PullRequest) error {
	pr := pullReq.(pullRequest)

	g.modLock()
	defer g.modUnlock()

	_, _, err := retry(ctx, func() (*github.PullRequest, *github.Response, error) {
		return g.ghClient.PullRequests.Edit(ctx, pr.ownerName, pr.repoName, pr.number, &github.PullRequest{
			State: &[]string{"closed"}[0],
		})
	})
	if err != nil {
		return err
	}

	_, err = retryWithoutReturn(ctx, func() (*github.Response, error) {
		return g.ghClient.Git.DeleteRef(ctx, pr.prOwnerName, pr.prRepoName, fmt.Sprintf("heads/%s", pr.branchName))
	})
	return err
}

// ForkRepository forks a repository. If newOwner is empty, fork on the logged in user
func (g *Github) ForkRepository(ctx context.Context, repo scm.Repository, newOwner string) (scm.Repository, error) {
	r := repo.(repository)

	g.modLock()
	defer g.modUnlock()

	createdRepo, _, err := retry(ctx, func() (*github.Repository, *github.Response, error) {
		return g.ghClient.Repositories.CreateFork(ctx, r.ownerName, r.name, &github.RepositoryCreateForkOptions{
			Organization: newOwner,
		})
	})
	if err != nil {
		if _, isAccepted := err.(*github.AcceptedError); !isAccepted {
			return nil, err
		}

		// Request to fork was accepted, but the repo was not created yet. Poll for the repo to be created
		var err error
		var repo *github.Repository
		for i := 0; i < 10; i++ {
			repo, _, err = retry(ctx, func() (*github.Repository, *github.Response, error) {
				return g.ghClient.Repositories.Get(ctx, createdRepo.GetOwner().GetLogin(), createdRepo.GetName())
			})
			if err != nil {
				time.Sleep(time.Second * 3)
				continue
			}
			// The fork does now exist
			return g.convertRepo(repo)
		}

		return nil, errors.New("time waiting for fork to complete was exceeded")
	}

	return g.convertRepo(createdRepo)
}

// GetAutocompleteOrganizations gets organizations for autocompletion
func (g *Github) GetAutocompleteOrganizations(ctx context.Context, _ string) ([]string, error) {
	orgs, _, err := retry(ctx, func() ([]*github.Organization, *github.Response, error) {
		return g.ghClient.Organizations.List(ctx, "", nil)
	})
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(orgs))
	for i, org := range orgs {
		ret[i] = org.GetLogin()
	}

	return ret, nil
}

// GetAutocompleteUsers gets users for autocompletion
func (g *Github) GetAutocompleteUsers(ctx context.Context, str string) ([]string, error) {
	users, _, err := retry(ctx, func() (*github.UsersSearchResult, *github.Response, error) {
		return g.ghClient.Search.Users(ctx, str, nil)
	})
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(users.Users))
	for i, user := range users.Users {
		ret[i] = user.GetLogin()
	}

	return ret, nil
}

// GetAutocompleteRepositories gets repositories for autocompletion
func (g *Github) GetAutocompleteRepositories(ctx context.Context, str string) ([]string, error) {
	var q string

	// If the user has already provided a org/user, it's much more effective to search based on that
	// comparared to a complete free text search
	splitted := strings.SplitN(str, "/", 2)
	switch {
	case len(splitted) == 2:
		// Search set the user or org (user/org in the search can be used interchangeable)
		q = fmt.Sprintf("user:%s %s in:name", splitted[0], splitted[1])
	default:
		q = fmt.Sprintf("%s in:name", str)
	}

	repos, _, err := retry(ctx, func() (*github.RepositoriesSearchResult, *github.Response, error) {
		return g.ghClient.Search.Repositories(ctx, q, nil)
	})
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(repos.Repositories))
	for i, repositories := range repos.Repositories {
		ret[i] = repositories.GetFullName()
	}

	return ret, nil
}

// modLock is a lock that should be used whenever a modifying request is made against the GitHub API.
// It works as a normal lock, but also makes sure that there is a buffer period of 1 second since the
// last critical section was left. This ensures that we always wait at least one seconds between modifying requests
// and does not hit GitHubs secondary rate limit:
// https://docs.github.com/en/rest/guides/best-practices-for-integrators#dealing-with-secondary-rate-limits
func (g *Github) modLock() {
	g.modMutex.Lock()
	shouldWait := time.Second - time.Since(g.lastModRequest)
	if shouldWait > 0 {
		log.Debugf("Waiting %s to not hit GitHub ratelimit", shouldWait)
		time.Sleep(shouldWait)
	}
}

func (g *Github) modUnlock() {
	g.lastModRequest = time.Now()
	g.modMutex.Unlock()
}

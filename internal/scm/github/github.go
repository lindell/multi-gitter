package github

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/lindell/multi-gitter/internal/git"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// New create a new Github client
func New(
	token string,
	baseURL string,
	transportMiddleware func(http.RoundTripper) http.RoundTripper,
	repoListing RepositoryListing,
	mergeTypes []git.MergeType,
	forkMode bool,
	forkOwner string,
) (*Github, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	tc.Transport = transportMiddleware(tc.Transport)

	var client *github.Client
	if baseURL != "" {
		var err error
		client, err = github.NewEnterpriseClient(baseURL, "", tc)
		if err != nil {
			return nil, err
		}
	} else {
		client = github.NewClient(tc)
	}

	return &Github{
		RepositoryListing: repoListing,
		MergeTypes:        mergeTypes,
		token:             token,
		Fork:              forkMode,
		ForkOwner:         forkOwner,
		ghClient:          client,
	}, nil
}

// Github contain github configuration
type Github struct {
	RepositoryListing
	MergeTypes []git.MergeType
	token      string

	// This determines if forks will be used when creating a prs.
	// In this package, it mainly determines which repos are possible to make changes on
	Fork bool

	// If set, the fork will happen to the ForkOwner value, and not the logged in user
	ForkOwner string

	ghClient *github.Client

	// Caching of the logged in user
	user string
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
func (g Github) GetRepositories(ctx context.Context) ([]git.Repository, error) {
	allRepos, err := g.getRepositories(ctx)
	if err != nil {
		return nil, err
	}

	repos := make([]git.Repository, 0, len(allRepos))
	for _, r := range allRepos {
		permissions := r.GetPermissions()

		if r.GetArchived() || r.GetDisabled() || !permissions["pull"] {
			continue
		}
		// The user needs push permissions or have defined that the pr should be on a fork
		if !g.Fork && !permissions["push"] {
			continue
		}

		newRepo, err := convertRepo(r, g.token)
		if err != nil {
			return nil, err
		}

		repos = append(repos, newRepo)
	}

	return repos, nil
}

func (g Github) getRepositories(ctx context.Context) ([]*github.Repository, error) {
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

func (g Github) getOrganizationRepositories(ctx context.Context, orgName string) ([]*github.Repository, error) {
	var repos []*github.Repository
	i := 1
	for {
		rr, _, err := g.ghClient.Repositories.ListByOrg(ctx, orgName, &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{
				Page:    i,
				PerPage: 100,
			},
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

func (g Github) getUserRepositories(ctx context.Context, user string) ([]*github.Repository, error) {
	var repos []*github.Repository
	i := 1
	for {
		rr, _, err := g.ghClient.Repositories.List(ctx, user, &github.RepositoryListOptions{
			ListOptions: github.ListOptions{
				Page:    i,
				PerPage: 100,
			},
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

func (g Github) getRepository(ctx context.Context, repoRef RepositoryReference) (*github.Repository, error) {
	repo, _, err := g.ghClient.Repositories.Get(ctx, repoRef.OwnerName, repoRef.Name)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// CreatePullRequest creates a pull request
func (g Github) CreatePullRequest(ctx context.Context, repo git.Repository, prRepo git.Repository, newPR git.NewPullRequest) (git.PullRequest, error) {
	r := repo.(repository)
	prR := prRepo.(repository)

	pr, err := g.createPullRequest(ctx, r, prR, newPR)
	if err != nil {
		return nil, err
	}

	if err := g.addReviewers(ctx, r, newPR, pr); err != nil {
		return nil, err
	}

	return convertPullRequest(pr), nil
}

func (g Github) createPullRequest(ctx context.Context, repo repository, prRepo repository, newPR git.NewPullRequest) (*github.PullRequest, error) {
	head := fmt.Sprintf("%s:%s", prRepo.ownerName, newPR.Head)
	pr, _, err := g.ghClient.PullRequests.Create(ctx, repo.ownerName, repo.name, &github.NewPullRequest{
		Title: &newPR.Title,
		Body:  &newPR.Body,
		Head:  &head,
		Base:  &newPR.Base,
	})
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (g Github) addReviewers(ctx context.Context, repo repository, newPR git.NewPullRequest, createdPR *github.PullRequest) error {
	if len(newPR.Reviewers) == 0 {
		return nil
	}
	_, _, err := g.ghClient.PullRequests.RequestReviewers(ctx, repo.ownerName, repo.name, createdPR.GetNumber(), github.ReviewersRequest{
		Reviewers: newPR.Reviewers,
	})
	return err
}

// GetPullRequests gets all pull requests of with a specific branch
func (g Github) GetPullRequests(ctx context.Context, branchName string) ([]git.PullRequest, error) {
	// TODO: If this is implemented with the GitHub v4 graphql api, it would be much faster

	repos, err := g.getRepositories(ctx)
	if err != nil {
		return nil, err
	}

	prStatuses := []git.PullRequest{}
	for _, r := range repos {
		repoOwner := r.GetOwner().GetLogin()
		repoName := r.GetName()
		log := log.WithField("repo", fmt.Sprintf("%s/%s", repoOwner, repoName))

		headOwner := repoOwner
		if g.Fork {
			if g.ForkOwner != "" {
				headOwner = g.ForkOwner
			} else {
				headOwner, err = g.loggedInUser(ctx)
				if err != nil {
					return nil, err
				}
			}
		}

		log.Debug("Fetching latest pull request")
		prs, _, err := g.ghClient.PullRequests.List(ctx, repoOwner, repoName, &github.PullRequestListOptions{
			Head:      fmt.Sprintf("%s:%s", headOwner, branchName),
			State:     "all",
			Direction: "desc",
			ListOptions: github.ListOptions{
				PerPage: 1,
			},
		})
		if err != nil {
			return nil, err
		}
		if len(prs) != 1 {
			continue
		}
		pr := prs[0]

		status, err := g.getPrStatus(ctx, pr)
		if err != nil {
			return nil, err
		}

		localPR := convertPullRequest(pr)
		localPR.status = status
		prStatuses = append(prStatuses, localPR)
	}

	return prStatuses, nil
}

func (g Github) loggedInUser(ctx context.Context) (string, error) {
	if g.user != "" {
		return g.user, nil
	}

	user, _, err := g.ghClient.Users.Get(ctx, "")
	if err != nil {
		return "", err
	}

	g.user = user.GetLogin()

	return g.user, nil
}

// MergePullRequest merges a pull request
func (g Github) MergePullRequest(ctx context.Context, pullReq git.PullRequest) error {
	pr := pullReq.(pullRequest)

	// We need to fetch the repo again since no AllowXMerge is present in listings of repositories
	repo, _, err := g.ghClient.Repositories.Get(ctx, pr.ownerName, pr.repoName)
	if err != nil {
		return err
	}

	// Filter out all merge types to only the allowed ones, but keep the order of the ones left
	mergeTypes := git.MergeTypeIntersection(g.MergeTypes, repoMergeTypes(repo))
	if len(mergeTypes) == 0 {
		return errors.New("none of the configured merge types was permitted")
	}

	_, _, err = g.ghClient.PullRequests.Merge(ctx, pr.ownerName, pr.repoName, pr.number, "", &github.PullRequestOptions{
		MergeMethod: mergeTypeGhName[mergeTypes[0]],
	})
	if err != nil {
		return err
	}

	_, err = g.ghClient.Git.DeleteRef(ctx, pr.prOwnerName, pr.prRepoName, fmt.Sprintf("heads/%s", pr.branchName))

	// Ignore errors about the reference not existing since it may be the case that GitHub has already deleted the branch
	if err != nil && !strings.Contains(err.Error(), "Reference does not exist") {
		return err
	}

	return nil
}

// ClosePullRequest closes a pull request
func (g Github) ClosePullRequest(ctx context.Context, pullReq git.PullRequest) error {
	pr := pullReq.(pullRequest)

	_, _, err := g.ghClient.PullRequests.Edit(ctx, pr.ownerName, pr.repoName, pr.number, &github.PullRequest{
		State: &[]string{"closed"}[0],
	})
	if err != nil {
		return err
	}

	_, err = g.ghClient.Git.DeleteRef(ctx, pr.prOwnerName, pr.prRepoName, fmt.Sprintf("heads/%s", pr.branchName))
	return err
}

// ForkRepository forks a repository. If newOwner is empty, fork on the logged in user
func (g Github) ForkRepository(ctx context.Context, repo git.Repository, newOwner string) (git.Repository, error) {
	r := repo.(repository)

	createdRepo, _, err := g.ghClient.Repositories.CreateFork(ctx, r.ownerName, r.name, &github.RepositoryCreateForkOptions{
		Organization: newOwner,
	})
	if err != nil {
		if _, isAccepted := err.(*github.AcceptedError); !isAccepted {
			return nil, err
		}

		// Request to fork was accepted, but the repo was not created yet. Poll for the repo to be created
		var err error
		var repo *github.Repository
		for i := 0; i < 10; i++ {
			repo, _, err = g.ghClient.Repositories.Get(ctx, createdRepo.GetOwner().GetLogin(), createdRepo.GetName())
			if err != nil {
				time.Sleep(time.Second * 3)
				continue
			}
			// The fork does now exist
			return convertRepo(repo, g.token)
		}

		return nil, errors.New("time waiting for fork to complete was exceeded")
	}

	return convertRepo(createdRepo, g.token)
}

// GetAutocompleteOrganizations gets organizations for autocompletion
func (g Github) GetAutocompleteOrganizations(ctx context.Context, _ string) ([]string, error) {
	orgs, _, err := g.ghClient.Organizations.List(ctx, "", nil)
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
func (g Github) GetAutocompleteUsers(ctx context.Context, str string) ([]string, error) {
	users, _, err := g.ghClient.Search.Users(ctx, str, nil)
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
func (g Github) GetAutocompleteRepositories(ctx context.Context, str string) ([]string, error) {
	var q string

	// If the user has already provided a org/user, it's much more effective to search based on that
	// comparared to a complete freetext search
	splitted := strings.SplitN(str, "/", 2)
	switch {
	case len(splitted) == 2:
		// Search set the user or org (user/org in the search can be used interchangeable)
		q = fmt.Sprintf("user:%s %s in:name", splitted[0], splitted[1])
	default:
		q = fmt.Sprintf("%s in:name", str)
	}

	repos, _, err := g.ghClient.Search.Repositories(ctx, q, nil)
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(repos.Repositories))
	for i, repositories := range repos.Repositories {
		ret[i] = repositories.GetFullName()
	}

	return ret, nil
}

func (g Github) getPrStatus(ctx context.Context, pr *github.PullRequest) (git.PullRequestStatus, error) {
	// Determine the status of the pr
	var status git.PullRequestStatus
	if pr.MergedAt != nil {
		status = git.PullRequestStatusMerged
	} else if pr.ClosedAt != nil {
		status = git.PullRequestStatusClosed
	} else {
		log.Debug("Fetching the combined status of the pull request")
		combinedStatus, _, err := g.ghClient.Repositories.GetCombinedStatus(ctx, pr.GetBase().GetUser().GetLogin(), pr.GetBase().GetRepo().GetName(), pr.GetHead().GetSHA(), nil)
		if err != nil {
			return git.PullRequestStatusUnknown, err
		}

		if combinedStatus.GetTotalCount() == 0 {
			status = git.PullRequestStatusSuccess
		} else {
			switch combinedStatus.GetState() {
			case "pending":
				status = git.PullRequestStatusPending
			case "success":
				status = git.PullRequestStatusSuccess
			case "failure", "error":
				status = git.PullRequestStatusError
			}
		}
	}

	return status, nil
}

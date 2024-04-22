package gitea

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	internalHTTP "github.com/lindell/multi-gitter/internal/http"
)

// New create a new Gitea client
func New(token, baseURL string, repoListing RepositoryListing, mergeTypes []scm.MergeType, sshAuth bool) (*Gitea, error) {
	gitea := &Gitea{
		RepositoryListing: repoListing,

		baseURL: baseURL,
		token:   token,

		MergeTypes: mergeTypes,
		SSHAuth:    sshAuth,
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

	MergeTypes []scm.MergeType
	SSHAuth    bool
}

// RepositoryListing contains information about which repositories that should be fetched
type RepositoryListing struct {
	Organizations []string
	Users         []string
	Repositories  []RepositoryReference
	Topics        []string
	SkipForks     bool
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
func (g *Gitea) GetRepositories(ctx context.Context) ([]scm.Repository, error) {
	allRepos, err := g.getRepositories(ctx)
	if err != nil {
		return nil, err
	}

	repos := make([]scm.Repository, 0, len(allRepos))
	for _, repo := range allRepos {
		log := log.WithField("repo", repo.FullName)

		if g.SkipForks && repo.Fork {
			log.Debug("Skipping repository since it's a fork")
			continue
		}

		if len(g.Topics) != 0 {
			topics, err := g.getRepoTopics(ctx, repo)
			if err != nil {
				return repos, fmt.Errorf("could not fetch repository topics: %w", err)
			}

			if !scm.RepoContainsTopic(topics, g.Topics) {
				log.Debug("Skipping repository since it does not match repository topics")
				continue
			}
		}

		convertedRepo, err := g.convertRepository(repo)
		if err != nil {
			return nil, err
		}
		repos = append(repos, convertedRepo)
	}

	return repos, nil
}

func (g *Gitea) getRepoTopics(ctx context.Context, repo *gitea.Repository) ([]string, error) {
	topics, _, err := g.giteaClient(ctx).ListRepoTopics(repo.Owner.UserName, repo.Name, gitea.ListRepoTopicsOptions{})
	return topics, err
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
func (g *Gitea) CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
	r := repo.(repository)
	prR := prRepo.(repository)

	head := newPR.Head
	if r.ownerName != prR.ownerName {
		head = fmt.Sprintf("%s:%s", prR.ownerName, newPR.Head)
	}

	prTitle := newPR.Title
	if newPR.Draft {
		prTitle = "WIP: " + prTitle // See https://docs.gitea.io/en-us/pull-request/
	}

	labels, err := g.getLabelsFromStrings(ctx, r, newPR.Labels)
	if err != nil {
		return nil, errors.WithMessage(err, "could not map labels")
	}

	pr, _, err := g.giteaClient(ctx).CreatePullRequest(r.ownerName, r.name, gitea.CreatePullRequestOption{
		Head:      head,
		Base:      newPR.Base,
		Title:     prTitle,
		Body:      newPR.Body,
		Assignees: newPR.Assignees,
		Labels:    labels,
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

	return pullRequest{
		repoName:    r.name,
		ownerName:   r.ownerName,
		branchName:  newPR.Head,
		prOwnerName: pr.Head.Repository.Owner.UserName,
		prRepoName:  pr.Head.Repository.Name,
		index:       pr.Index,
		webURL:      pr.HTMLURL,
	}, nil
}

func (g *Gitea) getLabelsFromStrings(ctx context.Context, repo repository, labelNames []string) ([]int64, error) {
	if len(labelNames) == 0 {
		return nil, nil
	}

	labels, _, err := g.giteaClient(ctx).ListRepoLabels(repo.ownerName, repo.name, gitea.ListLabelsOptions{})
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	labelMap := map[string]int64{}
	for _, label := range labels {
		labelMap[label.Name] = label.ID
	}

	// Get the ids of all labels, if the label-name does not exist, simply skip it
	var ret []int64
	for _, name := range labelNames {
		if id, ok := labelMap[name]; ok {
			ret = append(ret, id)
		}
	}

	return ret, nil
}

func (g *Gitea) setReviewers(ctx context.Context, repo repository, newPR scm.NewPullRequest, createdPR *gitea.PullRequest) error {
	if newPR.Reviewers == nil {
		return nil
	}

	reviews, _, err := g.giteaClient(ctx).ListPullReviews(repo.ownerName, repo.name, createdPR.Index, gitea.ListPullReviewsOptions{})
	if err != nil {
		return errors.Wrap(err, "could not list existing reviews on pull request")
	}
	reviews = slices.DeleteFunc(reviews, func(review *gitea.PullReview) bool {
		return review.State != gitea.ReviewStateRequestReview // can only remove reviews in requested state
	})
	existingReviewers := scm.Map(reviews, func(review *gitea.PullReview) string {
		return review.Reviewer.UserName
	})
	addedReviewers, removedReviewers := scm.Diff(existingReviewers, newPR.Reviewers)

	if len(addedReviewers) > 0 {
		_, err := g.giteaClient(ctx).CreateReviewRequests(repo.ownerName, repo.name, createdPR.Index, gitea.PullReviewRequestOptions{
			Reviewers: addedReviewers,
		})
		if err != nil {
			return errors.Wrap(err, "could not add reviewers to pull request")
		}
	}

	if len(removedReviewers) > 0 {
		_, err := g.giteaClient(ctx).DeleteReviewRequests(repo.ownerName, repo.name, createdPR.Index, gitea.PullReviewRequestOptions{
			Reviewers: removedReviewers,
		})
		if err != nil {
			return errors.Wrap(err, "could not remove reviewers from pull request")
		}
	}

	return nil
}

// UpdatePullRequest updates an existing pull request
func (g *Gitea) UpdatePullRequest(ctx context.Context, repo scm.Repository, pullReq scm.PullRequest, updatedPR scm.NewPullRequest) (scm.PullRequest, error) {
	r := repo.(repository)
	pr := pullReq.(pullRequest)

	prTitle := updatedPR.Title
	if updatedPR.Draft {
		prTitle = "WIP: " + prTitle // See https://docs.gitea.io/en-us/pull-request/
	}

	labels, err := g.getLabelsFromStrings(ctx, r, updatedPR.Labels)
	if err != nil {
		return nil, errors.WithMessage(err, "could not map labels")
	}

	giteaPr, _, err := g.giteaClient(ctx).EditPullRequest(r.ownerName, r.name, pr.index, gitea.EditPullRequestOption{
		Title:     prTitle,
		Body:      updatedPR.Body,
		Assignees: updatedPR.Assignees,
		Labels:    labels,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not update pull request")
	}

	if err := g.setReviewers(ctx, r, updatedPR, giteaPr); err != nil {
		return nil, err
	}

	return pullRequest{
		repoName:    r.name,
		ownerName:   r.ownerName,
		branchName:  giteaPr.Head.Name,
		prOwnerName: giteaPr.Head.Repository.Owner.UserName,
		prRepoName:  giteaPr.Head.Repository.Name,
		index:       giteaPr.Index,
		webURL:      giteaPr.HTMLURL,
	}, nil
}

// GetPullRequests gets all pull requests of with a specific branch
func (g *Gitea) GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error) {
	repos, err := g.getRepositories(ctx)
	if err != nil {
		return nil, err
	}

	prs := []scm.PullRequest{}
	for _, repo := range repos {
		pr, err := g.getPullRequest(ctx, branchName, repo.Owner.UserName, repo.Name, gitea.StateAll)
		if err != nil {
			return nil, err
		}
		if pr == nil {
			continue
		}

		convertedPR, err := g.convertPullRequest(ctx, pr)
		if err != nil {
			return nil, err
		}

		prs = append(prs, convertedPR)
	}

	return prs, nil
}

func (g *Gitea) convertPullRequest(ctx context.Context, pr *gitea.PullRequest) (pullRequest, error) {
	status, err := g.pullRequestStatus(ctx, pr)
	if err != nil {
		return pullRequest{}, err
	}

	return pullRequest{
		repoName:    pr.Base.Repository.Name,
		ownerName:   pr.Base.Repository.Owner.UserName,
		branchName:  pr.Head.Name,
		prOwnerName: pr.Head.Repository.Owner.UserName,
		prRepoName:  pr.Head.Repository.Name,
		status:      status,
		index:       pr.Index,
		webURL:      pr.HTMLURL,
	}, nil
}

func (g *Gitea) getPullRequest(ctx context.Context, branchName string, owner, repoName string, state gitea.StateType) (*gitea.PullRequest, error) {
	// We would like to be able to search for a pr with a specific head here, but current (2021-04-24), that option does not exist in the API
	prs, _, err := g.giteaClient(ctx).ListRepoPullRequests(owner, repoName, gitea.ListPullRequestsOptions{
		State: state,
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

func (g *Gitea) pullRequestStatus(ctx context.Context, pr *gitea.PullRequest) (scm.PullRequestStatus, error) {
	if pr.Merged != nil {
		return scm.PullRequestStatusMerged, nil
	}

	if pr.State == gitea.StateClosed {
		return scm.PullRequestStatusClosed, nil
	}

	status, _, err := g.giteaClient(ctx).GetCombinedStatus(pr.Base.Repository.Owner.UserName, pr.Base.Repository.Name, pr.Head.Sha)
	if err != nil {
		return scm.PullRequestStatusUnknown, err
	}

	if len(status.Statuses) == 0 {
		return scm.PullRequestStatusSuccess, nil
	}

	switch status.State {
	case gitea.StatusPending:
		return scm.PullRequestStatusPending, nil
	case gitea.StatusSuccess:
		return scm.PullRequestStatusSuccess, nil
	case gitea.StatusError, gitea.StatusFailure:
		return scm.PullRequestStatusError, nil
	}

	return scm.PullRequestStatusUnknown, nil
}

// GetOpenPullRequest gets a pull request for one specific repository
func (g *Gitea) GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
	r := repo.(repository)

	pr, err := g.getPullRequest(ctx, branchName, r.ownerName, r.name, gitea.StateOpen)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, nil
	}

	return g.convertPullRequest(ctx, pr)
}

// MergePullRequest merges a pull request
func (g *Gitea) MergePullRequest(ctx context.Context, pullReq scm.PullRequest) error {
	pr := pullReq.(pullRequest)

	repo, _, err := g.giteaClient(ctx).GetRepo(pr.ownerName, pr.repoName)
	if err != nil {
		return errors.Wrapf(err, "could not fetch %s/%s repository", pr.ownerName, pr.repoName)
	}

	// Filter out all merge types to only the allowed ones, but keep the order of the ones left
	mergeTypes := scm.MergeTypeIntersection(g.MergeTypes, repoMergeTypes(repo))
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
func (g *Gitea) ClosePullRequest(ctx context.Context, pullReq scm.PullRequest) error {
	pr := pullReq.(pullRequest)

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
func (g *Gitea) ForkRepository(ctx context.Context, repo scm.Repository, newOwner string) (scm.Repository, error) {
	r := repo.(repository)

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
		return g.convertRepository(existingRepo)
	}

	forkOptions := gitea.CreateForkOption{}
	if newOwner != "" {
		forkOptions.Organization = &newOwner
	}

	createdRepo, _, err := g.giteaClient(ctx).CreateFork(r.ownerName, r.name, forkOptions)
	if err != nil {
		return nil, err
	}

	return g.convertRepository(createdRepo)
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

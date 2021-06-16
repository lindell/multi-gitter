package gitlab

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"

	"github.com/lindell/multi-gitter/internal/domain"
	internalHTTP "github.com/lindell/multi-gitter/internal/http"
)

// New create a new Gitlab client
func New(token, baseURL string, repoListing RepositoryListing, config Config) (*Gitlab, error) {
	var options []gitlab.ClientOptionFunc
	if baseURL != "" {
		options = append(options, gitlab.WithBaseURL(baseURL))
	}

	options = append(options, gitlab.WithHTTPClient(&http.Client{
		Transport: internalHTTP.LoggingRoundTripper{},
	}))

	client, err := gitlab.NewClient(token, options...)
	if err != nil {
		return nil, err
	}

	return &Gitlab{
		RepositoryListing: repoListing,
		Config:            config,
		glClient:          client,
	}, nil
}

// Gitlab contain gitlab configuration
type Gitlab struct {
	RepositoryListing
	Config   Config
	glClient *gitlab.Client

	// Cached current user
	currentUser *gitlab.User
}

// RepositoryListing contains information about which repositories that should be fetched
type RepositoryListing struct {
	Groups   []string
	Users    []string
	Projects []ProjectReference
}

// Config includes extra config parameters for the GitLab client
type Config struct {
	IncludeSubgroups bool
}

// ProjectReference contains information to be able to reference a repository
type ProjectReference struct {
	OwnerName string
	Name      string
}

// ParseProjectReference parses a repository reference from the format "ownerName/repoName"
func ParseProjectReference(val string) (ProjectReference, error) {
	split := strings.Split(val, "/")
	if len(split) != 2 {
		return ProjectReference{}, fmt.Errorf("could not parse repository reference: %s", val)
	}
	return ProjectReference{
		OwnerName: split[0],
		Name:      split[1],
	}, nil
}

type repository struct {
	url           url.URL
	pid           int
	name          string
	ownerName     string
	defaultBranch string
}

func (r repository) URL(token string) string {
	// Set the token as https://oauth2:TOKEN@url
	r.url.User = url.UserPassword("oauth2", token)
	return r.url.String()
}

func (r repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r repository) FullName() string {
	return fmt.Sprintf("%s/%s", r.ownerName, r.name)
}

type pullRequest struct {
	ownerName  string
	repoName   string
	targetPID  int
	sourcePID  int
	branchName string
	iid        int
	webURL     string
	status     domain.PullRequestStatus
}

func (pr pullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.ownerName, pr.repoName, pr.iid)
}

func (pr pullRequest) Status() domain.PullRequestStatus {
	return pr.status
}

func (pr pullRequest) URL() string {
	return pr.webURL
}

// GetRepositories fetches repositories from all sources (groups/user/specific project)
func (g *Gitlab) GetRepositories(ctx context.Context) ([]domain.Repository, error) {
	allProjects, err := g.getProjects(ctx)
	if err != nil {
		return nil, err
	}

	repos := make([]domain.Repository, 0, len(allProjects))
	for _, project := range allProjects {
		p, err := convertProject(project)
		if err != nil {
			return nil, err
		}

		repos = append(repos, p)
	}

	return repos, nil
}

func (g *Gitlab) getProjects(ctx context.Context) ([]*gitlab.Project, error) {
	allProjects := []*gitlab.Project{}

	for _, group := range g.Groups {
		projects, err := g.getGroupProjects(ctx, group)
		if err != nil {
			return nil, err
		}
		allProjects = append(allProjects, projects...)
	}

	for _, user := range g.Users {
		projects, err := g.getUserProjects(ctx, user)
		if err != nil {
			return nil, err
		}
		allProjects = append(allProjects, projects...)
	}

	for _, project := range g.Projects {
		project, err := g.getProject(ctx, project)
		if err != nil {
			return nil, err
		}
		allProjects = append(allProjects, project)
	}

	// Remove duplicate projects
	projectMap := map[int]*gitlab.Project{}
	for _, proj := range allProjects {
		projectMap[proj.ID] = proj
	}
	allProjects = make([]*gitlab.Project, 0, len(projectMap))
	for _, proj := range projectMap {
		allProjects = append(allProjects, proj)
	}
	sort.Slice(allProjects, func(i, j int) bool {
		return allProjects[i].ID < allProjects[j].ID
	})

	return allProjects, nil
}

func (g *Gitlab) getGroupProjects(ctx context.Context, groupName string) ([]*gitlab.Project, error) {
	var allProjects []*gitlab.Project
	for i := 1; ; i++ {
		projects, _, err := g.glClient.Groups.ListGroupProjects(groupName, &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    i,
			},
			IncludeSubgroups: &g.Config.IncludeSubgroups,
		}, gitlab.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		allProjects = append(allProjects, projects...)

		if len(projects) < 100 {
			break
		}
	}
	return allProjects, nil
}

func (g *Gitlab) getProject(ctx context.Context, projRef ProjectReference) (*gitlab.Project, error) {
	project, _, err := g.glClient.Projects.GetProject(
		fmt.Sprintf("%s/%s", projRef.OwnerName, projRef.Name),
		nil,
		gitlab.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	return project, err
}

func (g *Gitlab) getUserProjects(ctx context.Context, username string) ([]*gitlab.Project, error) {
	var allProjects []*gitlab.Project
	for i := 1; ; i++ {
		projects, _, err := g.glClient.Projects.ListUserProjects(username, &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    i,
			},
		}, gitlab.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		allProjects = append(allProjects, projects...)

		if len(projects) < 100 {
			break
		}
	}
	return allProjects, nil
}

// CreatePullRequest creates a pull request
func (g *Gitlab) CreatePullRequest(ctx context.Context, repo domain.Repository, prRepo domain.Repository, newPR domain.NewPullRequest) (domain.PullRequest, error) {
	r := repo.(repository)
	prR := prRepo.(repository)

	// Convert from usernames to user ids
	var assigneeIDs []int
	if len(newPR.Reviewers) > 0 {
		var err error
		assigneeIDs, err = g.getUserIDs(ctx, newPR.Reviewers)
		if err != nil {
			return nil, err
		}
	}

	removeSourceBranch := true
	mr, _, err := g.glClient.MergeRequests.CreateMergeRequest(prR.pid, &gitlab.CreateMergeRequestOptions{
		Title:              &newPR.Title,
		Description:        &newPR.Body,
		SourceBranch:       &newPR.Head,
		TargetBranch:       &newPR.Base,
		TargetProjectID:    &r.pid,
		AssigneeIDs:        assigneeIDs,
		RemoveSourceBranch: &removeSourceBranch,
	})
	if err != nil {
		return nil, err
	}

	return pullRequest{
		repoName:   r.name,
		ownerName:  r.ownerName,
		targetPID:  mr.TargetProjectID,
		sourcePID:  mr.SourceProjectID,
		branchName: newPR.Head,
		iid:        mr.IID,
		webURL:     mr.WebURL,
	}, nil
}

func (g *Gitlab) getUserIDs(ctx context.Context, usernames []string) ([]int, error) {
	userIDs := make([]int, len(usernames))
	for i := range usernames {
		users, _, err := g.glClient.Users.ListUsers(&gitlab.ListUsersOptions{
			Username: &usernames[i],
		}, gitlab.WithContext(ctx))
		if err != nil {
			return nil, err
		}
		if len(users) != 1 {
			return nil, fmt.Errorf("could not find user: %s", usernames[i])
		}
		userIDs[i] = users[0].ID
	}
	return userIDs, nil
}

// GetPullRequests gets all pull requests of with a specific branch
func (g *Gitlab) GetPullRequests(ctx context.Context, branchName string) ([]domain.PullRequest, error) {
	projects, err := g.getProjects(ctx)
	if err != nil {
		return nil, err
	}

	prs := []domain.PullRequest{}
	for _, project := range projects {
		mr, err := g.getPullRequest(ctx, branchName, project)
		if err != nil {
			return nil, err
		}
		if mr == nil {
			continue
		}

		prs = append(prs, pullRequest{
			repoName:   project.Path,
			ownerName:  project.Namespace.Path,
			targetPID:  mr.TargetProjectID,
			sourcePID:  mr.SourceProjectID,
			branchName: branchName,
			status:     pullRequestStatus(mr),
			iid:        mr.IID,
			webURL:     mr.WebURL,
		})
	}

	return prs, nil
}

func (g *Gitlab) getPullRequest(ctx context.Context, branchName string, project *gitlab.Project) (*gitlab.MergeRequest, error) {
	mrs, _, err := g.glClient.MergeRequests.ListProjectMergeRequests(project.ID, &gitlab.ListProjectMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 1,
		},
		SourceBranch: &branchName,
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if len(mrs) == 0 {
		return nil, nil
	}

	mr, _, err := g.glClient.MergeRequests.GetMergeRequest(project.ID, mrs[0].IID, nil, gitlab.WithContext(ctx))
	if err != nil {
		return mrs[0], err
	}
	return mr, nil
}

func pullRequestStatus(mr *gitlab.MergeRequest) domain.PullRequestStatus {
	switch {
	case mr.MergedAt != nil:
		return domain.PullRequestStatusMerged
	case mr.ClosedAt != nil:
		return domain.PullRequestStatusClosed
	case mr.Pipeline == nil, mr.Pipeline.Status == "success":
		return domain.PullRequestStatusSuccess
	case mr.Pipeline.Status == "failed":
		return domain.PullRequestStatusError
	default:
		return domain.PullRequestStatusPending
	}
}

// MergePullRequest merges a pull request
func (g *Gitlab) MergePullRequest(ctx context.Context, pullReq domain.PullRequest) error {
	pr := pullReq.(pullRequest)

	shouldRemoveSourceBranch := true
	_, _, err := g.glClient.MergeRequests.AcceptMergeRequest(pr.targetPID, pr.iid, &gitlab.AcceptMergeRequestOptions{
		ShouldRemoveSourceBranch: &shouldRemoveSourceBranch,
	}, gitlab.WithContext(ctx))
	if err != nil {
		return err
	}

	return nil
}

// ClosePullRequest closes a pull request
func (g *Gitlab) ClosePullRequest(ctx context.Context, pullReq domain.PullRequest) error {
	pr := pullReq.(pullRequest)

	_, err := g.glClient.MergeRequests.DeleteMergeRequest(pr.targetPID, pr.iid, gitlab.WithContext(ctx))
	if err != nil {
		return err
	}

	_, err = g.glClient.Branches.DeleteBranch(pr.sourcePID, pr.branchName, gitlab.WithContext(ctx))
	if err != nil {
		return err
	}

	return nil
}

// ForkRepository forks a project
func (g *Gitlab) ForkRepository(ctx context.Context, repo domain.Repository, newOwner string) (domain.Repository, error) {
	r := repo.(repository)

	// Get the username of the fork (logged in user if none is set)
	ownerUsername := newOwner
	if newOwner == "" {
		currentUser, err := g.getCurrentUser(ctx)
		if err != nil {
			return nil, err
		}
		ownerUsername = currentUser.Username
	}

	// Check if the project already exist
	project, resp, err := g.glClient.Projects.GetProject(
		fmt.Sprintf("%s/%s", ownerUsername, r.name),
		nil,
		gitlab.WithContext(ctx),
	)
	if err == nil { // Already forked, just return it
		return convertProject(project)
	} else if resp.StatusCode != http.StatusNotFound { // If the error was that the project does not exist, continue to fork it
		return nil, err
	}

	newRepo, _, err := g.glClient.Projects.ForkProject(r.pid, &gitlab.ForkProjectOptions{
		Namespace: &newOwner,
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	for i := 0; i < 10; i++ {
		repo, _, err := g.glClient.Projects.GetProject(newRepo.ID, nil, gitlab.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		if repo.ImportStatus == "finished" {
			return convertProject(newRepo)
		}

		time.Sleep(time.Second * 3)
	}

	return nil, errors.New("time waiting for fork to complete was exceeded")
}

func (g *Gitlab) getCurrentUser(ctx context.Context) (*gitlab.User, error) {
	if g.currentUser != nil {
		return g.currentUser, nil
	}

	user, _, err := g.glClient.Users.CurrentUser(gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	g.currentUser = user

	return user, nil
}

func convertProject(project *gitlab.Project) (repository, error) {
	u, err := url.Parse(project.HTTPURLToRepo)
	if err != nil {
		return repository{}, err
	}

	return repository{
		url:           *u,
		pid:           project.ID,
		name:          project.Path,
		ownerName:     project.Namespace.Path,
		defaultBranch: project.DefaultBranch,
	}, nil
}

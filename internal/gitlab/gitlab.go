package gitlab

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/xanzy/go-gitlab"

	"github.com/lindell/multi-gitter/internal/domain"
)

// New create a new Gitlab client
func New(token, baseURL string, repoListing RepositoryListing) (*Gitlab, error) {
	var options []gitlab.ClientOptionFunc
	if baseURL != "" {
		options = append(options, gitlab.WithBaseURL(baseURL))
	}

	client, err := gitlab.NewClient(token, options...)
	if err != nil {
		return nil, err
	}

	return &Gitlab{
		RepositoryListing: repoListing,
		glClient:          client,
	}, nil
}

// Gitlab contain gitlab configuration
type Gitlab struct {
	RepositoryListing
	glClient *gitlab.Client
}

// RepositoryListing contains information about which repositories that should be fetched
type RepositoryListing struct {
	Groups   []string
	Users    []string
	Projects []ProjectReference
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
	pid        int
	branchName string
	iid        int
	status     domain.PullRequestStatus
}

func (pr pullRequest) String() string {
	return fmt.Sprintf("%s/%s #%d", pr.ownerName, pr.repoName, pr.iid)
}

func (pr pullRequest) Status() domain.PullRequestStatus {
	return pr.status
}

// GetRepositories fetches repositories from and organization
func (g *Gitlab) GetRepositories(ctx context.Context) ([]domain.Repository, error) {
	allProjects, err := g.getProjects(ctx)
	if err != nil {
		return nil, err
	}

	repos := make([]domain.Repository, 0, len(allProjects))
	for _, project := range allProjects {
		u, err := url.Parse(project.HTTPURLToRepo)
		if err != nil {
			return nil, err // TODO: better error
		}

		repos = append(repos, repository{
			url:           *u,
			pid:           project.ID,
			name:          project.Path,
			ownerName:     project.Namespace.Path,
			defaultBranch: project.DefaultBranch,
		})
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
func (g *Gitlab) CreatePullRequest(ctx context.Context, repo domain.Repository, newPR domain.NewPullRequest) error {
	r := repo.(repository)

	// Convert from usernames to user ids
	var assigneeIDs []int
	if len(newPR.Reviewers) > 0 {
		var err error
		assigneeIDs, err = g.getUserIDs(ctx, newPR.Reviewers)
		if err != nil {
			return err
		}
	}

	_, _, err := g.glClient.MergeRequests.CreateMergeRequest(r.pid, &gitlab.CreateMergeRequestOptions{
		Title:        &newPR.Title,
		Description:  &newPR.Body,
		SourceBranch: &newPR.Head,
		TargetBranch: &newPR.Base,
		AssigneeIDs:  assigneeIDs,
	})
	return err
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

// GetPullRequestStatuses gets the statuses of all pull requests of with a specific branch name in an organization
func (g *Gitlab) GetPullRequestStatuses(ctx context.Context, branchName string) ([]domain.PullRequest, error) {
	projects, err := g.getProjects(ctx)
	if err != nil {
		return nil, err
	}

	prs := make([]domain.PullRequest, len(projects))
	for i, project := range projects {
		status, iid, err := g.getPullRequestInfo(ctx, branchName, project)
		if err != nil {
			return nil, err
		}

		prs[i] = pullRequest{
			repoName:   project.Path,
			ownerName:  project.Namespace.Path,
			pid:        project.ID,
			branchName: branchName,
			status:     status,
			iid:        iid,
		}
	}

	return prs, nil
}

func (g *Gitlab) getPullRequestInfo(ctx context.Context, branchName string, project *gitlab.Project) (status domain.PullRequestStatus, id int, err error) {
	mrs, _, err := g.glClient.MergeRequests.ListProjectMergeRequests(project.ID, &gitlab.ListProjectMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 1,
		},
		SourceBranch: &branchName,
	}, gitlab.WithContext(ctx))
	if err != nil {
		return domain.PullRequestStatusUnknown, mrs[0].IID, err
	}

	if len(mrs) == 0 {
		return domain.PullRequestStatusUnknown, 0, nil
	}

	mr, _, err := g.glClient.MergeRequests.GetMergeRequest(project.ID, mrs[0].IID, nil, gitlab.WithContext(ctx))
	if err != nil {
		return domain.PullRequestStatusUnknown, mrs[0].IID, err
	}

	switch {
	case mr.MergedAt != nil:
		return domain.PullRequestStatusMerged, mr.IID, nil
	case mr.ClosedAt != nil:
		return domain.PullRequestStatusClosed, mr.IID, nil
	case mr.Pipeline == nil, mr.Pipeline.Status == "success":
		return domain.PullRequestStatusSuccess, mr.IID, nil
	case mr.Pipeline.Status == "failed":
		return domain.PullRequestStatusError, mr.IID, nil
	default:
		return domain.PullRequestStatusPending, mr.IID, nil
	}
}

// MergePullRequest merges a pull request
func (g *Gitlab) MergePullRequest(ctx context.Context, pullReq domain.PullRequest) error {
	pr := pullReq.(pullRequest)

	shouldRemoveSourceBranch := true
	_, _, err := g.glClient.MergeRequests.AcceptMergeRequest(pr.pid, pr.iid, &gitlab.AcceptMergeRequestOptions{
		ShouldRemoveSourceBranch: &shouldRemoveSourceBranch,
	}, gitlab.WithContext(ctx))
	if err != nil {
		return err
	}

	return nil
}

package gitlab

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	internalHTTP "github.com/lindell/multi-gitter/internal/http"
	"github.com/lindell/multi-gitter/internal/scm"
	log "github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
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
		token:             token,
	}, nil
}

// Gitlab contain gitlab configuration
type Gitlab struct {
	RepositoryListing
	Config   Config
	token    string
	glClient *gitlab.Client

	// Cached current user
	currentUser *gitlab.User
}

// RepositoryListing contains information about which repositories that should be fetched
type RepositoryListing struct {
	Groups    []string
	Users     []string
	Projects  []ProjectReference
	Topics    []string
	SkipForks bool
}

// Config includes extra config parameters for the GitLab client
type Config struct {
	IncludeSubgroups bool
	SSHAuth          bool
}

// ProjectReference contains information to be able to reference a repository
type ProjectReference struct {
	OwnerName string
	Name      string
}

// ParseProjectReference parses a repository reference from the format "ownerName/repoName"
func ParseProjectReference(val string) (ProjectReference, error) {
	lastSlashIndex := strings.LastIndex(val, "/")
	if lastSlashIndex == -1 {
		return ProjectReference{}, fmt.Errorf("could not parse repository reference: %s", val)
	}
	return ProjectReference{
		OwnerName: val[:lastSlashIndex],
		Name:      val[lastSlashIndex+1:],
	}, nil
}

// GetRepositories fetches repositories from all sources (groups/user/specific project)
func (g *Gitlab) GetRepositories(ctx context.Context) ([]scm.Repository, error) {
	allProjects, err := g.getProjects(ctx)
	if err != nil {
		return nil, err
	}

	repos := make([]scm.Repository, 0, len(allProjects))
	for _, project := range allProjects {
		log := log.WithField("repo", project.NameWithNamespace)
		if len(g.Topics) != 0 && !scm.RepoContainsTopic(project.Topics, g.Topics) {
			log.Debug("Skipping repository since it does not match repository topics")
			continue
		}
		if g.SkipForks && project.ForkedFromProject != nil {
			log.Debug("Skipping repository since it's a fork")
			continue
		}

		p, err := g.convertProject(project)
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
	withMergeRequestsEnabled := true
	archived := false
	for i := 1; ; i++ {
		projects, _, err := g.glClient.Groups.ListGroupProjects(groupName, &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    i,
			},
			Archived:                 &archived,
			IncludeSubGroups:         &g.Config.IncludeSubgroups,
			WithMergeRequestsEnabled: &withMergeRequestsEnabled,
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
	archived := false
	for i := 1; ; i++ {
		projects, _, err := g.glClient.Projects.ListUserProjects(username, &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    i,
			},
			Archived: &archived,
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
func (g *Gitlab) CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
	r := repo.(repository)
	prR := prRepo.(repository)

	reviewersIDs, err := g.getUserIds(ctx, newPR.Reviewers)
	if err != nil {
		return nil, err
	}

	assigneesIDs, err := g.getUserIds(ctx, newPR.Assignees)
	if err != nil {
		return nil, err
	}

	prTitle := newPR.Title
	if newPR.Draft {
		prTitle = "Draft: " + prTitle // See https://docs.gitlab.com/ee/user/project/merge_requests/drafts.html#mark-merge-requests-as-drafts
	}

	labels := gitlab.LabelOptions(newPR.Labels)
	removeSourceBranch := true
	mr, _, err := g.glClient.MergeRequests.CreateMergeRequest(prR.pid, &gitlab.CreateMergeRequestOptions{
		Title:              &prTitle,
		Description:        &newPR.Body,
		SourceBranch:       &newPR.Head,
		TargetBranch:       &newPR.Base,
		TargetProjectID:    &r.pid,
		ReviewerIDs:        &reviewersIDs,
		RemoveSourceBranch: &removeSourceBranch,
		Squash:             &r.shouldSquash,
		AssigneeIDs:        &assigneesIDs,
		Labels:             &labels,
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

func (g *Gitlab) getUserIds(ctx context.Context, usernames []string) ([]int, error) {
	// Convert from usernames to user ids
	var assigneeIDs []int

	if len(usernames) > 0 {
		var err error
		assigneeIDs, err = g.getUserIDs(ctx, usernames)
		if err != nil {
			return nil, err
		}
	}

	return assigneeIDs, nil
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

// UpdatePullRequest updates an existing pull request
func (g *Gitlab) UpdatePullRequest(ctx context.Context, repo scm.Repository, pullReq scm.PullRequest, updatedPR scm.NewPullRequest) (scm.PullRequest, error) {
	r := repo.(repository)
	pr := pullReq.(pullRequest)

	reviewersIDs, err := g.getUserIds(ctx, updatedPR.Reviewers)
	if err != nil {
		return nil, err
	}

	assigneesIDs, err := g.getUserIds(ctx, updatedPR.Assignees)
	if err != nil {
		return nil, err
	}

	prTitle := updatedPR.Title
	if updatedPR.Draft {
		prTitle = "Draft: " + prTitle // See https://docs.gitlab.com/ee/user/project/merge_requests/drafts.html#mark-merge-requests-as-drafts
	}

	labels := gitlab.LabelOptions(updatedPR.Labels)
	mr, _, err := g.glClient.MergeRequests.UpdateMergeRequest(pr.sourcePID, pr.iid, &gitlab.UpdateMergeRequestOptions{
		Title:       &prTitle,
		Description: &updatedPR.Body,
		ReviewerIDs: &reviewersIDs,
		AssigneeIDs: &assigneesIDs,
		Labels:      &labels,
	})
	if err != nil {
		return nil, err
	}

	return pullRequest{
		repoName:   r.name,
		ownerName:  r.ownerName,
		targetPID:  mr.TargetProjectID,
		sourcePID:  mr.SourceProjectID,
		branchName: updatedPR.Head,
		iid:        mr.IID,
		webURL:     mr.WebURL,
	}, nil
}

// GetPullRequests gets all pull requests of with a specific branch
func (g *Gitlab) GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error) {
	projects, err := g.getProjects(ctx)
	if err != nil {
		return nil, err
	}

	prs := []scm.PullRequest{}
	for _, project := range projects {
		mr, err := g.getPullRequest(ctx, branchName, project)
		if err != nil {
			return nil, err
		}
		if mr == nil {
			continue
		}

		prs = append(prs, convertMergeRequest(mr, project.Path, project.Namespace.FullPath))
	}

	return prs, nil
}

func convertMergeRequest(mr *gitlab.MergeRequest, repoName, ownerName string) pullRequest {
	return pullRequest{
		repoName:   repoName,
		ownerName:  ownerName,
		targetPID:  mr.TargetProjectID,
		sourcePID:  mr.SourceProjectID,
		branchName: mr.SourceBranch,
		status:     pullRequestStatus(mr),
		iid:        mr.IID,
		webURL:     mr.WebURL,
	}
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

func pullRequestStatus(mr *gitlab.MergeRequest) scm.PullRequestStatus {
	switch {
	case mr.MergedAt != nil:
		return scm.PullRequestStatusMerged
	case mr.ClosedAt != nil:
		return scm.PullRequestStatusClosed
	case mr.Pipeline == nil, mr.Pipeline.Status == "success":
		return scm.PullRequestStatusSuccess
	case mr.Pipeline.Status == "failed":
		return scm.PullRequestStatusError
	default:
		return scm.PullRequestStatusPending
	}
}

// GetOpenPullRequest gets a pull request for one specific repository
func (g *Gitlab) GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
	project := repo.(repository)

	state := "opened"
	mrs, _, err := g.glClient.MergeRequests.ListProjectMergeRequests(project.pid, &gitlab.ListProjectMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 1,
		},
		SourceBranch: &branchName,
		State:        &state,
	}, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if len(mrs) == 0 {
		return nil, nil
	}

	return convertMergeRequest(mrs[0], project.name, project.ownerName), nil
}

// MergePullRequest merges a pull request
func (g *Gitlab) MergePullRequest(ctx context.Context, pullReq scm.PullRequest) error {
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
func (g *Gitlab) ClosePullRequest(ctx context.Context, pullReq scm.PullRequest) error {
	pr := pullReq.(pullRequest)

	stateEvent := "close"
	_, _, err := g.glClient.MergeRequests.UpdateMergeRequest(pr.targetPID, pr.iid, &gitlab.UpdateMergeRequestOptions{
		StateEvent: &stateEvent,
	}, gitlab.WithContext(ctx))
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
func (g *Gitlab) ForkRepository(ctx context.Context, repo scm.Repository, newOwner string) (scm.Repository, error) {
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
		return g.convertProject(project)
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
			return g.convertProject(newRepo)
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

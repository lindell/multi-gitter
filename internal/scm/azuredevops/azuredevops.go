package azuredevops

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/location"
)

type AzureDevOps struct {
	RepositoryListing
	gitClient       git.Client
	identityClient  identity.Client
	locationClient  location.Client
	coreClient      core.Client
	SSHAuth         bool
	pat             string
	projectNameToID map[string]uuid.UUID
}

// RepositoryListing contains information about which repositories that should be fetched
type RepositoryListing struct {
	Projects     []string
	Repositories []RepositoryReference
}

// RepositoryReference contains information to be able to reference a repository
type RepositoryReference struct {
	ProjectName string
	Name        string
}

// New creates a new AzureDevOps client
func New(pat, baseURL string, sshAuth bool, fork bool, repoListing RepositoryListing) (*AzureDevOps, error) {
	connection := azuredevops.NewPatConnection(baseURL, pat)
	gitClient, err := git.NewClient(context.Background(), connection)
	if err != nil {
		return nil, err
	}
	locationClient := location.NewClient(context.Background(), connection)

	coreClient, err := core.NewClient(context.Background(), connection)
	if err != nil {
		return nil, err
	}
	identityClient, err := identity.NewClient(context.Background(), connection)
	if err != nil {
		return nil, err
	}
	projectMap := make(map[string]uuid.UUID)

	if fork {
		projects, err := coreClient.GetProjects(context.Background(), core.GetProjectsArgs{})
		if err != nil {
			return nil, err
		}
		projectMap = make(map[string]uuid.UUID, len(projects.Value))

		for _, project := range projects.Value {
			projectMap[*project.Name] = *project.Id
		}
	}

	azureDevOps := &AzureDevOps{
		RepositoryListing: repoListing,
		gitClient:         gitClient,
		identityClient:    identityClient,
		locationClient:    locationClient,
		coreClient:        coreClient,
		SSHAuth:           sshAuth,
		pat:               pat,
		projectNameToID:   projectMap,
	}

	return azureDevOps, nil
}

func ParseRepositoryReference(repo string) (RepositoryReference, error) {
	split := strings.Split(repo, "/")
	if len(split) != 2 {
		return RepositoryReference{}, fmt.Errorf("could not parse repository reference: %s", repo)
	}

	return RepositoryReference{
		ProjectName: split[0],
		Name:        split[1],
	}, nil
}

// GetRepositories fetches repositories from all sources (groups/user/specific project)
func (a *AzureDevOps) GetRepositories(ctx context.Context) ([]scm.Repository, error) {
	projects, err := a.getRepos(ctx)

	if err != nil {
		return nil, err
	}
	repos := make([]scm.Repository, 0, len(projects))

	for _, repo := range projects {
		rCopy := repo
		r := a.convertRepo(&rCopy)
		repos = append(repos, r)
	}

	return repos, nil
}

func (a *AzureDevOps) getRepos(ctx context.Context) ([]git.GitRepository, error) {
	allRepos, err := a.gitClient.GetRepositories(ctx, git.GetRepositoriesArgs{})
	if err != nil {
		return nil, err
	}
	filteredRepos := make([]git.GitRepository, 0)

	for _, repo := range *allRepos {
		added := false

		// Check if the repository belongs to one of the configured projects
		for _, project := range a.RepositoryListing.Projects {
			if *repo.Project.Name == project {
				filteredRepos = append(filteredRepos, repo)
				added = true
				break
			}
		}

		if added {
			continue
		}

		// Check if the repository is one of the configured repositories
		for _, repository := range a.RepositoryListing.Repositories {
			if *repo.Project.Name == repository.ProjectName && *repo.Name == repository.Name {
				filteredRepos = append(filteredRepos, repo)
				break
			}
		}
	}

	return filteredRepos, nil
}

func (a *AzureDevOps) CreatePullRequest(ctx context.Context, _ scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error) {
	prr := prRepo.(repository)

	reviewerIDs, err := a.getUserIds(ctx, append(newPR.Reviewers, newPR.Assignees...))
	if err != nil {
		return nil, err
	}
	reviewers := make([]git.IdentityRefWithVote, len(reviewerIDs))
	for i, reviewerID := range reviewerIDs {
		reviewers[i] = git.IdentityRefWithVote{
			Id: &reviewerID,
		}
	}

	prTitle := newPR.Title
	if newPR.Draft {
		prTitle = "Draft: " + prTitle
	}

	activeLabel := true
	labels := make([]core.WebApiTagDefinition, len(newPR.Labels))
	removeSourceBranch := true

	for i, label := range newPR.Labels {
		lCopy := label
		labels[i] = core.WebApiTagDefinition{
			Active: &activeLabel,
			Name:   &lCopy,
		}
	}
	srn := PrependPrefixIfNeeded(newPR.Head) // Pass the value of newPR.Head
	trn := PrependPrefixIfNeeded(newPR.Base)

	requestArgs := git.CreatePullRequestArgs{
		GitPullRequestToCreate: &git.GitPullRequest{
			Title:         &prTitle,
			Description:   &newPR.Body,
			SourceRefName: &srn,
			TargetRefName: &trn,
			CompletionOptions: &git.GitPullRequestCompletionOptions{
				DeleteSourceBranch: &removeSourceBranch,
			},
			Labels:    &labels,
			Reviewers: &reviewers,
			IsDraft:   &newPR.Draft,
		},
		RepositoryId: &prr.rid,
		Project:      &prr.ownerName,
	}

	pr, err := a.gitClient.CreatePullRequest(ctx, requestArgs)
	if err != nil {
		return nil, err
	}
	return pullRequest{
		ownerName:  prr.ownerName,
		repoName:   prr.name,
		repoID:     prr.rid,
		branchName: newPR.Head,
		id:         *pr.PullRequestId,
	}, nil
}

func (a *AzureDevOps) UpdatePullRequest(ctx context.Context, repo scm.Repository, pullReq scm.PullRequest, updatedPR scm.NewPullRequest) (scm.PullRequest, error) {
	r := repo.(repository)
	pr := pullReq.(pullRequest)

	reviewerIDs, err := a.getUserIds(ctx, append(updatedPR.Reviewers, updatedPR.Assignees...))
	if err != nil {
		return nil, err
	}
	reviewers := make([]git.IdentityRefWithVote, len(reviewerIDs))
	for i, reviewerID := range reviewerIDs {
		rCopy := reviewerID
		reviewers[i] = git.IdentityRefWithVote{
			Id: &rCopy,
		}
	}
	prTitle := updatedPR.Title
	if updatedPR.Draft {
		prTitle = "Draft: " + prTitle
	}

	activeLabel := true
	labels := make([]core.WebApiTagDefinition, len(updatedPR.Labels))

	for i, label := range updatedPR.Labels {
		lCopy := label
		labels[i] = core.WebApiTagDefinition{
			Active: &activeLabel,
			Name:   &lCopy,
		}
	}

	srn := PrependPrefixIfNeeded(updatedPR.Head) // Pass the value of newPR.Head
	trn := PrependPrefixIfNeeded(updatedPR.Base)

	updateResult, err := a.gitClient.UpdatePullRequest(ctx, git.UpdatePullRequestArgs{
		GitPullRequestToUpdate: &git.GitPullRequest{
			Title:         &prTitle,
			Description:   &updatedPR.Body,
			SourceRefName: &srn,
			TargetRefName: &trn,
			Labels:        &labels,
			Reviewers:     &reviewers,
		},
		RepositoryId:  &pr.repoID,
		PullRequestId: &pr.id,
	})

	if err != nil {
		return nil, err
	}
	return pullRequest{
		ownerName:  r.ownerName,
		repoName:   r.name,
		repoID:     r.rid,
		branchName: updatedPR.Head,
		id:         *updateResult.PullRequestId,
	}, nil
}

func (a *AzureDevOps) getUserIds(ctx context.Context, usernames []string) ([]string, error) {
	userIds := make([]string, len(usernames))

	searchFilter := "General"

	for i, username := range usernames {
		uCopy := username
		users, err := a.identityClient.ReadIdentities(ctx, identity.ReadIdentitiesArgs{
			SearchFilter: &searchFilter,
			FilterValue:  &uCopy,
		})

		if err != nil {
			return nil, err
		}
		if len(*users) != 1 {
			return nil, fmt.Errorf("could not find user: %s", usernames[i])
		}

		userIds[i] = (*users)[0].Id.String()
	}
	return userIds, nil
}

func (a *AzureDevOps) GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error) {
	repos, err := a.getRepos(ctx)
	if err != nil {
		return nil, err
	}
	bn := PrependPrefixIfNeeded(branchName)

	prs := []scm.PullRequest{}
	for _, repo := range repos {
		pr, err := a.getPullRequest(ctx, bn, repo.Id.String())
		if err != nil {
			return nil, err
		}
		if pr == nil {
			continue
		}

		prs = append(prs, convertPullRequest(*pr))
	}

	return prs, nil
}

func convertPullRequest(pr git.GitPullRequest) pullRequest {
	return pullRequest{
		ownerName:               *pr.Repository.Project.Name,
		repoName:                *pr.Repository.Name,
		repoID:                  pr.Repository.Id.String(),
		branchName:              *pr.SourceRefName,
		id:                      *pr.PullRequestId,
		status:                  pullRequestStatus(pr),
		lastMergeSourceCommitID: *pr.LastMergeSourceCommit.CommitId,
	}
}

func (a *AzureDevOps) getPullRequest(ctx context.Context, branchName string, rid string) (*git.GitPullRequest, error) {
	status := &git.PullRequestStatusValues.Active

	prs, err := a.gitClient.GetPullRequests(ctx, git.GetPullRequestsArgs{
		RepositoryId: &rid,
		SearchCriteria: &git.GitPullRequestSearchCriteria{
			SourceRefName: &branchName,
			Status:        status,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(*prs) == 0 {
		return nil, nil
	}

	return &(*prs)[0], nil
}

func pullRequestStatus(pr git.GitPullRequest) scm.PullRequestStatus {
	switch {
	case *pr.MergeStatus == git.PullRequestAsyncStatusValues.Succeeded:
		return scm.PullRequestStatusSuccess
	case *pr.Status == git.PullRequestStatusValues.Abandoned || *pr.Status == git.PullRequestStatusValues.Completed:
		return scm.PullRequestStatusClosed
	case *pr.MergeStatus == git.PullRequestAsyncStatusValues.Conflicts ||
		*pr.MergeStatus == git.PullRequestAsyncStatusValues.Failure ||
		*pr.MergeStatus == git.PullRequestAsyncStatusValues.RejectedByPolicy:
		return scm.PullRequestStatusError
	case *pr.Status == git.PullRequestStatusValues.Active:
		return scm.PullRequestStatusPending
	default:
		return scm.PullRequestStatusUnknown
	}
}

// GetOpenPullRequest gets a pull request for one specific repository
func (a *AzureDevOps) GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error) {
	r := repo.(repository)

	prs, err := a.gitClient.GetPullRequests(ctx, git.GetPullRequestsArgs{
		RepositoryId: &r.rid,
		SearchCriteria: &git.GitPullRequestSearchCriteria{
			Status:        &git.PullRequestStatusValues.Active,
			SourceRefName: &branchName,
		},
	})
	if err != nil {
		return nil, err
	}
	if len(*prs) == 0 {
		return nil, nil
	}

	return convertPullRequest((*prs)[0]), nil
}

// MergePullRequest merges a pull request
func (a *AzureDevOps) MergePullRequest(ctx context.Context, pullReq scm.PullRequest) error {
	pr := pullReq.(pullRequest)

	_, err := a.gitClient.UpdatePullRequest(ctx, git.UpdatePullRequestArgs{
		GitPullRequestToUpdate: &git.GitPullRequest{
			Status: &git.PullRequestStatusValues.Completed,
			LastMergeSourceCommit: &git.GitCommitRef{
				CommitId: &pr.lastMergeSourceCommitID,
			},
		},
		RepositoryId:  &pr.repoID,
		PullRequestId: &pr.id,
	})
	if err != nil {
		return err
	}
	return nil
}

// ClosePullRequest closes a pull request
func (a *AzureDevOps) ClosePullRequest(ctx context.Context, pullReq scm.PullRequest) error {
	pr := pullReq.(pullRequest)

	_, err := a.gitClient.UpdatePullRequest(ctx, git.UpdatePullRequestArgs{
		GitPullRequestToUpdate: &git.GitPullRequest{
			Status: &git.PullRequestStatusValues.Abandoned,
		},
		RepositoryId:  &pr.repoID,
		PullRequestId: &pr.id,
	})
	if err != nil {
		return err
	}
	newObjectID := "0000000000000000000000000000000000000000"
	refUpdate := git.GitRefUpdate{
		Name:        &pr.branchName,
		OldObjectId: &pr.lastMergeSourceCommitID,
		NewObjectId: &newObjectID,
	}

	refUpdates := []git.GitRefUpdate{refUpdate}

	_, err = a.gitClient.UpdateRefs(ctx, git.UpdateRefsArgs{
		RepositoryId: &pr.repoID,
		RefUpdates:   &refUpdates,
	})

	if err != nil {
		return err
	}
	return nil
}

// ForkRepository forks a project
func (a *AzureDevOps) ForkRepository(ctx context.Context, repo scm.Repository, newProject string) (scm.Repository, error) {
	r := repo.(repository)

	currentUser, err := a.getCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	np := newProject
	if np == "" {
		np = r.ownerName
	}

	repoName := r.name + "." + strings.ReplaceAll(currentUser, " ", ".")

	// Check if the forked repo already exists
	existingFork, err := a.gitClient.GetRepository(ctx, git.GetRepositoryArgs{
		Project:      &np,
		RepositoryId: &repoName,
	})

	if err == nil {
		return a.convertRepo(existingFork), nil
	}

	repoUUID, err := uuid.Parse(r.rid)
	if err != nil {
		return nil, err
	}
	pid := a.projectNameToID[np]
	parentPid := a.projectNameToID[r.ownerName]

	// Fork the repository
	forkedRepo, err := a.gitClient.CreateRepository(ctx, git.CreateRepositoryArgs{
		GitRepositoryToCreate: &git.GitRepositoryCreateOptions{
			Name: &repoName,
			Project: &core.TeamProjectReference{
				Id: &pid,
			},
			ParentRepository: &git.GitRepositoryRef{
				Id: &repoUUID,
				Project: &core.TeamProjectReference{
					Id: &parentPid,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	// Convert the forked repository to scm.Repository
	return a.convertRepo(forkedRepo), nil
}

func (a *AzureDevOps) getCurrentUser(ctx context.Context) (string, error) {
	currentUser, err := a.locationClient.GetConnectionData(ctx, location.GetConnectionDataArgs{})
	if err != nil {
		return "", err
	}
	return *currentUser.AuthenticatedUser.ProviderDisplayName, nil
}

func PrependPrefixIfNeeded(s string) string {
	if !strings.HasPrefix(s, "refs/heads/") {
		s = "refs/heads/" + s
	}
	return s
}

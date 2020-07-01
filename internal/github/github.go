package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/lindell/multi-gitter/internal/domain"
)

// New create a new Github client
func New(token, baseURL string, repoListing RepositoryListing) (*Github, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	tc.Transport = loggingRoundTripper{
		next: tc.Transport,
	}

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
		ghClient:          client,
	}, nil
}

// Github contain github configuration
type Github struct {
	RepositoryListing
	ghClient *github.Client
}

// RepositoryListing contains information about which repositories that should be fetched
type RepositoryListing struct {
	Organizations []string
	Users         []string
}

type pullRequest struct {
	ID     int64 `json:"id"`
	Number int   `json:"number"`
}

// GetRepositories fetches repositories from and organization
func (g Github) GetRepositories(ctx context.Context) ([]domain.Repository, error) {
	allRepos, err := g.getRepositories(ctx)
	if err != nil {
		return nil, err
	}

	repos := make([]domain.Repository, 0, len(allRepos))
	for _, r := range allRepos {
		permissions := r.GetPermissions()
		if !r.GetArchived() && !r.GetDisabled() && permissions["pull"] && permissions["push"] {
			repos = append(repos, domain.Repository{
				URL:           r.GetCloneURL(),
				Name:          r.GetName(),
				OwnerName:     r.GetOwner().GetLogin(),
				DefaultBranch: r.GetDefaultBranch(),
			})
		}
	}

	return repos, nil
}

func (g Github) getRepositories(ctx context.Context) ([]*github.Repository, error) {
	allRepos := []*github.Repository{}

	for _, org := range g.Organizations {
		repos, err := g.getOrganizationRepositories(ctx, org)
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

	return allRepos, nil
}

// GetRepositories fetches repositories from and organization
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

// CreatePullRequest creates a pull request
func (g Github) CreatePullRequest(ctx context.Context, repository domain.Repository, newPR domain.NewPullRequest) error {
	pr, err := g.createPullRequest(ctx, repository, newPR)
	if err != nil {
		return err
	}

	if err := g.addReviewers(ctx, repository, newPR, pr); err != nil {
		return err
	}

	return nil
}

func (g Github) createPullRequest(ctx context.Context, repo domain.Repository, newPR domain.NewPullRequest) (pullRequest, error) {
	pr, _, err := g.ghClient.PullRequests.Create(ctx, repo.OwnerName, repo.Name, &github.NewPullRequest{
		Title: &newPR.Title,
		Body:  &newPR.Body,
		Head:  &newPR.Head,
		Base:  &newPR.Base,
	})
	if err != nil {
		return pullRequest{}, err
	}

	return pullRequest{
		ID:     pr.GetID(),
		Number: pr.GetNumber(),
	}, nil
}

func (g Github) addReviewers(ctx context.Context, repo domain.Repository, newPR domain.NewPullRequest, createdPR pullRequest) error {
	if len(newPR.Reviewers) == 0 {
		return nil
	}
	_, _, err := g.ghClient.PullRequests.RequestReviewers(ctx, repo.OwnerName, repo.Name, createdPR.Number, github.ReviewersRequest{
		Reviewers: newPR.Reviewers,
	})
	return err
}

// GetPullRequestStatuses gets the statuses of all pull requests of with a specific branch name in an organization
func (g Github) GetPullRequestStatuses(ctx context.Context, branchName string) ([]domain.PullRequest, error) {
	// TODO: If this is implemented with the GitHub v4 graphql api, it would be much faster

	repos, err := g.getRepositories(ctx)
	if err != nil {
		return nil, err
	}

	prStatuses := []domain.PullRequest{}
	for _, r := range repos {
		repoOwner := r.GetOwner().GetLogin()
		repoName := r.GetName()
		log := log.WithField("repo", fmt.Sprintf("%s/%s", repoOwner, repoName))
		log.Debug("Fetching latest pull request")
		prs, _, err := g.ghClient.PullRequests.List(ctx, repoOwner, repoName, &github.PullRequestListOptions{
			Head:      fmt.Sprintf("%s:%s", repoOwner, branchName),
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

		// Determine the status of the pr
		var status domain.PullRequestStatus
		if pr.MergedAt != nil {
			status = domain.PullRequestStatusMerged
		} else if pr.ClosedAt != nil {
			status = domain.PullRequestStatusClosed
		} else {
			log.Debug("Fetching the combined status of the pull request")
			combinedStatus, _, err := g.ghClient.Repositories.GetCombinedStatus(ctx, repoOwner, repoName, pr.GetHead().GetSHA(), nil)
			if err != nil {
				return nil, err
			}

			if combinedStatus.GetTotalCount() == 0 {
				status = domain.PullRequestStatusSuccess
			} else {
				switch combinedStatus.GetState() {
				case "pending":
					status = domain.PullRequestStatusPending
				case "success":
					status = domain.PullRequestStatusSuccess
				case "failure", "error":
					status = domain.PullRequestStatusError
				}
			}
		}

		prStatuses = append(prStatuses, domain.PullRequest{
			OwnerName:  repoOwner,
			RepoName:   repoName,
			BranchName: pr.GetHead().GetRef(),
			Number:     pr.GetNumber(),
			Status:     status,
		})
	}

	return prStatuses, nil
}

// MergePullRequest merges a pull request
func (g Github) MergePullRequest(ctx context.Context, pr domain.PullRequest) error {
	_, _, err := g.ghClient.PullRequests.Merge(ctx, pr.OwnerName, pr.RepoName, pr.Number, "", nil)
	if err != nil {
		return err
	}

	_, err = g.ghClient.Git.DeleteRef(ctx, pr.OwnerName, pr.RepoName, fmt.Sprintf("heads/%s", pr.BranchName))
	return err
}

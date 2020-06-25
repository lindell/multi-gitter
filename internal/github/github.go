package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"

	"github.com/lindell/multi-gitter/internal/domain"
)

// New create a new Github client
func New(token, baseURL string) (*Github, error) {
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
		ghClient: client,
	}, nil
}

// Github contain github configuration
type Github struct {
	ghClient *github.Client
}

type pullRequest struct {
	ID     int64 `json:"id"`
	Number int   `json:"number"`
}

// GetRepositories fetches repositories from and organization
func (g Github) GetRepositories(ctx context.Context, orgName string) ([]domain.Repository, error) {
	allRepos := []domain.Repository{}
	for i := 1; ; i++ {
		repos, err := g.getRepositories(ctx, orgName, i)
		if err != nil {
			return nil, err
		} else if len(repos) == 0 {
			break
		}
		allRepos = append(allRepos, repos...)
	}
	return allRepos, nil
}

func (g Github) getRepositories(ctx context.Context, orgName string, page int) ([]domain.Repository, error) {
	rr, _, err := g.ghClient.Repositories.ListByOrg(ctx, orgName, &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: 100,
		},
	})
	if err != nil {
		return nil, err
	}

	repos := make([]domain.Repository, 0, len(rr))
	for _, r := range rr {
		if !r.GetArchived() && !r.GetDisabled() {
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
func (g Github) GetPullRequestStatuses(ctx context.Context, orgName, branchName string) ([]domain.PullRequest, error) {
	// TODO: If this is implemented with the GitHub v4 graphql api, it would be much faster

	var repos []*github.Repository
	for {
		rr, _, err := g.ghClient.Repositories.ListByOrg(ctx, orgName, &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{
				Page:    1,
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
	}

	prStatuses := []domain.PullRequest{}
	for _, r := range repos {
		prs, _, err := g.ghClient.PullRequests.List(ctx, orgName, r.GetName(), &github.PullRequestListOptions{
			Head:      fmt.Sprintf("%s:%s", orgName, branchName),
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
			combinedStatus, _, err := g.ghClient.Repositories.GetCombinedStatus(ctx, orgName, r.GetName(), pr.GetHead().GetSHA(), nil)
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
			OwnerName:  r.GetOwner().GetLogin(),
			RepoName:   r.GetName(),
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

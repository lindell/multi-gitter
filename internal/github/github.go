package github

import (
	"context"
	"errors"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"

	"github.com/lindell/multi-gitter/internal/domain"
)

func New(token, baseURL string) (*Github, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

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

type repository struct {
	SSH           string
	Name          string
	OwnerName     string
	DefaultBranch string
}

func (r repository) GetURL() string {
	return r.SSH
}

func (r repository) GetBranch() string {
	return r.DefaultBranch
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
			repos = append(repos, repository{
				SSH:           r.GetSSHURL(),
				Name:          r.GetName(),
				OwnerName:     r.GetOwner().GetLogin(),
				DefaultBranch: r.GetDefaultBranch(),
			})
		}
	}
	return repos, nil
}

// CreatePullRequest creates a pull request
func (g Github) CreatePullRequest(ctx context.Context, repo domain.Repository, newPR domain.NewPullRequest) error {
	repository, ok := repo.(repository)
	if !ok {
		return errors.New("the repository needs to originate from this package")
	}

	pr, err := g.createPullRequest(ctx, repository, newPR)
	if err != nil {
		return err
	}

	if err := g.addReviewers(ctx, repository, newPR, pr); err != nil {
		return err
	}

	return nil
}

func (g Github) createPullRequest(ctx context.Context, repo repository, newPR domain.NewPullRequest) (pullRequest, error) {
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

func (g Github) addReviewers(ctx context.Context, repo repository, newPR domain.NewPullRequest, createdPR pullRequest) error {
	_, _, err := g.ghClient.PullRequests.RequestReviewers(ctx, repo.OwnerName, repo.Name, createdPR.Number, github.ReviewersRequest{
		Reviewers: newPR.Reviewers,
	})
	return err
}

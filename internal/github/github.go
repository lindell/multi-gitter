package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/lindell/multi-gitter/internal/domain"
)

type Config struct {
	BaseURL string
	Token   string // Personal access token
}

// DefaultConfig contains values for the github.com api
// The access token is still always needed
var DefaultConfig = Config{
	BaseURL: "https://api.github.com/",
}

type createPrRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

type pr struct {
	ID     int `json:"id"`
	Number int `json:"number"`
}

type addReviewersRequest struct {
	Reviewers []string `json:"reviewers"`
}

type repository struct {
	SSH           string `json:"ssh_url"`
	Slug          string `json:"full_name"`
	DefaultBranch string `json:"default_branch"`
}

func (r repository) GetURL() string {
	return r.SSH
}

func (r repository) GetBranch() string {
	return r.DefaultBranch
}

type OrgRepoGetter struct {
	Config

	Organization string
}

func (g OrgRepoGetter) GetRepositories() ([]domain.Repository, error) {
	url := fmt.Sprintf("%sorgs/%s/repos", g.BaseURL, g.Organization)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "token "+g.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unknown status code %d when fetching repositories from %s", resp.StatusCode, url)
	}

	var rr []repository
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		return nil, err
	}

	// Transform the slice of repositories struct into a slice of the interface repositories
	repos := make([]domain.Repository, len(rr))
	for i := range rr {
		repos[i] = rr[i]
	}
	return repos, nil
}

type PullRequestCreator struct {
	Config
}

func (g PullRequestCreator) CreatePullRequest(repo domain.Repository, newPR domain.NewPullRequest) error {
	repository, ok := repo.(repository)
	if !ok {
		return errors.New("the repository needs to originate from this package")
	}

	pr, err := g.createPullRequest(repository, newPR)
	if err != nil {
		return err
	}

	if err := g.addReviewers(repository, newPR, pr); err != nil {
		return err
	}

	return nil
}

func (g PullRequestCreator) createPullRequest(repo repository, newPR domain.NewPullRequest) (pr, error) {
	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(createPrRequest{
		Title: newPR.Title,
		Body:  newPR.Body,
		Head:  newPR.Head,
		Base:  newPR.Base,
	})

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%srepos/%s/pulls", g.BaseURL, repo.Slug), buf)
	if err != nil {
		return pr{}, err
	}
	req.Header.Add("Authorization", "token "+g.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return pr{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return pr{}, fmt.Errorf("could not create pull request, got status code: %d", resp.StatusCode)
	}

	var pullRequest pr
	if err := json.NewDecoder(resp.Body).Decode(&pullRequest); err != nil {
		return pr{}, err
	}
	return pullRequest, nil
}

func (g PullRequestCreator) addReviewers(repo repository, newPR domain.NewPullRequest, createdPR pr) error {
	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(addReviewersRequest{
		Reviewers: newPR.Reviewers,
	})

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%srepos/%s/pulls/%d/requested_reviewers", g.BaseURL, repo.Slug, createdPR.Number), buf)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "token "+g.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("could add reviewers to pull request, got status code: %d", resp.StatusCode)
	}

	return nil
}

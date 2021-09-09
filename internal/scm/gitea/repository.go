package gitea

import (
	"fmt"
	"net/url"

	"code.gitea.io/sdk/gitea"
)

func convertRepository(repo *gitea.Repository, token string) (repository, error) {
	u, err := url.Parse(repo.CloneURL)
	if err != nil {
		return repository{}, err
	}

	return repository{
		url:           *u,
		name:          repo.Name,
		ownerName:     repo.Owner.UserName,
		defaultBranch: repo.DefaultBranch,
		token:         token,
	}, nil
}

type repository struct {
	url           url.URL
	name          string
	ownerName     string
	defaultBranch string
	token         string
}

func (r repository) CloneURL() string {
	// Set the token as https://oauth2:TOKEN@url
	r.url.User = url.UserPassword("oauth2", r.token)
	return r.url.String()
}

func (r repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r repository) FullName() string {
	return fmt.Sprintf("%s/%s", r.ownerName, r.name)
}

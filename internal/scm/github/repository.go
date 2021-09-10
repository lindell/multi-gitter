package github

import (
	"fmt"
	"net/url"

	"github.com/google/go-github/v38/github"
	"github.com/pkg/errors"
)

func convertRepo(r *github.Repository, token string) (repository, error) {
	u, err := url.Parse(r.GetCloneURL())
	if err != nil {
		return repository{}, errors.Wrap(err, "could not parse github clone error")
	}

	return repository{
		url:           *u,
		name:          r.GetName(),
		ownerName:     r.GetOwner().GetLogin(),
		defaultBranch: r.GetDefaultBranch(),
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
	// Set the token as https://TOKEN@url
	r.url.User = url.User(r.token)
	return r.url.String()
}

func (r repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r repository) FullName() string {
	return fmt.Sprintf("%s/%s", r.ownerName, r.name)
}

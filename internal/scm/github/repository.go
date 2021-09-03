package github

import (
	"fmt"
	"net/url"

	"github.com/google/go-github/v38/github"
	"github.com/pkg/errors"
)

func convertRepo(r *github.Repository) (Repository, error) {
	u, err := url.Parse(r.GetCloneURL())
	if err != nil {
		return Repository{}, errors.Wrap(err, "could not parse github clone error")
	}

	return Repository{
		url:           *u,
		name:          r.GetName(),
		ownerName:     r.GetOwner().GetLogin(),
		defaultBranch: r.GetDefaultBranch(),
	}, nil
}

type Repository struct {
	url           url.URL
	name          string
	ownerName     string
	defaultBranch string
}

func (r Repository) URL(token string) string {
	// Set the token as https://TOKEN@url
	r.url.User = url.User(token)
	return r.url.String()
}

func (r Repository) URLWithUsername(_, _ string) string {
	return ""
}

func (r Repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r Repository) FullName() string {
	return fmt.Sprintf("%s/%s", r.ownerName, r.name)
}

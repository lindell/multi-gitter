package github

import (
	"fmt"
	"net/url"

	"github.com/google/go-github/v58/github"
	"github.com/pkg/errors"
)

func (g *Github) convertRepo(r *github.Repository) (repository, error) {
	var repoURL string
	if g.SSHAuth {
		repoURL = r.GetSSHURL()
	} else {
		u, err := url.Parse(r.GetCloneURL())
		if err != nil {
			return repository{}, errors.Wrap(err, "could not parse github clone error")
		}
		// Set the token as https://oauth2@TOKEN@url
		u.User = url.UserPassword("oauth2", g.token)
		repoURL = u.String()
	}

	return repository{
		url:           repoURL,
		name:          r.GetName(),
		ownerName:     r.GetOwner().GetLogin(),
		defaultBranch: r.GetDefaultBranch(),
	}, nil
}

type repository struct {
	url           string
	name          string
	ownerName     string
	defaultBranch string
}

func (r repository) CloneURL() string {
	return r.url
}

func (r repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r repository) FullName() string {
	return fmt.Sprintf("%s/%s", r.ownerName, r.name)
}

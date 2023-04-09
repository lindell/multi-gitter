package gitea

import (
	"fmt"
	"net/url"

	"code.gitea.io/sdk/gitea"
)

func (g *Gitea) convertRepository(repo *gitea.Repository) (repository, error) {
	var repoURL string
	if g.SSHAuth {
		repoURL = repo.SSHURL
	} else {
		u, err := url.Parse(repo.CloneURL)
		if err != nil {
			return repository{}, err
		}
		// Set the token as https://oauth2:TOKEN@url
		u.User = url.UserPassword("oauth2", g.token)
		repoURL = u.String()
	}

	return repository{
		url:           repoURL,
		name:          repo.Name,
		ownerName:     repo.Owner.UserName,
		defaultBranch: repo.DefaultBranch,
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

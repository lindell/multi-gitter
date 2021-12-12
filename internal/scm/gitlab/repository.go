package gitlab

import (
	"fmt"
	"net/url"

	"github.com/xanzy/go-gitlab"
)

func (g *Gitlab) convertProject(project *gitlab.Project) (repository, error) {
	var cloneURL string
	if g.Config.SSHAuth {
		cloneURL = project.SSHURLToRepo
	} else {
		u, err := url.Parse(project.HTTPURLToRepo)
		if err != nil {
			return repository{}, err
		}
		u.User = url.UserPassword("oauth2", g.token)
		cloneURL = u.String()
	}

	return repository{
		url:           cloneURL,
		pid:           project.ID,
		name:          project.Path,
		ownerName:     project.Namespace.Path,
		defaultBranch: project.DefaultBranch,
	}, nil
}

type repository struct {
	url           string
	pid           int
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

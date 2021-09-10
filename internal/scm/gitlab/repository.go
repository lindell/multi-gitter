package gitlab

import (
	"fmt"
	"net/url"

	"github.com/xanzy/go-gitlab"
)

func convertProject(project *gitlab.Project, token string) (repository, error) {
	u, err := url.Parse(project.HTTPURLToRepo)
	if err != nil {
		return repository{}, err
	}

	return repository{
		url:           *u,
		pid:           project.ID,
		name:          project.Path,
		ownerName:     project.Namespace.Path,
		defaultBranch: project.DefaultBranch,
		token:         token,
	}, nil
}

type repository struct {
	url           url.URL
	pid           int
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

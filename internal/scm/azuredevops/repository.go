package azuredevops

import (
	"fmt"
	"net/url"
	"path"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/pkg/errors"
)

func (a *AzureDevOps) convertRepo(repo *git.GitRepository) (repository, error) {
	var cloneURL string
	if a.SSHAuth {
		cloneURL = *repo.SshUrl
	} else {
		u, err := url.Parse(*repo.RemoteUrl)
		if err != nil {
			return repository{}, errors.Wrap(err, "could not parse Azure Devops remote url")
		}
		// Set the token as https://ACCESSTOKEN@remote-url
		u.User = url.User(a.pat)
		cloneURL = u.String()
	}

	defaultBranch := ""
	if repo.DefaultBranch != nil {
		defaultBranch = path.Base(*repo.DefaultBranch)
	}

	return repository{
		url:           cloneURL,
		rid:           repo.Id.String(),
		name:          *repo.Name,
		ownerName:     *repo.Project.Name,
		defaultBranch: defaultBranch,
	}, nil
}

type repository struct {
	url           string
	rid           string
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

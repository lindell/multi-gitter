package azuredevops

import (
	"fmt"
	"path"
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
)

func (a *AzureDevOps) convertRepo(repo *git.GitRepository) repository {
	var cloneURL string
	if a.SSHAuth {
		cloneURL = *repo.SshUrl
	} else {
		cloneURL = fmt.Sprintf("https://%s@%s", a.pat, strings.TrimPrefix(*repo.RemoteUrl, "https://"))
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
	}
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

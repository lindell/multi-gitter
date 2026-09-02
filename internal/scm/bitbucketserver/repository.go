package bitbucketserver

import (
	"net/url"
	"path"
	"strings"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/pkg/errors"
)

func (b *BitbucketServer) convertRepository(bitbucketRepository *bitbucketv1.Repository, defaultBranch bitbucketv1.Branch) (*repository, error) {
	var cloneURL string

	if b.sshAuth {
		cloneURL = findLinkType(bitbucketRepository.Links.Clone, cloneSSHType)
		if cloneURL == "" {
			return nil, errors.Errorf("unable to find clone url for repository %s using clone type %s", bitbucketRepository.Name, cloneSSHType)
		}
	} else {
		httpURL := findLinkType(bitbucketRepository.Links.Clone, cloneHTTPType)
		if httpURL == "" {
			return nil, errors.Errorf("unable to find clone url for repository %s using clone type %s", bitbucketRepository.Name, cloneHTTPType)
		}
		parsedURL, err := url.Parse(httpURL)
		if err != nil {
			return nil, err
		}

		parsedURL.User = url.UserPassword(b.username, b.token)
		cloneURL = parsedURL.String()
	}

	repo := repository{
		name:          bitbucketRepository.Slug,
		project:       bitbucketRepository.Project.Key,
		defaultBranch: defaultBranch.DisplayID,
		cloneURL:      cloneURL,
		webURL:        b.repositoryWebURL(bitbucketRepository.Project.Key, bitbucketRepository.Slug),
	}

	return &repo, nil
}

func findLinkType(links []bitbucketv1.CloneLink, cloneType string) string {
	for _, clone := range links {
		if strings.EqualFold(clone.Name, cloneType) {
			return clone.Href
		}
	}

	return ""
}

// repository contains information about a bitbucket repository
type repository struct {
	name          string
	project       string
	defaultBranch string
	cloneURL      string
	webURL        string
}

func (r repository) CloneURL() string {
	return r.cloneURL
}

func (r repository) BranchURL(branchName string) string {
	if r.webURL == "" {
		return ""
	}

	return r.webURL + "/browse?at=" + url.QueryEscape("refs/heads/"+branchName)
}

func (r repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r repository) FullName() string {
	return r.project + "/" + r.name
}

func (b *BitbucketServer) repositoryWebURL(projectKey, slug string) string {
	if b.baseURL == nil {
		return ""
	}

	browserURL := *b.baseURL
	browserURL.Path = strings.TrimSuffix(browserURL.Path, "/rest")
	browserURL.Path = path.Join(browserURL.Path, "projects", projectKey, "repos", slug)
	browserURL.RawPath = ""
	browserURL.RawQuery = ""
	browserURL.Fragment = ""

	return browserURL.String()
}

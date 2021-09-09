package bitbucketserver

import (
	"net/url"
	"strings"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/pkg/errors"
)

func convertRepository(bitbucketRepository *bitbucketv1.Repository, defaultBranch bitbucketv1.Branch, username, token string) (*repository, error) {
	var cloneURL *url.URL
	var err error
	for _, clone := range bitbucketRepository.Links.Clone {
		if strings.EqualFold(clone.Name, cloneType) {
			cloneURL, err = url.Parse(clone.Href)
			if err != nil {
				return nil, err
			}

			break
		}
	}

	if cloneURL == nil {
		return nil, errors.Errorf("unable to find clone url for repostory %s using clone type %s", bitbucketRepository.Name, cloneType)
	}

	repo := repository{
		name:          bitbucketRepository.Slug,
		project:       bitbucketRepository.Project.Key,
		defaultBranch: defaultBranch.DisplayID,
		username:      username,
		token:         token,
		cloneURL:      cloneURL,
	}

	return &repo, nil
}

// repository contains information about a bitbucket repository
type repository struct {
	name          string
	project       string
	defaultBranch string
	username      string
	token         string
	cloneURL      *url.URL
}

func (r repository) CloneURL() string {
	cloneURL := r.cloneURL

	if r.username != "" && r.token != "" {
		cloneURL.User = url.UserPassword(r.username, r.token)
	}

	return cloneURL.String()
}

func (r repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r repository) FullName() string {
	return r.project + "/" + r.name
}

package bitbucketserver

import "net/url"

// Repository contains information about a bitbucket Repository
type Repository struct {
	name          string
	project       string
	defaultBranch string
	cloneURL      *url.URL
}

func (r Repository) URL(_ string) string {
	return ""
}

func (r Repository) URLWithUsername(username, token string) string {
	cloneURL := r.cloneURL

	if username != "" && token != "" {
		cloneURL.User = url.UserPassword(username, token)
	}

	return cloneURL.String()
}

func (r Repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r Repository) FullName() string {
	return r.project + "/" + r.name
}

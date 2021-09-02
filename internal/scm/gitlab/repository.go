package gitlab

import (
	"fmt"
	"net/url"
)

type Repository struct {
	url           url.URL
	pid           int
	name          string
	ownerName     string
	defaultBranch string
}

func (r Repository) URL(token string) string {
	// Set the token as https://oauth2:TOKEN@url
	r.url.User = url.UserPassword("oauth2", token)
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

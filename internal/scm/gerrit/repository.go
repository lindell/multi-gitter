package gerrit

import "fmt"

type repository struct {
	url           string
	webURL        string
	name          string
	defaultBranch string
}

func (r repository) CloneURL() string {
	return r.url
}

func (r repository) BranchURL(_ string) string {
	if r.webURL == "" {
		return ""
	}

	return fmt.Sprintf("%s,+branches", r.webURL)
}

func (r repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r repository) FullName() string {
	return r.name
}

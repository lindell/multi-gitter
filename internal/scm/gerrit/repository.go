package gerrit

type repository struct {
	url           string
	name          string
	defaultBranch string
}

func (r repository) CloneURL() string {
	return r.url
}

func (r repository) BranchURL(branchName string) string {
	// Not yet implemented
	return ""
}

func (r repository) DefaultBranch() string {
	return r.defaultBranch
}

func (r repository) FullName() string {
	return r.name
}

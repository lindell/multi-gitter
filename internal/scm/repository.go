package scm

// Repository provides all the information needed about a git repository
type Repository interface {
	// CloneURL returns the clone address of the repository
	CloneURL() string
	// DefaultBranch returns the name of the default branch of the repository
	DefaultBranch() string
	// FullName returns the full id of the repository, usually ownerName/repoName
	FullName() string
}

func RepoContainsTopic(slice1 []string, slice2 []string) bool {
	slice1Map := map[string]struct{}{}
	for _, v := range slice1 {
		slice1Map[v] = struct{}{}
	}

	for _, v := range slice2 {
		if _, ok := slice1Map[v]; ok {
			return true
		}
	}

	return false
}

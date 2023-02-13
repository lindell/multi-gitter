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

func RepoContainsTopic(repoTopics []string, filterTopics []string) bool {
	repoTopicsMap := map[string]struct{}{}
	for _, v := range repoTopics {
		repoTopicsMap[v] = struct{}{}
	}

	for _, v := range filterTopics {
		if _, ok := repoTopicsMap[v]; ok {
			return true
		}
	}

	return false
}

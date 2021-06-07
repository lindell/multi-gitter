package domain

// Repository contains all information about a git repository
type Repository interface {
	URL(token string) string
	DefaultBranch() string
	// Returns the full id of the repository, usually ownerName/repoName
	FullName() string
}

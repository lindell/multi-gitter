package git

// Repository provides all the information needed about a git repository
type Repository interface {
	// CloneURL returns the clone address of the repository
	CloneURL() string
	// DefaultBranch returns the name of the default branch of the repository
	DefaultBranch() string
	// FullName returns the full id of the repository, usually ownerName/repoName
	FullName() string
}

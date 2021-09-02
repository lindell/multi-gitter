package repository

// Data contains all information about a git repository
type Data interface {
	// URL returns the clone address of the repository
	URL(token string) string
	// URLWithUsername returns the clone address of the repository for SCMs that only support username and token
	URLWithUsername(username, token string) string
	// DefaultBranch returns the name of the default branch of the repository
	DefaultBranch() string
	// FullName returns the full id of the repository, usually ownerName/repoName
	FullName() string
}

package git

// Changes represents the changes made to a repository
type Changes struct {
	// Map of file paths to the changes made to the file
	// The key is the file path and the value is the change
	Additions map[string][]byte

	// List of file paths that were deleted
	Deletions []string

	// OldHash is the hash of the previous commit
	OldHash string
}

type LastCommitChecker interface {
	LastCommitChanges() (Changes, error)
}

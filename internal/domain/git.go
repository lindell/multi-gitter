package domain

// GitConfig is configration for any git implementation
type GitConfig struct {
	// Absolute path to the directory
	Directory string
	// The fetch depth used when cloning, if set to 0, the entire history will be used
	FetchDepth int
}

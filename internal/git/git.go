package git

// Config is configuration for any git implementation
type Config struct {
	// Absolute path to the directory
	Directory string
	// The fetch depth used when cloning, if set to 0, the entire history will be used
	FetchDepth int
}

package domain

// GitConfig is configration for any git implementation
type GitConfig struct {
	// Absolute path to the directory
	Directory string

	FetchDepth int
}

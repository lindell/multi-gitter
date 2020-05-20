package domain

// Repository is an interface with generic methods for any repository
type Repository interface {
	GetURL() string
	GetBranch() string
}

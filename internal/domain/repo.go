package domain

type Repository interface {
	GetURL() string
	GetBranch() string
}

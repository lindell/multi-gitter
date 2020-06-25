package domain

import "fmt"

// Repository contains all information about a git repository
type Repository struct {
	URL           string
	Name          string
	OwnerName     string
	DefaultBranch string
}

func (r Repository) FullName() string {
	return fmt.Sprintf("%s/%s", r.OwnerName, r.Name)
}

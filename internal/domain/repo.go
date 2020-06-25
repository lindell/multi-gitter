package domain

import "fmt"

// Repository contains all information about a git repository
type Repository struct {
	URL           string
	Name          string
	OwnerName     string
	DefaultBranch string
}

// FullName returns the full repository name including the owner
func (r Repository) FullName() string {
	return fmt.Sprintf("%s/%s", r.OwnerName, r.Name)
}

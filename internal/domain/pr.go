package domain

// NewPullRequest is the data needed to create a new pull request
type NewPullRequest struct {
	Title string
	Body  string
	Head  string
	Base  string

	Reviewers []string // The username of all reviewers
}

package domain

// NewPullRequest is the data needed to create a new pull request
type NewPullRequest struct {
	Title string
	Body  string
	Head  string
	Base  string

	Reviewers []string // The username of all reviewers
}

// PullRequestStatus is the status of a pull request, including statuses of the last commit
type PullRequestStatus int

// All PullRequestStatuses
const (
	PullRequestStatusUnknown PullRequestStatus = iota
	PullRequestStatusSuccess
	PullRequestStatusPending
	PullRequestStatusError
	PullRequestStatusMerged
	PullRequestStatusClosed
)

func (s PullRequestStatus) String() string {
	switch s {
	case PullRequestStatusUnknown:
		return "Unknown"
	case PullRequestStatusSuccess:
		return "Success"
	case PullRequestStatusPending:
		return "Pending"
	case PullRequestStatusError:
		return "Error"
	case PullRequestStatusMerged:
		return "Merged"
	case PullRequestStatusClosed:
		return "Closed"
	}
	return "Unknown"
}

// PullRequest represents a pull request
type PullRequest struct {
	OwnerName  string
	RepoName   string
	BranchName string
	Number     int
	Status     PullRequestStatus
}

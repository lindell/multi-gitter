package git

import (
	"fmt"
	"strings"
)

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
type PullRequest interface {
	Status() PullRequestStatus
	String() string
}

// MergeType is the way a pull request is "merged" into the base branch
type MergeType int

// All MergeTypes
const (
	MergeTypeUnknown MergeType = iota
	MergeTypeMerge
	MergeTypeRebase
	MergeTypeSquash
)

// ParseMergeType parses a merge type
func ParseMergeType(typ string) (MergeType, error) {
	switch strings.ToLower(typ) {
	case "merge":
		return MergeTypeMerge, nil
	case "rebase":
		return MergeTypeRebase, nil
	case "squash":
		return MergeTypeSquash, nil
	}
	return MergeTypeUnknown, fmt.Errorf(`not a valid merge type: "%s"`, typ)
}

// MergeTypeIntersection calculates the intersection of two merge type slices,
// The order of the first slice will be preserved
func MergeTypeIntersection(mergeTypes1, mergeTypes2 []MergeType) []MergeType {
	res := []MergeType{}
	for _, mt := range mergeTypes1 {
		for _, mt2 := range mergeTypes2 {
			if mt == mt2 {
				res = append(res, mt)
			}
		}
	}
	return res
}

package pullrequest

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

// Status is the status of a pull request, including statuses of the last commit
type Status int

// All PullRequestStatuses
const (
	StatusUnknown Status = iota
	StatusSuccess
	StatusPending
	StatusError
	StatusMerged
	StatusClosed
)

func (s Status) String() string {
	switch s {
	case StatusUnknown:
		return "Unknown"
	case StatusSuccess:
		return "Success"
	case StatusPending:
		return "Pending"
	case StatusError:
		return "Error"
	case StatusMerged:
		return "Merged"
	case StatusClosed:
		return "Closed"
	}
	return "Unknown"
}

// PullRequest represents a pull request
type PullRequest interface {
	Status() Status
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

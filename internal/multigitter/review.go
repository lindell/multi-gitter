package multigitter

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Reviewer reviews pull requests in an organization
type Reviewer struct {
	VersionController VersionController
	FeatureBranch     string
	Comment           string
	All               bool
	Batch             string
	Pager             string
	DisablePaging     bool
	IncludeApproved   bool
}

const (
	header        = "# *********************************************************"
	reviewOptions = "[a/r/c/n/d/h (help)]"
	batchApprove  = "approve"
	batchReject   = "reject"
	batchComment  = "comment"
)

type approvedPR struct {
	scm.PullRequest
	approved bool
}

// Reviewer reviews pull requests in an organization
func (s Reviewer) Review(ctx context.Context) error {
	prs, err := s.VersionController.GetPullRequests(ctx, s.FeatureBranch)
	if err != nil {
		return err
	}

	var approvedPrs []approvdPr

	for _, pr := range prs {
		log := log.WithField("pr", pr.String())
		approved, err := s.VersionController.IsPullRequestApprovedByMe(ctx, pr)
		if err != nil {
			log.Errorf("Failed to retrieve pull request reviews: %s", err.Error())
			continue
		}

		if approved && !s.IncludeApproved {
			continue
		}

		approvedPrs = append(approvedPrs, approvdPr{
			PullRequest: pr,
			approved:    approved,
		})
	}

	if len(approvedPrs) == 0 {
		fmt.Println("All pull requests are approved by you. You can still view and comment them by using review --include-approved")
		return nil
	}

	var reviewDiffs string

	for _, pr := range approvedPrs {
		log := log.WithField("pr", pr.String())
		
		diff, err := s.VersionController.DiffPullRequest(ctx, pr.PullRequest)
		if err != nil {
			log.Errorf("Error occurred while retrieving diff: %s", err.Error())
			continue
		}

		if !s.All && s.Batch == "" {
			if err := s.leaveReview(ctx, strings.NewReader(fmt.Sprintf("%s\n# %s\n%s\n%s", header, pr.String(), header, diff)), pr); err != nil {
				log.Errorf("Error occurred while reviewing pr: %s", err.Error())
			}
		} else {
			reviewDiffs = fmt.Sprintf("%s\n%s\n# %s\n%s\n%s", reviewDiffs, header, pr.String(), header, diff)
		}
	}

	if s.All || s.Batch != "" {
		return s.leaveReviews(ctx, strings.NewReader(reviewDiffs), approvedPrs...)
	}

	return nil
}

func (s Reviewer) printDiff(r io.ReadSeeker) error {
	var err error

	// attempt to page the diff
	if !s.DisablePaging {
		_, _ = r.Seek(0, 0)
		err = s.pageTmpFile(r)
	}

	// if paging failed (or is diabled) the diff gets dumped to stdout instead
	if err != nil || s.DisablePaging {
		_, _ = r.Seek(0, 0)
		b, err := io.ReadAll(r)
		if err != nil {
			return err
		}

		fmt.Println(string(b))
	}

	return nil
}

// write file to tmp file so a pager can also be a tool which can't directly read from stdin
func (s Reviewer) pageTmpFile(r io.Reader) error {
	file, err := os.CreateTemp(os.TempDir(), "*.diff")
	if err != nil {
		return err
	}

	defer os.Remove(file.Name())

	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}

	cmd := exec.Command(s.Pager, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (s Reviewer) leaveReview(ctx context.Context, r io.ReadSeeker, pr approvdPr) error {
	if err := s.printDiff(r); err != nil {
		return err
	}

	var action string
	for {
		fmt.Printf("Leave a review on %s %s: ", pr.String(), reviewOptions)
		fmt.Scanln(&action)
		if repeat, _ := s.reviewAction(ctx, r, action, pr); !repeat {
			break
		}
	}

	return nil
}

func (s Reviewer) leaveReviews(ctx context.Context, r io.ReadSeeker, prs ...approvdPr) error {
	if err := s.printDiff(r); err != nil {
		return err
	}

	if s.Batch != "" {
		switch s.Batch {
		case batchApprove:
			_, err := s.reviewAction(ctx, r, "a", prs...)
			return err
		case batchReject:
			_, err := s.reviewAction(ctx, r, "r", prs...)
			return err
		case batchComment:
			_, err := s.reviewAction(ctx, r, "c", prs...)
			return err
		default:
			return errors.New("unknown batch operation")
		}
	}

	var action string
	for {
		fmt.Printf("Leave a review on all pull requests %s: ", reviewOptions)
		fmt.Scanln(&action)

		if repeat, _ := s.reviewAction(ctx, r, action, prs...); !repeat {
			break
		}
	}

	return nil
}

func (s Reviewer) reviewAction(ctx context.Context, r io.ReadSeeker, action string, prs ...approvdPr) (bool, error) {
	switch action {
	case "a", "r", "c":
		comment := s.getComment()
		for _, pr := range prs {
			log.Infof("Reviewing %s pull request", pr.String())

			switch action {
			case "a":
				err := s.approve(ctx, pr, comment)
				if err != nil {
					log.Errorf("Error occurred while approving: %s", err.Error())
				}

			case "r":
				err := s.VersionController.RejectPullRequest(ctx, pr.PullRequest, comment)
				if err != nil {
					log.Errorf("Error occurred while rejecting: %s", err.Error())
				}
			case "c":
				err := s.VersionController.CommentPullRequest(ctx, pr.PullRequest, comment)
				if err != nil {
					log.Errorf("Error occurred while commenting: %s", err.Error())
				}
			}
		}

		return false, nil
	case "n":
		return false, nil
	case "d":
		return true, s.printDiff(r)
	default:
		fmt.Println("a - Approve pull request")
		fmt.Println("r - Request changes")
		fmt.Println("c - Only leave a comment")
		fmt.Println("n - Skip review")
		fmt.Println("d - Show the pull request diff again")
		fmt.Println("h - Display these optons")
		return true, nil
	}
}

// approve gates double approvement submits
func (s Reviewer) approve(ctx context.Context, pr approvdPr, comment string) error {
	if pr.approved {
		log.Debug("Skip approved pull request", pr.String())
		return nil
	}

	return s.VersionController.ApprovePullRequest(ctx, pr.PullRequest, comment)
}

func (s Reviewer) getComment() string {
	if s.Batch != "" {
		return s.Comment
	}

	comment := s.Comment
	fmt.Printf("Leave a comment [%s]: ", comment)
	fmt.Scanln(&comment)
	return comment
}

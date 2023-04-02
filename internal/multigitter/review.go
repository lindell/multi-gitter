package multigitter

import (
	"context"
	"fmt"
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
	BatchOperation    *BatchOperation
	Pager             string
	DisablePaging     bool
	IncludeApproved   bool
}

const (
	header        = "# *********************************************************"
	reviewOptions = "[a/r/c/n/d/h (help)]"
)

type reviewPR struct {
	scm.PullRequest
	approved bool
	diff     string
}

func (pr reviewPR) formatReviewDiff() string {
	return fmt.Sprintf("%s\n# %s\n%s\n%s", header, pr.String(), header, pr.diff)
}

// Reviewer reviews pull requests in an organization
func (s Reviewer) Review(ctx context.Context) error {
	if s.BatchOperation != nil {
		return s.reviewAllPrs(ctx, *s.BatchOperation)
	} else if s.All {
		return s.reviewAllPrsWithQuestion(ctx)
	} else {
		return s.reviewAllPrsIndividually(ctx)
	}
}

func (s Reviewer) reviewAllPrsIndividually(ctx context.Context) error {
	prs, err := s.getReviewablePRs(ctx)
	if err != nil {
		return err
	}

	for _, pr := range prs {
		if err := s.leaveReview(ctx, pr.diff, pr); err != nil {
			log.Errorf("Error occurred while reviewing pr: %s", err.Error())
		}
	}

	return nil
}

func (s Reviewer) reviewAllPrs(ctx context.Context, batchOperation BatchOperation) error {
	prs, err := s.getReviewablePRs(ctx)
	if err != nil {
		return err
	}

	diff := prDiffs(prs)

	switch batchOperation {
	case BatchOperationApprove:
		_, err := s.reviewAction(ctx, diff, "a", prs...)
		return err
	case BatchOperationReject:
		_, err := s.reviewAction(ctx, diff, "r", prs...)
		return err
	case BatchOperationComment:
		_, err := s.reviewAction(ctx, diff, "c", prs...)
		return err
	default:
		return errors.New("unknown batch operation")
	}
}

func (s Reviewer) reviewAllPrsWithQuestion(ctx context.Context) error {
	prs, err := s.getReviewablePRs(ctx)
	if err != nil {
		return err
	}

	diff := prDiffs(prs)
	if err := s.printDiff(diff); err != nil {
		return err
	}

	var action string
	for {
		fmt.Printf("Leave a review on all pull requests %s: ", reviewOptions)
		fmt.Scanln(&action)

		if repeat, _ := s.reviewAction(ctx, diff, action, prs...); !repeat {
			break
		}
	}

	return nil
}

func prDiffs(prs []reviewPR) string {
	var reviewDiffs strings.Builder
	for _, pr := range prs {
		reviewDiffs.WriteString(pr.formatReviewDiff())
		reviewDiffs.WriteString("\n")
	}
	return reviewDiffs.String()
}

func (s Reviewer) getReviewablePRs(ctx context.Context) ([]reviewPR, error) {
	prs, err := s.VersionController.GetPullRequests(ctx, s.FeatureBranch)
	if err != nil {
		return nil, err
	}

	if len(prs) == 0 {
		return nil, fmt.Errorf("no open pull requests found matching the branch %s", s.FeatureBranch)
	}

	var reviewPrs []reviewPR

	for _, pr := range prs {
		log := log.WithField("pr", pr.String())
		log.Debug("Retrieving pull request reviews")

		approved, err := s.VersionController.IsPullRequestApprovedByMe(ctx, pr)
		if err != nil {
			log.Errorf("Failed to retrieve pull request reviews: %s", err.Error())
			continue
		}

		if approved && !s.IncludeApproved {
			continue
		}

		log.Debug("Retrieving pull request diff")

		diff, err := s.VersionController.DiffPullRequest(ctx, pr)
		if err != nil {
			log.Errorf("Error occurred while retrieving diff: %s", err.Error())
			continue
		}

		reviewPrs = append(reviewPrs, reviewPR{
			PullRequest: pr,
			approved:    approved,
			diff:        diff,
		})
	}

	if len(reviewPrs) == 0 {
		return nil, fmt.Errorf("all pull requests are approved by you. You can still view and comment them by using review --include-approved")
	}

	return reviewPrs, nil
}

func (s Reviewer) printDiff(diff string) error {
	var err error

	// attempt to page the diff
	if !s.DisablePaging {
		err = s.pageTmpFile(diff)
	}

	// if paging failed (or is diabled) the diff gets dumped to stdout instead
	if err != nil || s.DisablePaging {
		fmt.Println(diff)
	}

	return nil
}

// write file to tmp file so a pager can also be a tool which can't directly read from stdin
func (s Reviewer) pageTmpFile(diff string) error {
	file, err := os.CreateTemp(os.TempDir(), "*.diff")
	if err != nil {
		return err
	}

	defer os.Remove(file.Name())

	_, err = file.Write([]byte(diff))
	if err != nil {
		return err
	}

	cmd := exec.Command(s.Pager, file.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (s Reviewer) leaveReview(ctx context.Context, diff string, pr reviewPR) error {
	if err := s.printDiff(diff); err != nil {
		return err
	}

	var action string
	for {
		fmt.Printf("Leave a review on %s %s: ", pr.String(), reviewOptions)
		fmt.Scanln(&action)
		if repeat, _ := s.reviewAction(ctx, diff, action, pr); !repeat {
			break
		}
	}

	return nil
}

func (s Reviewer) reviewAction(ctx context.Context, diff string, action string, prs ...reviewPR) (bool, error) {
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
		return true, s.printDiff(diff)
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
func (s Reviewer) approve(ctx context.Context, pr reviewPR, comment string) error {
	if pr.approved {
		log.Debug("Skip approved pull request", pr.String())
		return nil
	}

	return s.VersionController.ApprovePullRequest(ctx, pr.PullRequest, comment)
}

func (s Reviewer) getComment() string {
	if s.BatchOperation != nil {
		return s.Comment
	}

	comment := s.Comment
	fmt.Printf("Leave a comment [%s]: ", comment)
	fmt.Scanln(&comment)
	return comment
}

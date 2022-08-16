package multigitter

import (
	"context"

	"github.com/lindell/multi-gitter/internal/scm"
	log "github.com/sirupsen/logrus"
)

// Merger merges pull requests in an organization
type Merger struct {
	VersionController VersionController

	FeatureBranch string
}

// Merge merges pull requests in an organization
func (s Merger) Merge(ctx context.Context) error {
	prs, err := s.VersionController.GetPullRequests(ctx, s.FeatureBranch)
	if err != nil {
		return err
	}

	successPrs := make([]scm.PullRequest, 0, len(prs))
	for _, pr := range prs {
		if pr.Status() == scm.PullRequestStatusSuccess && pr.Status() != scm.PullRequestStatusMerged {
			successPrs = append(successPrs, pr)
		}
	}

	log.Infof("Merging %d pull requests", len(successPrs))

	for _, pr := range successPrs {
		log.WithField("pr", pr.String()).Infof("Merging")
		err := s.VersionController.MergePullRequest(ctx, pr)
		if err != nil {
			return err
		}
	}

	return nil
}

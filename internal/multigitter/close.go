package multigitter

import (
	"context"

	"github.com/lindell/multi-gitter/internal/domain"
	log "github.com/sirupsen/logrus"
)

// Closer closes pull requests
type Closer struct {
	VersionController VersionController

	FeatureBranch string
}

// Close closes pull requests
func (s Closer) Close(ctx context.Context) error {
	prs, err := s.VersionController.GetPullRequests(ctx, s.FeatureBranch)
	if err != nil {
		return err
	}

	openPRs := make([]domain.PullRequest, 0, len(prs))
	for _, pr := range prs {
		if pr.Status() != domain.PullRequestStatusClosed && pr.Status() != domain.PullRequestStatusMerged {
			openPRs = append(openPRs, pr)
		}
	}

	log.Infof("Closing %d pull requests", len(openPRs))

	for _, pr := range openPRs {
		log.WithField("pr", pr.String()).Infof("Closing")
		err := s.VersionController.ClosePullRequest(ctx, pr)
		if err != nil {
			return err
		}
	}

	return nil
}

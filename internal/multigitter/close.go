package multigitter

import (
	"context"

	"github.com/lindell/multi-gitter/internal/scm"
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

	openPRs := make([]scm.PullRequest, 0, len(prs))
	for _, pr := range prs {
		if pr.Status() != scm.PullRequestStatusClosed && pr.Status() != scm.PullRequestStatusMerged {
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

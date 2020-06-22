package multigitter

import (
	"context"

	"github.com/lindell/multi-gitter/internal/domain"
)

// Merger merges pull requests in an organization
type Merger struct {
	VersionController VersionController

	FeatureBranch string
	OrgName       string
}

// Merge merges pull requests in an organization
func (s Merger) Merge(ctx context.Context) error {
	prs, err := s.VersionController.GetPullRequestStatuses(ctx, s.OrgName, s.FeatureBranch)
	if err != nil {
		return err
	}

	successPrs := make([]domain.PullRequest, 0, len(prs))
	for _, pr := range prs {
		if pr.Status == domain.PullRequestStatusSuccess {
			successPrs = append(successPrs, pr)
		}
	}

	for _, pr := range successPrs {
		err := s.VersionController.MergePullRequest(ctx, pr)
		if err != nil {
			return err
		}
	}

	return nil
}

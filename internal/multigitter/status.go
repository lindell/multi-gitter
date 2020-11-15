package multigitter

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Statuser checks the statuses of pull requests
type Statuser struct {
	VersionController VersionController

	FeatureBranch string
}

// Statuses checks the statuses of pull requests
func (s Statuser) Statuses(ctx context.Context) error {
	prs, err := s.VersionController.GetPullRequestStatuses(ctx, s.FeatureBranch)
	if err != nil {
		return err
	}

	for _, pr := range prs {
		fmt.Fprintf(log.StandardLogger().Out, "%s: %s\n", pr.String(), pr.Status())
	}

	return nil
}

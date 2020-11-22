package multigitter

import (
	"context"
	"fmt"
	"io"
)

// Statuser checks the statuses of pull requests
type Statuser struct {
	VersionController VersionController

	Output io.Writer

	FeatureBranch string
}

// Statuses checks the statuses of pull requests
func (s Statuser) Statuses(ctx context.Context) error {
	prs, err := s.VersionController.GetPullRequestStatuses(ctx, s.FeatureBranch)
	if err != nil {
		return err
	}

	for _, pr := range prs {
		fmt.Fprintf(s.Output, "%s: %s\n", pr.String(), pr.Status())
	}

	return nil
}

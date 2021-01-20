package multigitter

import (
	"context"
	"fmt"
	"io"

	"github.com/lindell/multi-gitter/internal/multigitter/terminal"
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
		if urler, ok := pr.(urler); ok {
			fmt.Fprintf(s.Output, "%s: %s\n", terminal.Link(pr.String(), urler.URL()), pr.Status())
		} else {
			fmt.Fprintf(s.Output, "%s: %s\n", pr.String(), pr.Status())
		}
	}

	return nil
}

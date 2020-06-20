package multigitter

import (
	"context"
	"fmt"
)

type Statuser struct {
	VersionController VersionController

	FeatureBranch string
	OrgName       string
}

func (s Statuser) Statuses(ctx context.Context) error {
	statuses, err := s.VersionController.GetPullRequestStatuses(ctx, s.OrgName, s.FeatureBranch)
	if err != nil {
		return err
	}

	for _, status := range statuses {
		fmt.Printf("%s: %s\n", status.RepoName, status.Status)
	}

	return nil
}

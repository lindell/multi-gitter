package github

import (
	"testing"

	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/stretchr/testify/assert"
)

func Test_convertGraphQLPullRequest(t *testing.T) {
	scenarios := []struct {
		name     string
		pr       graphqlPR
		expected pullRequest
	}{{
		name: "should return status 'closed' when the PR branch is deleted",
		pr: graphqlPR{
			Number:      1,
			HeadRefName: "dummy_branch",
			Closed:      true,
			URL:         "http://dummy.url",
			Merged:      false,
			BaseRepository: struct {
				Name  string "json:\"name\""
				Owner struct {
					Login string "json:\"login\""
				} "json:\"owner\""
			}{
				Name: "base_repo",
				Owner: struct {
					Login string "json:\"login\""
				}{Login: "dummy_user"},
			},
			HeadRepository: struct {
				Name  string "json:\"name\""
				Owner struct {
					Login string "json:\"login\""
				} "json:\"owner\""
			}{
				Name: "pr_owner",
				Owner: struct {
					Login string "json:\"login\""
				}{Login: "dummy_owner"},
			},
		},
		expected: pullRequest{
			status:      scm.PullRequestStatusClosed,
			ownerName:   "dummy_user",
			repoName:    "base_repo",
			branchName:  "dummy_branch",
			prOwnerName: "dummy_owner",
			prRepoName:  "pr_owner",
			number:      1,
			guiURL:      "http://dummy.url",
		},
	}}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			got := convertGraphQLPullRequest(scenario.pr)
			if got != scenario.expected {
				assert.Equal(t, scenario.expected, got)
			}
		})
	}
}

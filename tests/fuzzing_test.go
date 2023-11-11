package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/lindell/multi-gitter/cmd"
	"github.com/lindell/multi-gitter/tests/vcmock"
	"github.com/stretchr/testify/assert"
)

func FuzzRun(f *testing.F) {
	f.Add(
		"assignee1,assignee2",           // assignees
		"commit message",                // commit-message
		1,                               // concurrent
		"skip",                          // conflict-strategy
		false,                           // draft
		false,                           // dry-run
		1,                               // fetch-depth
		false,                           // fork
		"fork-owner",                    // fork-owner
		"go",                            // git-type
		"label1,label2",                 // labels
		"text",                          // log-format
		"info",                          // log-level
		1,                               // max-reviewers
		1,                               // max-team-reviewers
		"pr-body",                       // pr-body
		"pr-title",                      // pr-title
		"reviewer1,reviewer2",           // reviewers
		false,                           // skip-forks
		false,                           // skip-pr
		"should-not-change",             // skip-repo
		"team-reviewer1,team-reviewer1", // team-reviewers
		"topic1,topic2",                 // topic
	)
	f.Fuzz(func(
		t *testing.T,

		assignees string,
		commitMessage string,
		concurrent int,
		conflictStrategy string,
		draft bool,
		dryRun bool,
		fetchDepth int,
		fork bool,
		forkOwner string,
		gitType string,
		labels string,
		logFormat string,
		logLevel string,
		maxReviewers int,
		maxTeamReviewers int,
		prBody string,
		prTitle string,
		reviewers string,
		skipForks bool,
		skipPr bool,
		skipRepo string,
		teamReviewers string,
		topic string,
	) {
		vcMock := &vcmock.VersionController{}
		defer vcMock.Clean()
		cmd.OverrideVersionController = vcMock

		tmpDir, err := os.MkdirTemp(os.TempDir(), "multi-git-test-run-")
		defer os.RemoveAll(tmpDir)
		assert.NoError(t, err)

		workingDir, err := os.Getwd()
		assert.NoError(t, err)

		changerBinaryPath := normalizePath(filepath.Join(workingDir, changerBinaryPath))

		changeRepo := createRepo(t, "owner", "should-change", "i like apples")
		changeRepo2 := createRepo(t, "owner", "should-change-2", "i like my apple")
		noChangeRepo := createRepo(t, "owner", "should-not-change", "i like oranges")
		vcMock.AddRepository(changeRepo)
		vcMock.AddRepository(changeRepo2)
		vcMock.AddRepository(noChangeRepo)

		runOutFile := filepath.Join(tmpDir, "run-out.txt")
		runLogFile := filepath.Join(tmpDir, "run-log.txt")

		command := cmd.RootCmd()
		command.SetArgs([]string{"run",
			"--output", runOutFile,
			"--log-file", runLogFile,
			"--author-name", "Test Author",
			"--author-email", "test@example.com",
			"--assignees", assignees,
			"--commit-message", commitMessage,
			"--concurrent", fmt.Sprint(concurrent),
			"--conflict-strategy", conflictStrategy,
			fmt.Sprintf("--draft=%t", draft),
			fmt.Sprintf("--dry-run=%t", dryRun),
			"--fetch-depth", fmt.Sprint(fetchDepth),
			fmt.Sprintf("--fork=%t", fork),
			"--fork-owner", forkOwner,
			"--git-type", gitType,
			"--labels", labels,
			"--log-format", logFormat,
			"--log-level", logLevel,
			"--max-reviewers", fmt.Sprint(maxReviewers),
			"--max-team-reviewers", fmt.Sprint(maxTeamReviewers),
			"--pr-body", prBody,
			"--pr-title", prTitle,
			"--reviewers", reviewers,
			fmt.Sprintf("--skip-forks=%t", skipForks),
			fmt.Sprintf("--skip-pr=%t", skipPr),
			"--skip-repo", skipRepo,
			"--team-reviewers", teamReviewers,
			"--topic", topic,
			changerBinaryPath,
		})
		err = command.Execute()
		if err != nil {
			assert.NotContains(t, err.Error(), "panic")
		}

		// Verify that the output was correct
		runOutData, _ := os.ReadFile(runOutFile)
		assert.NotContains(t, string(runOutData), "panic")
		runLogData, _ := os.ReadFile(runLogFile)
		assert.NotContains(t, string(runLogData), "panic")
	})
}

package tests

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/lindell/multi-gitter/cmd"
	"github.com/lindell/multi-gitter/internal/git"
	"github.com/lindell/multi-gitter/tests/vcmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStory tests the common usecase: run, status, merge, status
func TestStory(t *testing.T) {
	vcMock := &vcmock.VersionController{}
	defer vcMock.Clean()
	cmd.OverrideVersionController = vcMock

	tmpDir, err := ioutil.TempDir(os.TempDir(), "multi-git-test-run-")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)

	workingDir, err := os.Getwd()
	assert.NoError(t, err)

	changerBinaryPath := filepath.ToSlash(filepath.Join(workingDir, changerBinaryPath))

	changeRepo := createRepo(t, "owner", "should-change", "i like apples")
	changeRepo2 := createRepo(t, "owner", "should-change-2", "i like my apple")
	noChangeRepo := createRepo(t, "owner", "should-not-change", "i like oranges")
	vcMock.AddRepository(changeRepo)
	vcMock.AddRepository(changeRepo2)
	vcMock.AddRepository(noChangeRepo)

	runOutFile := filepath.Join(tmpDir, "run-log.txt")

	command := cmd.RootCmd()
	command.SetArgs([]string{"run",
		"--output", runOutFile,
		"--author-name", "Test Author",
		"--author-email", "test@example.com",
		"-B", "custom-branch-name",
		"-m", "test",
		changerBinaryPath,
	})
	err = command.Execute()
	assert.NoError(t, err)

	// Verify that the data of the original branch is intact
	data, err := ioutil.ReadFile(filepath.Join(changeRepo.Path, fileName))
	assert.NoError(t, err)
	assert.Equal(t, []byte("i like apples"), data)

	// Verify that the new branch is changed
	changeBranch(t, changeRepo.Path, "custom-branch-name", false)
	data, err = ioutil.ReadFile(filepath.Join(changeRepo.Path, fileName))
	assert.NoError(t, err)
	assert.Equal(t, []byte("i like bananas"), data)

	// Verify that the output was correct
	runOutData, err := ioutil.ReadFile(runOutFile)
	require.NoError(t, err)
	assert.Equal(t, `No data was changed:
  owner/should-not-change
Repositories with a successful run:
  owner/should-change #1
  owner/should-change-2 #2
`, string(runOutData))

	//
	// PullRequestStatus
	//
	statusOutFile := filepath.Join(tmpDir, "status-log.txt")

	command = cmd.RootCmd()
	command.SetArgs([]string{"status",
		"--output", statusOutFile,
		"-B", "custom-branch-name",
	})
	err = command.Execute()
	assert.NoError(t, err)

	// Verify that the output was correct
	statusOutData, err := ioutil.ReadFile(statusOutFile)
	require.NoError(t, err)
	assert.Equal(t, "owner/should-change #1: Pending\nowner/should-change-2 #2: Pending\n", string(statusOutData))

	// One of the created PRs is set to succeeded
	vcMock.SetPRStatus("should-change", "custom-branch-name", git.PullRequestStatusSuccess)

	//
	// Merge
	//
	mergeLogFile := filepath.Join(tmpDir, "merge-log.txt")

	command = cmd.RootCmd()
	command.SetArgs([]string{"merge",
		"--log-file", mergeLogFile,
		"-B", "custom-branch-name",
	})
	err = command.Execute()
	assert.NoError(t, err)

	// Verify that the output was correct
	mergeLogData, err := ioutil.ReadFile(mergeLogFile)
	require.NoError(t, err)
	assert.Contains(t, string(mergeLogData), "Merging 1 pull requests")
	assert.Contains(t, string(mergeLogData), "Merging pr=\"owner/should-change #1\"")

	//
	// After Merge PullRequestStatus
	//
	afterMergeStatusOutFile := filepath.Join(tmpDir, "after-merge-status-log.txt")

	command = cmd.RootCmd()
	command.SetArgs([]string{"status",
		"--output", afterMergeStatusOutFile,
		"-B", "custom-branch-name",
	})
	err = command.Execute()
	assert.NoError(t, err)

	// Verify that the output was correct
	afterMergeStatusOutData, err := ioutil.ReadFile(afterMergeStatusOutFile)
	require.NoError(t, err)
	assert.Equal(t, "owner/should-change #1: Merged\nowner/should-change-2 #2: Pending\n", string(afterMergeStatusOutData))

	//
	// Close
	//
	closeLogFile := filepath.Join(tmpDir, "close-log.txt")

	command = cmd.RootCmd()
	command.SetArgs([]string{"close",
		"--log-file", closeLogFile,
		"-B", "custom-branch-name",
	})
	err = command.Execute()
	assert.NoError(t, err)

	// Verify that the output was correct
	closeLogData, err := ioutil.ReadFile(closeLogFile)
	require.NoError(t, err)
	assert.Contains(t, string(closeLogData), "Closing 1 pull request")
	assert.Contains(t, string(closeLogData), "Closing pr=\"owner/should-change-2 #2\"")

	//
	// After Close PullRequestStatus
	//
	afterCloseStatusOutFile := filepath.Join(tmpDir, "after-close-status-log.txt")

	command = cmd.RootCmd()
	command.SetArgs([]string{"status",
		"--output", afterCloseStatusOutFile,
		"-B", "custom-branch-name",
	})
	err = command.Execute()
	assert.NoError(t, err)

	// Verify that the output was correct
	afterCloseStatusOutData, err := ioutil.ReadFile(afterCloseStatusOutFile)
	require.NoError(t, err)
	assert.Equal(t, "owner/should-change #1: Merged\nowner/should-change-2 #2: Closed\n", string(afterCloseStatusOutData))
}

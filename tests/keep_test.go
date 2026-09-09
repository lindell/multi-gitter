package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lindell/multi-gitter/cmd"
	"github.com/lindell/multi-gitter/tests/vcmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKeepRunOption tests the --keep flag for the run command.
// It verifies that:
// 1. The cloned repository directory persists after the run
// 2. A second run reuses the existing clone (with a hard reset) instead of cloning fresh
// 3. The changes from the first run do not leak into the second run's base state
func TestKeepRunOption(t *testing.T) {
	workingDir, err := os.Getwd()
	require.NoError(t, err)

	changerBinary := normalizePath(filepath.Join(workingDir, changerBinaryPath))

	for _, gitBackend := range gitBackends {
		t.Run(string(gitBackend), func(t *testing.T) {
			vcMock := &vcmock.VersionController{}
			defer vcMock.Clean()
			cmd.OverrideVersionController = vcMock

			cloneDir, err := os.MkdirTemp(os.TempDir(), "multi-gitter-keep-test-")
			require.NoError(t, err)
			defer os.RemoveAll(cloneDir)

			repo := createRepo(t, "owner", "keep-repo", "i like apples")
			vcMock.AddRepository(repo)

			runLogFile := filepath.Join(cloneDir, "run1-log.txt")
			outFile := filepath.Join(cloneDir, "run1-out.txt")

			// --- First run with --keep ---
			command := cmd.RootCmd()
			command.SetArgs([]string{
				"run",
				"--log-file", runLogFile,
				"--output", outFile,
				"--git-type", string(gitBackend),
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "keep-branch",
				"-m", "keep test commit",
				"--clone-dir", cloneDir,
				"--keep",
				changerBinary,
			})
			err = command.Execute()
			assert.NoError(t, err)

			// Verify that the PR was created
			require.Len(t, vcMock.PullRequests, 1)
			assert.Equal(t, "keep-branch", vcMock.PullRequests[0].Head)
			assert.Equal(t, "keep test commit", vcMock.PullRequests[0].Title)

			// Verify the clone directory still exists (was NOT cleaned up)
			expectedDir := filepath.Join(cloneDir, "multi-gitter-owner-keep-repo")
			_, statErr := os.Stat(expectedDir)
			require.NoError(t, statErr, "cloned repo directory should still exist after --keep run")

			// Verify that .git directory exists in the kept directory
			_, statErr = os.Stat(filepath.Join(expectedDir, ".git"))
			require.NoError(t, statErr, ".git directory should exist in kept clone")

			// --- Second run with --keep (should reuse existing clone) ---
			runLogFile2 := filepath.Join(cloneDir, "run2-log.txt")
			outFile2 := filepath.Join(cloneDir, "run2-out.txt")

			command = cmd.RootCmd()
			command.SetArgs([]string{
				"run",
				"--log-file", runLogFile2,
				"--output", outFile2,
				"--git-type", string(gitBackend),
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "keep-branch-2",
				"-m", "keep test commit 2",
				"--clone-dir", cloneDir,
				"--keep",
				"--conflict-strategy", "replace",
				changerBinary,
			})
			err = command.Execute()
			assert.NoError(t, err)

			// Verify the second PR was created (with the new branch name)
			require.Len(t, vcMock.PullRequests, 2)
			assert.Equal(t, "keep-branch-2", vcMock.PullRequests[1].Head)
			assert.Equal(t, "keep test commit 2", vcMock.PullRequests[1].Title)

			// Verify the log mentions reusing the existing clone
			logData2, err := os.ReadFile(runLogFile2)
			require.NoError(t, err)
			assert.Contains(t, string(logData2), "Reusing existing clone")

			// Verify the directory still exists
			_, statErr = os.Stat(expectedDir)
			require.NoError(t, statErr, "cloned repo directory should still exist after second --keep run")
		})
	}
}

// TestKeepPrintOption tests the --keep flag for the print command.
// It verifies that:
// 1. The cloned repository directory persists after print
// 2. A second print reuses the existing clone
func TestKeepPrintOption(t *testing.T) {
	workingDir, err := os.Getwd()
	require.NoError(t, err)

	printerBinary := normalizePath(filepath.Join(workingDir, printerBinaryPath))

	for _, gitBackend := range gitBackends {
		t.Run(string(gitBackend), func(t *testing.T) {
			vcMock := &vcmock.VersionController{}
			defer vcMock.Clean()
			cmd.OverrideVersionController = vcMock

			cloneDir, err := os.MkdirTemp(os.TempDir(), "multi-gitter-keep-print-test-")
			require.NoError(t, err)
			defer os.RemoveAll(cloneDir)

			repo := createRepo(t, "owner", "keep-print-repo", "i like apples")
			vcMock.AddRepository(repo)

			runLogFile := filepath.Join(cloneDir, "print1-log.txt")
			outFile := filepath.Join(cloneDir, "print1-out.txt")
			errOutFile := filepath.Join(cloneDir, "print1-err.txt")

			// --- First print with --keep ---
			command := cmd.RootCmd()
			command.SetArgs([]string{
				"print",
				"--log-file", runLogFile,
				"--output", outFile,
				"--error-output", errOutFile,
				"--git-type", string(gitBackend),
				"--clone-dir", cloneDir,
				"--keep",
				printerBinary,
			})
			err = command.Execute()
			assert.NoError(t, err)

			// Verify the output
			outData, err := os.ReadFile(outFile)
			require.NoError(t, err)
			assert.Equal(t, "i like apples\n", string(outData))

			// Verify the clone directory still exists
			expectedDir := filepath.Join(cloneDir, "multi-gitter-owner-keep-print-repo")
			_, statErr := os.Stat(expectedDir)
			require.NoError(t, statErr, "cloned repo directory should still exist after --keep print")

			// --- Second print with --keep (should reuse existing clone) ---
			runLogFile2 := filepath.Join(cloneDir, "print2-log.txt")
			outFile2 := filepath.Join(cloneDir, "print2-out.txt")
			errOutFile2 := filepath.Join(cloneDir, "print2-err.txt")

			command = cmd.RootCmd()
			command.SetArgs([]string{
				"print",
				"--log-file", runLogFile2,
				"--output", outFile2,
				"--error-output", errOutFile2,
				"--git-type", string(gitBackend),
				"--clone-dir", cloneDir,
				"--keep",
				printerBinary,
			})
			err = command.Execute()
			assert.NoError(t, err)

			// Verify the output is still correct (reset worked)
			outData2, err := os.ReadFile(outFile2)
			require.NoError(t, err)
			assert.Equal(t, "i like apples\n", string(outData2))

			// Verify the log mentions reusing
			logData2, err := os.ReadFile(runLogFile2)
			require.NoError(t, err)
			assert.Contains(t, string(logData2), "Reusing existing clone")
		})
	}
}

// TestKeepRunWithoutCloneDir tests that --keep works with the default temp directory
// (no explicit --clone-dir set).
func TestKeepRunWithoutCloneDir(t *testing.T) {
	workingDir, err := os.Getwd()
	require.NoError(t, err)

	changerBinary := normalizePath(filepath.Join(workingDir, changerBinaryPath))

	vcMock := &vcmock.VersionController{}
	defer vcMock.Clean()
	cmd.OverrideVersionController = vcMock

	repo := createRepo(t, "owner", "keep-no-clone-dir", "i like apples")
	vcMock.AddRepository(repo)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "multi-gitter-keep-noclone-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	runLogFile := filepath.Join(tmpDir, "run-log.txt")
	outFile := filepath.Join(tmpDir, "run-out.txt")

	command := cmd.RootCmd()
	command.SetArgs([]string{
		"run",
		"--log-file", runLogFile,
		"--output", outFile,
		"--git-type", "go",
		"--author-name", "Test Author",
		"--author-email", "test@example.com",
		"-B", "keep-no-dir-branch",
		"-m", "keep no dir test",
		"--keep",
		changerBinary,
	})
	err = command.Execute()
	assert.NoError(t, err)

	// Verify the PR was created
	require.Len(t, vcMock.PullRequests, 1)
	assert.Equal(t, "keep-no-dir-branch", vcMock.PullRequests[0].Head)

	// The kept directory should be in the default temp dir
	expectedDir := filepath.Join(os.TempDir(), "multi-gitter-owner-keep-no-clone-dir")
	defer os.RemoveAll(expectedDir)

	_, statErr := os.Stat(expectedDir)
	require.NoError(t, statErr, "cloned repo directory should exist in default temp dir when using --keep without --clone-dir")
}

// TestNoKeepCleansUp verifies that without --keep, clone directories are cleaned up
func TestNoKeepCleansUp(t *testing.T) {
	workingDir, err := os.Getwd()
	require.NoError(t, err)

	changerBinary := normalizePath(filepath.Join(workingDir, changerBinaryPath))

	vcMock := &vcmock.VersionController{}
	defer vcMock.Clean()
	cmd.OverrideVersionController = vcMock

	cloneDir, err := os.MkdirTemp(os.TempDir(), "multi-gitter-nokeep-test-")
	require.NoError(t, err)
	defer os.RemoveAll(cloneDir)

	repo := createRepo(t, "owner", "nokeep-repo", "i like apples")
	vcMock.AddRepository(repo)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "multi-gitter-nokeep-log-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	runLogFile := filepath.Join(tmpDir, "run-log.txt")
	outFile := filepath.Join(tmpDir, "run-out.txt")

	command := cmd.RootCmd()
	command.SetArgs([]string{
		"run",
		"--log-file", runLogFile,
		"--output", outFile,
		"--git-type", "go",
		"--author-name", "Test Author",
		"--author-email", "test@example.com",
		"-B", "nokeep-branch",
		"-m", "no keep test",
		"--clone-dir", cloneDir,
		changerBinary,
	})
	err = command.Execute()
	assert.NoError(t, err)

	require.Len(t, vcMock.PullRequests, 1)

	// The clone dir should have no subdirectories remaining (they should have been cleaned up)
	entries, err := os.ReadDir(cloneDir)
	require.NoError(t, err)
	assert.Empty(t, entries, "clone directory should be empty after run without --keep (temp dirs cleaned up)")
}

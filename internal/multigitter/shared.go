package multigitter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/lindell/multi-gitter/internal/git"
	"github.com/pkg/errors"
)

type urler interface {
	URL() string
}

func transformExecError(err error) error {
	var sysErr syscall.Errno
	if ok := errors.As(err, &sysErr); ok {
		if sysErr.Error() == "exec format error" {
			return errors.New("the script or program is in the wrong format")
		}
	}
	return err
}

// Git is a git implementation
type Git interface {
	Clone(ctx context.Context, url string, baseName string) error
	ChangeBranch(branchName string) error
	Changes() (bool, error)
	Commit(commitAuthor *git.CommitAuthor, commitMessage string) error
	BranchExist(remoteName, branchName string) (bool, error)
	Push(ctx context.Context, remoteName string, force bool) error
	AddRemote(name, url string) error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func getStackTrace(err error) string {
	if err, ok := err.(stackTracer); ok {
		trace := ""
		for _, f := range err.StackTrace() {
			trace += fmt.Sprintf("%+s:%d\n", f, f)
		}
		return trace
	}
	return ""
}

// ConflictStrategy define how a conflict of an already existing branch should be handled
type ConflictStrategy int

const (
	// ConflictStrategySkip will skip the run for if the branch does already exist
	ConflictStrategySkip ConflictStrategy = iota + 1
	// ConflictStrategyReplace will ignore any existing branch and replace it with new changes
	ConflictStrategyReplace
)

// ParseConflictStrategy parses a conflict strategy from a string
func ParseConflictStrategy(str string) (ConflictStrategy, error) {
	switch str {
	default:
		return ConflictStrategy(0), fmt.Errorf("could not parse \"%s\" as conflict strategy", str)
	case "skip":
		return ConflictStrategySkip, nil
	case "replace":
		return ConflictStrategyReplace, nil
	}
}

// createTempDir creates a temporary directory in the given directory.
// If the given directory is an empty string, it will use the os.TempDir()
func createTempDir(cloneDir string) (string, error) {
	if cloneDir == "" {
		cloneDir = os.TempDir()
	}

	absDir, err := makeAbsolutePath(cloneDir)
	if err != nil {
		return "", err
	}

	err = createDirectoryIfDoesntExist(absDir)
	if err != nil {
		return "", err
	}

	tmpDir, err := os.MkdirTemp(absDir, "multi-git-changer-")
	if err != nil {
		return "", err
	}

	return tmpDir, nil
}

func createDirectoryIfDoesntExist(directoryPath string) error {
	// Check if the directory exists
	if _, err := os.Stat(directoryPath); !os.IsNotExist(err) {
		return nil
	}

	// Create the directory
	err := os.MkdirAll(directoryPath, 0600)
	if err != nil {
		return err
	}

	return nil
}

// makeAbsolutePath creates an absolute path from a relative path
func makeAbsolutePath(path string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "could not get the working directory")
	}

	if !filepath.IsAbs(path) {
		return filepath.Join(workingDir, path), nil
	}

	return path, nil
}

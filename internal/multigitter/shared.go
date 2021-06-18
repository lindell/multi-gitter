package multigitter

import (
	"fmt"
	"syscall"

	"github.com/lindell/multi-gitter/internal/domain"
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
	Clone(url string, baseName string) error
	ChangeBranch(branchName string) error
	Changes() (bool, error)
	Commit(commitAuthor *domain.CommitAuthor, commitMessage string) error
	BranchExist(remoteName, branchName string) (bool, error)
	Push(remoteName string) error
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

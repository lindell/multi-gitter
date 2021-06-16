package multigitter

import (
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

type git interface {
	Clone(url string, baseName string) error
	ChangeBranch(branchName string) error
	Changes() (bool, error)
	Commit(commitAuthor *domain.CommitAuthor, commitMessage string) error
	BranchExist(remoteName, branchName string) (bool, error)
	Push(remoteName string) error
	AddRemote(name, url string) error
}

package multigitter

import (
	"syscall"

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

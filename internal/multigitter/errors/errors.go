package errors

import "errors"

// Constant errors used in multi-gitter
var (
	ErrAborted     = errors.New("run was never started because of aborted execution")
	ErrRejected    = errors.New("changes were not included since they were manually rejected")
	ErrNoChange    = errors.New("no data was changed")
	ErrBranchExist = errors.New("the new branch does already exist")
)

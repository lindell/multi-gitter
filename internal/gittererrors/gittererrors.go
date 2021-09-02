package gittererrors

// Error exist to create constant errors
type Error string

func (e Error) Error() string {
	return string(e)
}

// Constant errors
const (
	NoChangeError    Error = "no data was changed"
	ExitCodeError    Error = "the program exited with a non zero exit code"
	BranchExistError Error = "the new branch does already exist"
)

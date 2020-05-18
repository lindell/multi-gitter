package domain

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	NoChangeError Error = "no data was changed"
	ExitCodeError Error = "the program exited with a non zero exit code"
)

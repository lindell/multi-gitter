package git

// CommitAuthor is the data (name and email) used when a commit is made
type CommitAuthor struct {
	Name  string
	Email string
}

type Changes struct {
	Additions []Change
	Deletions []Change
}

type Change struct {
	Path     string
	Contents []byte
}

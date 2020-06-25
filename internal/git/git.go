package git

import (
	"net/url"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/lindell/multi-gitter/internal/domain"
)

// Git is an implementation of git that executes git as a command
// This has drawbacks, but the big benefit is that the configuration probably already present can be reused
type Git struct {
	Directory string // The (temporary) directory that should be worked within
	Repo      string // The "url" to the repo, any format can be used as long as it's pushable
	NewBranch string // The name of the new branch that new changes will be pushed to
	Token     string
}

// Clone a repository
func (g Git) Clone() error {
	u, err := url.Parse(g.Repo)
	if err != nil {
		return err
	}

	// Set the token as https://TOKEN@url
	u.User = url.User(g.Token)

	r, err := git.PlainClone(g.Directory, false, &git.CloneOptions{
		URL:        u.String(),
		RemoteName: "origin",
	})
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(g.NewBranch),
		Create: true,
	})
	if err != nil {
		return err
	}

	return err
}

// Commit and push if any changes has been made to the directory
func (g Git) Commit(commitMessage string) error {
	r, err := git.PlainOpen(g.Directory)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	status, err := w.Status()
	if err != nil {
		return err
	}

	if status.IsClean() {
		return domain.NoChangeError
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	_, err = w.Commit(commitMessage, &git.CommitOptions{})
	if err != nil {
		return err
	}

	err = r.Push(&git.PushOptions{})
	if err != nil {
		return err
	}

	return nil
}

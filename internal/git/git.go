package git

import (
	"net/url"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Git is an implementation of git that executes git as a command
// This has drawbacks, but the big benefit is that the configuration probably already present can be reused
type Git struct {
	Directory string // The (temporary) directory that should be worked within
	Repo      string // The "url" to the repo, any format can be used as long as it's pushable
	NewBranch string // The name of the new branch that new changes will be pushed to
	Token     string

	repo *git.Repository // The repository after the clone has been made
}

// Clone a repository
func (g *Git) Clone() error {
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
	g.repo = r

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

// Changes detect if any changes has been made in the directory
func (g *Git) Changes() (bool, error) {
	w, err := g.repo.Worktree()
	if err != nil {
		return false, err
	}

	status, err := w.Status()
	if err != nil {
		return false, err
	}

	return !status.IsClean(), nil
}

// Commit and push all changes
func (g *Git) Commit(commitMessage string) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	_, err = w.Commit(commitMessage, &git.CommitOptions{})
	if err != nil {
		return err
	}

	err = g.repo.Push(&git.PushOptions{})
	if err != nil {
		return err
	}

	return nil
}

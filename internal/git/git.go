package git

import (
	"bytes"
	"fmt"
	"net/url"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
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
		Depth:      10,
	})
	if err != nil {
		return err
	}
	g.repo = r

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

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(g.NewBranch),
		Create: true,
		Keep:   true,
	})
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	// Get the current hash to be able to diff it with the committed changes later
	oldHead, err := g.repo.Head()
	if err != nil {
		return err
	}
	oldHash := oldHead.Hash()

	hash, err := w.Commit(commitMessage, &git.CommitOptions{})
	if err != nil {
		return err
	}

	commit, err := g.repo.CommitObject(hash)
	if err != nil {
		return err
	}

	_ = g.logDiff(oldHash, commit.Hash)

	return nil
}

func (g *Git) logDiff(aHash, bHash plumbing.Hash) error {
	if !log.IsLevelEnabled(log.GetLevel()) {
		return nil
	}

	aCommit, err := g.repo.CommitObject(aHash)
	if err != nil {
		return err
	}
	aTree, err := aCommit.Tree()
	if err != nil {
		return err
	}

	bCommit, err := g.repo.CommitObject(bHash)
	if err != nil {
		return err
	}
	bTree, err := bCommit.Tree()
	if err != nil {
		return err
	}

	patch, err := aTree.Patch(bTree)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	err = patch.Encode(buf)
	if err != nil {
		return err
	}
	log.Debug(buf.String())

	return nil
}

// BranchExist checks if the new branch exists
func (g *Git) BranchExist() (bool, error) {
	_, err := g.repo.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", g.NewBranch)), false)
	if err == plumbing.ErrReferenceNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Push the committed changes to the remote
func (g *Git) Push() error {
	return g.repo.Push(&git.PushOptions{})
}

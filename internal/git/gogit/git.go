package gogit

import (
	"bytes"
	"time"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/go-git/go-git/v5/plumbing/object"
	internalgit "github.com/lindell/multi-gitter/internal/git"
	"github.com/pkg/errors"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
)

// Git is an implementation of git that used go-git
type Git struct {
	Directory  string // The (temporary) directory that should be worked within
	FetchDepth int    // Limit fetching to the specified number of commits

	repo *git.Repository // The repository after the clone has been made
}

// Clone a repository
func (g *Git) Clone(url string, baseName string) error {
	r, err := git.PlainClone(g.Directory, false, &git.CloneOptions{
		URL:           url,
		RemoteName:    "origin",
		Depth:         g.FetchDepth,
		ReferenceName: plumbing.NewBranchReferenceName(baseName),
		SingleBranch:  true,
	})
	if err != nil {
		return errors.Wrap(err, "could not clone from the remote")
	}

	g.repo = r

	return nil
}

// ChangeBranch changes the branch
func (g *Git) ChangeBranch(branchName string) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Create: true,
	})
	if err != nil {
		return err
	}

	return nil
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
func (g *Git) Commit(commitAuthor *internalgit.CommitAuthor, commitMessage string) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	// Make sure gitignore is used
	patterns, err := gitignore.ReadPatterns(w.Filesystem, nil)
	if err != nil {
		return err
	}
	w.Excludes = patterns

	err = w.AddWithOptions(&git.AddOptions{
		All: true,
	})
	if err != nil {
		return err
	}

	// Get the current hash to be able to diff it with the committed changes later
	oldHead, err := g.repo.Head()
	if err != nil {
		return err
	}
	oldHash := oldHead.Hash()

	var author *object.Signature
	if commitAuthor != nil {
		author = &object.Signature{
			Name:  commitAuthor.Name,
			Email: commitAuthor.Email,
			When:  time.Now(),
		}
	}

	hash, err := w.Commit(commitMessage, &git.CommitOptions{
		Author: author,
	})
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
	if !log.IsLevelEnabled(log.DebugLevel) {
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
func (g *Git) BranchExist(remoteName, branchName string) (bool, error) {
	remote, err := g.repo.Remote(remoteName)
	if err != nil {
		return false, err
	}

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return false, err
	}
	for _, r := range refs {
		if r.Name().Short() == branchName {
			return true, nil
		}
	}
	return false, nil
}

// Push the committed changes to the remote
func (g *Git) Push(remoteName string) error {
	return g.repo.Push(&git.PushOptions{
		RemoteName: remoteName,
	})
}

// AddRemote adds a new remote
func (g *Git) AddRemote(name, url string) error {
	_, err := g.repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	})
	return err
}

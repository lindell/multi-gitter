package git

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/lindell/multi-gitter/internal/domain"
	"github.com/pkg/errors"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	log "github.com/sirupsen/logrus"
)

// Git is an implementation of git that executes git as a command
// This has drawbacks, but the big benefit is that the configuration probably already present can be reused
type Git struct {
	Directory  string // The (temporary) directory that should be worked within
	Repo       string // The "url" to the repo, any format can be used as long as it's pushable
	FetchDepth int    // Limit fetching to the specified number of commits

	repo *git.Repository // The repository after the clone has been made
}

// Clone a repository
func (g *Git) Clone(baseName, headName string) error {
	r, err := git.PlainClone(g.Directory, false, &git.CloneOptions{
		URL:           g.Repo,
		RemoteName:    "origin",
		Depth:         g.FetchDepth,
		ReferenceName: plumbing.NewBranchReferenceName(baseName),
		SingleBranch:  true,
	})
	if err != nil {
		return errors.Wrap(err, "could not clone from the remote")
	}

	if headName != "" {
		err = r.Fetch(&git.FetchOptions{
			RefSpecs: []config.RefSpec{
				config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/remotes/origin/%s", headName, headName)),
			},
			Depth: g.FetchDepth,
		})
		if err != nil {
			if _, ok := err.(git.NoMatchingRefSpecError); !ok {
				return err
			}
		}
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
func (g *Git) Commit(commitAuthor *domain.CommitAuthor, commitMessage string) error {
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
func (g *Git) BranchExist(branchName string) (bool, error) {
	_, err := g.repo.Reference(plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", branchName)), false)
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

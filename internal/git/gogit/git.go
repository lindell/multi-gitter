package gogit

import (
	"bytes"
	"context"
	"io"
	"slices"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	internalgit "github.com/lindell/multi-gitter/internal/git"
	"github.com/pkg/errors"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	log "github.com/sirupsen/logrus"
)

// Git is an implementation of git that used go-git
type Git struct {
	Directory   string // The (temporary) directory that should be worked within
	FetchDepth  int    // Limit fetching to the specified number of commits
	Credentials *internalgit.Credentials

	repo *git.Repository // The repository after the clone has been made
}

// Clone a repository
func (g *Git) Clone(ctx context.Context, url string, baseName string) error {
	r, err := git.PlainCloneContext(ctx, g.Directory, false, &git.CloneOptions{
		URL:           url,
		RemoteName:    "origin",
		Depth:         g.FetchDepth,
		ReferenceName: plumbing.NewBranchReferenceName(baseName),
		SingleBranch:  true,
		Auth:          g.credentials(),
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

	status, err := w.Status()
	if err != nil {
		return err
	}

	// This is a workaround for a bug in go-git where "add all" does not add deleted files
	// If https://github.com/go-git/go-git/issues/223 is fixed, this can be removed
	for file, s := range status {
		if s.Worktree == git.Deleted {
			_, err = w.Add(file)
			if err != nil {
				return err
			}
		}
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

	refs, err := remote.List(&git.ListOptions{
		Auth: g.credentials(),
	})
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
func (g *Git) Push(ctx context.Context, remoteName, remoteReference string, force bool) error {
	var refSpecs []config.RefSpec

	if remoteReference != "" {
		// go-git doesn't support refSpec like HEAD:<name>, so first we need to resolve SHA1 commit related to HEAD
		head, err := g.repo.Head()
		if err != nil {
			return errors.Wrap(err, "Unable to get HEAD")
		}
		refSpecs = []config.RefSpec{
			config.RefSpec(head.Hash().String() + ":" + remoteReference),
		}
	}
	return g.repo.PushContext(ctx, &git.PushOptions{
		RemoteName: remoteName,
		Force:      force,
		Auth:       g.credentials(),
		RefSpecs:   refSpecs,
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

func (g *Git) credentials() *http.BasicAuth {
	if g.Credentials == nil {
		return nil
	}

	return &http.BasicAuth{
		Username: g.Credentials.Username,
		Password: g.Credentials.Password,
	}
}

// LatestCommitHash returns the latest commit hash
func (g *Git) LatestCommitHash() (string, error) {
	head, err := g.repo.Head()
	if err != nil {
		return "", err
	}
	return head.Hash().String(), nil
}

// ChangesSinceCommit returns the changes made in commits since the given commit hash
func (g *Git) ChangesSinceCommit(sinceCommitHash string) ([]internalgit.Changes, error) {
	iter, err := g.repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	toCommit, err := iter.Next()
	if err != nil {
		return nil, errors.WithMessage(err, "could not get current commit")
	}

	// Go through all commits until we reach the sinceCommitHash
	// This is is by default only the latest commit, but if we use manual commits,
	// there might be multiple commits to go through.
	allChanges := []internalgit.Changes{}
	for {
		fromCommit, err := iter.Next()
		if err != nil {
			return nil, errors.WithMessage(err, "could not get last commit")
		}

		changes, err := g.changesBetweenCommits(context.Background(), fromCommit, toCommit)
		if err != nil {
			return nil, errors.WithMessage(err, "could not get changes")
		}
		allChanges = append(allChanges, changes)

		if sinceCommitHash == fromCommit.Hash.String() {
			break
		}
		toCommit = fromCommit
	}

	// Reverse the order of the changes to get the earliest commit first
	slices.Reverse(allChanges)

	return allChanges, nil
}

func (g *Git) changesBetweenCommits(_ context.Context, from, to *object.Commit) (internalgit.Changes, error) {
	toTree, err := to.Tree()
	if err != nil {
		return internalgit.Changes{}, errors.WithMessage(err, "could not get current tree")
	}
	fromTree, err := from.Tree()
	if err != nil {
		return internalgit.Changes{}, errors.WithMessage(err, "could not get current tree")
	}

	changes, err := fromTree.Diff(toTree)
	if err != nil {
		return internalgit.Changes{}, errors.WithMessage(err, "could not get diff")
	}

	additions := map[string][]byte{}
	deletions := []string{}
	for _, change := range changes {
		action, err := change.Action()
		if err != nil {
			return internalgit.Changes{}, errors.WithMessage(err, "could not get action")
		}

		if action == merkletrie.Insert || action == merkletrie.Modify {
			_, to, err := change.Files()
			if err != nil {
				return internalgit.Changes{}, errors.WithMessage(err, "could not get files")
			}

			reader, err := to.Reader()
			if err != nil {
				return internalgit.Changes{}, errors.WithMessage(err, "could not get reader")
			}
			bytes, err := io.ReadAll(reader)
			reader.Close()
			if err != nil {
				return internalgit.Changes{}, errors.WithMessage(err, "could not read file")
			}

			additions[change.To.Name] = bytes
		} else if action == merkletrie.Delete {
			deletions = append(deletions, change.From.Name)
		}
	}

	return internalgit.Changes{
		Message:   strings.TrimSpace(to.Message),
		Additions: additions,
		Deletions: deletions,
		OldHash:   from.Hash.String(),
	}, nil
}

package tests

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/lindell/multi-gitter/tests/vcmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRepo(t *testing.T, name, dataInFile string) vcmock.Repository {
	tmpDir, err := createDummyRepo(dataInFile)
	require.NoError(t, err)

	return vcmock.Repository{
		Name: name,
		Path: tmpDir,
	}
}

func createDummyRepo(dataInFile string) (string, error) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "multi-git-test-")
	if err != nil {
		return "", err
	}

	repo, err := git.PlainInit(tmpDir, false)
	if err != nil {
		return "", err
	}

	testFilePath := path.Join(tmpDir, fileName)

	err = ioutil.WriteFile(testFilePath, []byte(dataInFile), 0600)
	if err != nil {
		return "", err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	if _, err = wt.Add("."); err != nil {
		return "", err
	}

	_, err = wt.Commit("First commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", err
	}

	return tmpDir, nil
}

func changeBranch(t *testing.T, path string, branchName string) {
	repo, err := git.PlainOpen(path)
	assert.NoError(t, err)

	wt, err := repo.Worktree()
	assert.NoError(t, err)

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
	})
	assert.NoError(t, err)
}

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

const fileName = "test.txt"

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

func changeBranch(t *testing.T, path string, branchName string, create bool) {
	repo, err := git.PlainOpen(path)
	assert.NoError(t, err)

	wt, err := repo.Worktree()
	assert.NoError(t, err)

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Create: create,
	})
	assert.NoError(t, err)
}

func changeTestFile(t *testing.T, basePath string, content string, commitMessage string) {
	repo, err := git.PlainOpen(basePath)
	require.NoError(t, err)

	testFilePath := path.Join(basePath, fileName)

	err = ioutil.WriteFile(testFilePath, []byte(content), 0600)
	require.NoError(t, err)

	wt, err := repo.Worktree()
	require.NoError(t, err)

	_, err = wt.Add(".")
	require.NoError(t, err)

	_, err = wt.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)
}

func readTestFile(t *testing.T, basePath string) string {
	testFilePath := path.Join(basePath, fileName)

	b, err := ioutil.ReadFile(testFilePath)
	require.NoError(t, err)

	return string(b)
}

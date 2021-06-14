package cmdgit

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/lindell/multi-gitter/internal/domain"
)

type Git struct {
	Directory  string // The (temporary) directory that should be worked within
	FetchDepth int    // Limit fetching to the specified number of commits
}

var errRe = regexp.MustCompile(`(^|\n)(error|fatal): (.+)`)

func (g *Git) run(args ...string) (string, string, error) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	cmd := exec.Command("git", args...)
	cmd.Dir = g.Directory
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	err := cmd.Run()
	if err != nil {
		msg := fmt.Sprintf(`"%s" existed with %d`,
			strings.Join(append([]string{"git"}, args...), " "),
			cmd.ProcessState.ExitCode(),
		)

		fmt.Println(stderr.String())

		matches := errRe.FindStringSubmatch(stderr.String())
		if matches != nil {
			msg += ": " + matches[3]
		}

		return "", "", errors.New(msg)
	}
	return stdout.String(), stderr.String(), nil
}

// Clone
func (g *Git) Clone(url string, baseName string) error {
	args := []string{"clone", url, "--branch", baseName, "--single-branch"}
	if g.FetchDepth > 0 {
		args = append(args, "--depth", fmt.Sprint(g.FetchDepth))
	}
	args = append(args, g.Directory)

	_, _, err := g.run(args...)
	return err
}

// ChangeBranch
func (g *Git) ChangeBranch(branchName string) error {
	_, _, err := g.run("checkout", "-b", branchName)
	return err
}

// Changes
func (g *Git) Changes() (bool, error) {
	stdOut, _, err := g.run("status", "-s")
	return len(stdOut) > 0, err
}

// Commit
func (g *Git) Commit(commitAuthor *domain.CommitAuthor, commitMessage string) error {
	_, _, err := g.run("add", ".")
	if err != nil {
		return err
	}

	_, _, err = g.run("commit", "-m", commitMessage)
	return err
}

// BranchExist
func (g *Git) BranchExist(remoteName, branchName string) (bool, error) {
	stdOut, _, err := g.run("ls-remote", "-q", "-h")
	if err != nil {
		return false, err
	}
	return strings.Contains(stdOut, fmt.Sprintf("refs/heads/%s", branchName)), nil
}

// Push
func (g *Git) Push(remoteName string) error {
	_, _, err := g.run("push", remoteName, "HEAD")
	return err
}

// AddRemote
func (g *Git) AddRemote(name, url string) error {
	_, _, err := g.run("remote", "add", name, url)
	return err
}

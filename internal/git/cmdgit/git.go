package cmdgit

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/lindell/multi-gitter/internal/domain"
)

type Git struct {
	Directory  string // The (temporary) directory that should be worked within
	FetchDepth int    // Limit fetching to the specified number of commits
}

var errRe = regexp.MustCompile(`(^|\n)(error|fatal): (.+)`)

func (g *Git) run(cmd *exec.Cmd) (string, string, error) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	cmd.Dir = g.Directory
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	err := cmd.Run()
	if err != nil {
		matches := errRe.FindStringSubmatch(stderr.String())
		if matches != nil {
			return "", "", errors.New(matches[3])
		}

		msg := fmt.Sprintf(`git command existed with %d`,
			cmd.ProcessState.ExitCode(),
		)

		return "", "", errors.New(msg)
	}
	return stdout.String(), stderr.String(), nil
}

// Clone a repository
func (g *Git) Clone(url string, baseName string) error {
	args := []string{"clone", url, "--branch", baseName, "--single-branch"}
	if g.FetchDepth > 0 {
		args = append(args, "--depth", fmt.Sprint(g.FetchDepth))
	}
	args = append(args, g.Directory)

	cmd := exec.Command("git", args...)
	_, _, err := g.run(cmd)
	return err
}

// ChangeBranch changes the branch
func (g *Git) ChangeBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	_, _, err := g.run(cmd)
	return err
}

// Changes detect if any changes has been made in the directory
func (g *Git) Changes() (bool, error) {
	cmd := exec.Command("git", "status", "-s")
	stdOut, _, err := g.run(cmd)
	return len(stdOut) > 0, err
}

// Commit and push all changes
func (g *Git) Commit(commitAuthor *domain.CommitAuthor, commitMessage string) error {
	cmd := exec.Command("git", "add", ".")
	_, _, err := g.run(cmd)
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "commit", "-m", commitMessage)

	if commitAuthor != nil {
		cmd.Env = append(cmd.Env,
			"GIT_AUTHOR_NAME="+commitAuthor.Name,
			"GIT_AUTHOR_NAME="+commitAuthor.Email,
			"GIT_COMMITTER_NAME="+commitAuthor.Name,
			"GIT_COMMITTER_EMAIL="+commitAuthor.Email,
		)
	}

	_, _, err = g.run(cmd)

	if err := g.logDiff(); err != nil {
		return err
	}

	return err
}

func (g *Git) logDiff() error {
	if !log.IsLevelEnabled(log.DebugLevel) {
		return nil
	}

	cmd := exec.Command("git", "diff", "HEAD~1")
	stdout, _, err := g.run(cmd)
	if err != nil {
		return err
	}

	log.Debug(stdout)

	return nil
}

// BranchExist checks if the new branch exists
func (g *Git) BranchExist(remoteName, branchName string) (bool, error) {
	cmd := exec.Command("git", "ls-remote", "-q", "-h")
	stdOut, _, err := g.run(cmd)
	if err != nil {
		return false, err
	}
	return strings.Contains(stdOut, fmt.Sprintf("refs/heads/%s", branchName)), nil
}

// Push the committed changes to the remote
func (g *Git) Push(remoteName string) error {
	cmd := exec.Command("git", "push", remoteName, "HEAD")
	_, _, err := g.run(cmd)
	return err
}

// AddRemote adds a new remote
func (g *Git) AddRemote(name, url string) error {
	cmd := exec.Command("git", "remote", "add", name, url)
	_, _, err := g.run(cmd)
	return err
}

package cmdgit

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/lindell/multi-gitter/internal/git"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Git is an implementation of git that executes git as commands
type Git struct {
	Directory  string // The (temporary) directory that should be worked within
	FetchDepth int    // Limit fetching to the specified number of commits
}

var errRe = regexp.MustCompile(`(^|\n)(error|fatal): (.+)`)

func (g *Git) run(cmd *exec.Cmd) (string, error) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	cmd.Dir = g.Directory
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	err := cmd.Run()
	if err != nil {
		matches := errRe.FindStringSubmatch(stderr.String())
		if matches != nil {
			return "", errors.New(matches[3])
		}

		msg := fmt.Sprintf(`git command exited with %d (%s)`,
			cmd.ProcessState.ExitCode(),
			stderr.String(),
		)

		return "", errors.New(msg)
	}
	return stdout.String(), nil
}

// Clone a repository
func (g *Git) Clone(ctx context.Context, url string, baseName string) error {
	args := []string{"clone", url, "--branch", baseName, "--single-branch"}
	if g.FetchDepth > 0 {
		args = append(args, "--depth", fmt.Sprint(g.FetchDepth))
	}
	args = append(args, g.Directory)

	cmd := exec.CommandContext(ctx, "git", args...)
	_, err := g.run(cmd)
	return err
}

// ChangeBranch changes the branch
func (g *Git) ChangeBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	_, err := g.run(cmd)
	return err
}

// Changes detect if any changes has been made in the directory
func (g *Git) Changes() (bool, error) {
	cmd := exec.Command("git", "status", "-s")
	stdOut, err := g.run(cmd)
	return len(stdOut) > 0, err
}

// Commit and push all changes
func (g *Git) Commit(commitAuthor *git.CommitAuthor, commitMessage string) error {
	cmd := exec.Command("git", "add", ".")
	_, err := g.run(cmd)
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "commit", "--no-verify", "-m", commitMessage)

	if commitAuthor != nil {
		cmd.Env = append(cmd.Env,
			"GIT_AUTHOR_NAME="+commitAuthor.Name,
			"GIT_AUTHOR_EMAIL="+commitAuthor.Email,
			"GIT_COMMITTER_NAME="+commitAuthor.Name,
			"GIT_COMMITTER_EMAIL="+commitAuthor.Email,
		)
	}

	_, err = g.run(cmd)
	if err != nil {
		return err
	}

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
	stdout, err := g.run(cmd)
	if err != nil {
		return err
	}

	log.Debug(stdout)

	return nil
}

// BranchExist checks if the new branch exists
func (g *Git) BranchExist(remoteName, branchName string) (bool, error) {
	cmd := exec.Command("git", "ls-remote", "-q", "-h", remoteName)
	stdOut, err := g.run(cmd)
	if err != nil {
		return false, err
	}
	return strings.Contains(stdOut, fmt.Sprintf("\trefs/heads/%s\n", branchName)), nil
}

// Push the committed changes to the remote
func (g *Git) Push(ctx context.Context, remoteName string, force bool) error {
	args := []string{"push", "--no-verify", remoteName}
	if force {
		args = append(args, "--force")
	}
	args = append(args, "HEAD")

	cmd := exec.CommandContext(ctx, "git", args...)
	_, err := g.run(cmd)
	return err
}

// AddRemote adds a new remote
func (g *Git) AddRemote(name, url string) error {
	cmd := exec.Command("git", "remote", "add", name, url)
	_, err := g.run(cmd)
	return err
}

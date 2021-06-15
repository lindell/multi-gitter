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

func (g *Git) run(args ...string) (string, string, error) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	cmd := exec.Command("git", args...)
	cmd.Dir = g.Directory
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	err := cmd.Run()
	if err != nil {
		matches := errRe.FindStringSubmatch(stderr.String())
		if matches != nil {
			return "", "", errors.New(matches[3])
		}

		msg := fmt.Sprintf(`"%s" existed with %d`,
			strings.Join(append([]string{"git"}, args...), " "),
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

	_, _, err := g.run(args...)
	return err
}

// ChangeBranch changes the branch
func (g *Git) ChangeBranch(branchName string) error {
	_, _, err := g.run("checkout", "-b", branchName)
	return err
}

// Changes detect if any changes has been made in the directory
func (g *Git) Changes() (bool, error) {
	stdOut, _, err := g.run("status", "-s")
	return len(stdOut) > 0, err
}

// Commit and push all changes
func (g *Git) Commit(commitAuthor *domain.CommitAuthor, commitMessage string) error {
	_, _, err := g.run("add", ".")
	if err != nil {
		return err
	}

	args := []string{"commit", "-m", commitMessage}

	if commitAuthor != nil {
		args = append(args, fmt.Sprintf(`--author="%s"`, fmt.Sprintf("%s <%s>", commitAuthor.Name, commitAuthor.Email)))
	}

	fmt.Println("MARKLAR", args)

	_, _, err = g.run(args...)

	if err := g.logDiff(); err != nil {
		return err
	}

	return err
}

func (g *Git) logDiff() error {
	if !log.IsLevelEnabled(log.DebugLevel) {
		return nil
	}

	stdout, _, err := g.run("diff", "HEAD~1")
	if err != nil {
		return err
	}

	log.Debug(stdout)

	return nil
}

// BranchExist checks if the new branch exists
func (g *Git) BranchExist(remoteName, branchName string) (bool, error) {
	stdOut, _, err := g.run("ls-remote", "-q", "-h")
	if err != nil {
		return false, err
	}
	return strings.Contains(stdOut, fmt.Sprintf("refs/heads/%s", branchName)), nil
}

// Push the committed changes to the remote
func (g *Git) Push(remoteName string) error {
	_, _, err := g.run("push", remoteName, "HEAD")
	return err
}

// AddRemote adds a new remote
func (g *Git) AddRemote(name, url string) error {
	_, _, err := g.run("remote", "add", name, url)
	return err
}

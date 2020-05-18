package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/lindell/multi-gitter/internal/domain"
)

// Git is an implementation of git that executes git as a command
// This has drawbacks, but the big benefit is that the configuration probably already present can be reused
type Git struct {
	Directory string // The (temporary) directory that should be worked within
	Repo      string // The "url" to the repo, any format can be used as long as it's pushable
	NewBranch string // The name of the new branch that new changes will be pushed to
}

// errorWrap converts errors a failed command into more a more useful error
func errorWrap(err error) error {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("git command existed with status code %d:\n%s\n", exitErr.ExitCode(), exitErr.Stderr)
	}
	return err
}

func (g Git) Clone() error {
	cmd := g.command("git", "clone", g.Repo, g.Directory)
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (g Git) Commit(commitMessage string) error {
	cmd := g.command("git", "add", ".")
	cmd.Dir = g.Directory
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = g.command("git", "checkout", "-b", g.NewBranch)
	cmd.Dir = g.Directory
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = g.command("git", "commit", "-F", "-")
	cmd.Dir = g.Directory
	cmd.Stdin = strings.NewReader(commitMessage)
	err = cmd.Run()
	if err != nil {
		return domain.NoChangeError
	}

	cmd = g.command("git", "push", "-u", "origin", g.NewBranch)
	cmd.Dir = g.Directory
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

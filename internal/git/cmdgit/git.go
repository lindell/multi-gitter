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
	logGitExecution(cmd, stdout, stderr)
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

func logGitExecution(cmd *exec.Cmd, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	log.WithFields(log.Fields{
		"cmd":    cmd.String(),
		"stdout": stdout.String(),
		"stderr": stderr.String(),
	}).Trace("cmdgit")
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

// FetchAndResetToDefault fetches the latest changes from origin and performs a hard reset
// to the default branch. This is used with the --keep option to reuse already-cloned repositories.
func (g *Git) FetchAndResetToDefault(ctx context.Context, baseName string) error {
	// Checkout the base branch (in case we are on a feature branch)
	cmd := exec.CommandContext(ctx, "git", "checkout", baseName)
	_, err := g.run(cmd)
	if err != nil {
		return errors.Wrap(err, "could not checkout base branch")
	}

	// Fetch latest changes
	cmd = exec.CommandContext(ctx, "git", "fetch", "origin", baseName)
	_, err = g.run(cmd)
	if err != nil {
		return errors.Wrap(err, "could not fetch from origin")
	}

	// Hard reset to origin's base branch
	cmd = exec.CommandContext(ctx, "git", "reset", "--hard", "origin/"+baseName)
	_, err = g.run(cmd)
	if err != nil {
		return errors.Wrap(err, "could not hard reset to origin base branch")
	}

	// Clean untracked files
	cmd = exec.CommandContext(ctx, "git", "clean", "-fd")
	_, err = g.run(cmd)
	if err != nil {
		return errors.Wrap(err, "could not clean untracked files")
	}

	return nil
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
func (g *Git) Push(ctx context.Context, remoteName, remoteReference string, force bool) error {
	args := []string{"push", "--no-verify", remoteName}
	if force {
		args = append(args, "--force")
	}
	refSpec := "HEAD"
	if remoteReference != "" {
		refSpec = refSpec + ":" + remoteReference
	}
	args = append(args, refSpec)

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

// LatestCommitHash returns the latest commit hash
func (g *Git) LatestCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	stdOut, err := g.run(cmd)
	return strings.TrimSpace(stdOut), err
}

// ChangesSinceCommit returns the changes made in commits since the given commit hash
func (g *Git) ChangesSinceCommit(sinceCommitHash string) ([]git.Changes, error) {
	// Get the list of commits from sinceCommitHash to HEAD, one for each line
	cmd := exec.Command("git", "rev-list", "--reverse", sinceCommitHash+"..HEAD")
	stdOut, err := g.run(cmd)
	if err != nil {
		return nil, errors.WithMessage(err, "could not get commit list")
	}

	commitHashes := strings.Split(strings.TrimSpace(stdOut), "\n")
	if len(commitHashes) == 0 || commitHashes[0] == "" {
		return nil, errors.New("no commits found")
	}

	allChanges := []git.Changes{}
	previousHash := sinceCommitHash

	for _, commitHash := range commitHashes {
		commitHash = strings.TrimSpace(commitHash)
		if commitHash == "" {
			continue
		}

		changes, err := g.getChangesBetweenCommits(previousHash, commitHash)
		if err != nil {
			return nil, err
		}

		allChanges = append(allChanges, changes)
		previousHash = commitHash
	}

	return allChanges, nil
}

func (g *Git) getChangesBetweenCommits(fromHash, toHash string) (git.Changes, error) {
	// Get commit message
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%B", toHash)
	commitMessage, err := g.run(cmd)
	if err != nil {
		return git.Changes{}, errors.WithMessage(err, "could not get commit message")
	}
	commitMessage = strings.TrimRight(commitMessage, "\n")

	// Get the list of files changed
	cmd = exec.Command("git", "diff", "--name-status", fromHash, toHash)
	diffOutput, err := g.run(cmd)
	if err != nil {
		return git.Changes{}, errors.WithMessage(err, "could not get diff")
	}

	additions := map[string][]byte{}
	deletions := []string{}

	for _, line := range strings.Split(strings.TrimSpace(diffOutput), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		status := parts[0]
		filePath := parts[1]

		switch status {
		case "A", "M": // Added or Modified
			cmd = exec.Command("git", "show", toHash+":"+filePath)
			content, err := g.run(cmd)
			if err != nil {
				return git.Changes{}, errors.WithMessage(err, fmt.Sprintf("could not get content of %s", filePath))
			}
			additions[filePath] = []byte(content)
		case "D": // Deleted
			deletions = append(deletions, filePath)
		}
	}

	return git.Changes{
		Message:   commitMessage,
		Additions: additions,
		Deletions: deletions,
		OldHash:   fromHash,
	}, nil
}

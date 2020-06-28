package multigitter

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/lindell/multi-gitter/internal/domain"
	"github.com/lindell/multi-gitter/internal/git"
)

// VersionController fetches repositories
type VersionController interface {
	GetRepositories(ctx context.Context, orgName string) ([]domain.Repository, error)
	CreatePullRequest(ctx context.Context, repo domain.Repository, newPR domain.NewPullRequest) error
	GetPullRequestStatuses(ctx context.Context, orgName, branchName string) ([]domain.PullRequest, error)
	MergePullRequest(ctx context.Context, pr domain.PullRequest) error
}

// Runner conains fields to be able to do the run
type Runner struct {
	VersionController VersionController

	ScriptPath    string // Must be absolute path
	FeatureBranch string
	Token         string

	OrgName          string
	CommitMessage    string
	PullRequestTitle string
	PullRequestBody  string
	Reviewers        []string
	MaxReviewers     int // If set to zero, all reviewers will be used
	DryRun           bool
}

// Run runs a script for multiple repositories and creates PRs with the changes made
func (r Runner) Run(ctx context.Context) error {
	repos, err := r.VersionController.GetRepositories(ctx, r.OrgName)
	if err != nil {
		return err
	}

	rc := newRepoCounter()
	defer rc.printRepos()

	log.Infof("Running on %d repositories", len(repos))

	for _, repo := range repos {
		logger := log.WithField("repo", repo.FullName())
		err := r.runSingleRepo(ctx, repo)
		if exitErr, ok := err.(*exec.ExitError); ok {
			logger.Infof("Got exit code %d", exitErr.ExitCode())
			rc.exitCodeRepos[exitErr.ExitCode()] = append(rc.exitCodeRepos[exitErr.ExitCode()], repo)
			continue
		} else if err == domain.NoChangeError {
			logger.Info("No change done on the repo by the script")
			rc.noChangeRepos = append(rc.noChangeRepos, repo)
			continue
		} else if err == domain.BranchExistError {
			logger.Info("Branch already exist")
			rc.existingBranchRepos = append(rc.existingBranchRepos, repo)
			continue
		} else if err != nil {
			return err
		}

		rc.successRepos = append(rc.successRepos, repo)
	}

	return nil
}

func getReviewers(reviewers []string, maxReviewers int) []string {
	if maxReviewers == 0 || len(reviewers) <= maxReviewers {
		return reviewers
	}

	rand.Shuffle(len(reviewers), func(i, j int) { reviewers[i], reviewers[j] = reviewers[j], reviewers[i] })

	return reviewers[0:maxReviewers]
}

func (r Runner) runSingleRepo(ctx context.Context, repo domain.Repository) error {
	logger := log.WithField("repo", repo.FullName())
	logger.Info("Cloning and running script")

	tmpDir, err := ioutil.TempDir(os.TempDir(), "multi-git-changer-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	sourceController := &git.Git{
		Directory: tmpDir,
		Repo:      repo.URL,
		NewBranch: r.FeatureBranch,
		Token:     r.Token,
	}

	err = sourceController.Clone()
	if err != nil {
		return err
	}

	// Run the command that might or might not change the content of the repo
	// If the command return a non zero exit code, abort.
	cmd := exec.Command(r.ScriptPath)
	cmd.Dir = tmpDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("REPOSITORY_NAME=%s", repo.Name),
	)

	writer := newLogger()
	defer writer.Close()
	cmd.Stdout = writer
	cmd.Stderr = writer

	err = cmd.Run()
	if err != nil {
		return err
	}

	if changed, err := sourceController.Changes(); err != nil {
		return err
	} else if !changed {
		return domain.NoChangeError
	}

	branchExist, err := sourceController.BranchExist()
	if err != nil {
		return err
	} else if branchExist {
		return domain.BranchExistError
	}

	err = sourceController.Commit(r.CommitMessage)
	if err != nil {
		return err
	}

	if r.DryRun {
		logger.Info("Skipping pushing changes because of dry run")
		return nil
	}

	err = sourceController.Push()
	if err != nil {
		return err
	}

	logger.Info("Change done, creating pull request")
	err = r.VersionController.CreatePullRequest(ctx, repo, domain.NewPullRequest{
		Title:     r.PullRequestTitle,
		Body:      r.PullRequestBody,
		Head:      r.FeatureBranch,
		Base:      repo.DefaultBranch,
		Reviewers: getReviewers(r.Reviewers, r.MaxReviewers),
	})
	if err != nil {
		return err
	}

	return nil
}

func newLogger() io.WriteCloser {
	reader, writer := io.Pipe()

	// Print each line that is outputted by the script
	go func() {
		buf := bufio.NewReader(reader)
		for {
			line, err := buf.ReadString('\n')
			if line != "" {
				log.Infof("Script output: %s", line)
			}
			if err != nil {
				return
			}
		}
	}()

	return writer
}

type repoCounter struct {
	exitCodeRepos       map[int][]domain.Repository
	noChangeRepos       []domain.Repository
	successRepos        []domain.Repository
	existingBranchRepos []domain.Repository
}

func newRepoCounter() *repoCounter {
	return &repoCounter{
		exitCodeRepos: map[int][]domain.Repository{},
	}
}

func (r *repoCounter) printRepos() {
	var exitInfo string
	for exitCode := range r.exitCodeRepos {
		exitInfo += fmt.Sprintf("Repositories with exit code %d:\n", exitCode)
		for _, repo := range r.exitCodeRepos[exitCode] {
			exitInfo += fmt.Sprintf("  %s\n", repo.FullName())
		}
	}

	if len(r.noChangeRepos) > 0 {
		exitInfo += "Repositories where nothing was changed:\n"
		for _, repo := range r.noChangeRepos {
			exitInfo += fmt.Sprintf("  %s\n", repo.FullName())
		}
	}

	if len(r.existingBranchRepos) > 0 {
		exitInfo += "Repositories where the new branch already existed:\n"
		for _, repo := range r.existingBranchRepos {
			exitInfo += fmt.Sprintf("  %s\n", repo.FullName())
		}
	}

	if len(r.successRepos) > 0 {
		exitInfo += "Repositories with a successful run:\n"
		for _, repo := range r.successRepos {
			exitInfo += fmt.Sprintf("  %s\n", repo.FullName())
		}
	}

	if exitInfo != "" {
		fmt.Print(exitInfo)
	}
}

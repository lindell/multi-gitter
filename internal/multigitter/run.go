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
}

// Run runs a script for multiple repositories and creates PRs with the changes made
func (r Runner) Run(ctx context.Context) error {
	repos, err := r.VersionController.GetRepositories(ctx, r.OrgName)
	if err != nil {
		return err
	}

	log.Infof("Running on %d repositories", len(repos))

	exitCodeRepos := map[int][]domain.Repository{}
	noChangeRepos := []domain.Repository{}
	successRepos := []domain.Repository{}

	defer func() {
		var exitInfo string
		for exitCode := range exitCodeRepos {
			exitInfo += fmt.Sprintf("Repositories with exit code %d:\n", exitCode)
			for _, repo := range exitCodeRepos[exitCode] {
				exitInfo += fmt.Sprintf("  %s\n", repo.FullName())
			}
		}

		if len(noChangeRepos) > 0 {
			exitInfo += "Repositories where nothing was changed:\n"
			for _, repo := range noChangeRepos {
				exitInfo += fmt.Sprintf("  %s\n", repo.FullName())
			}
		}

		if len(successRepos) > 0 {
			exitInfo += "Repositories with a successful run:\n"
			for _, repo := range successRepos {
				exitInfo += fmt.Sprintf("  %s\n", repo.FullName())
			}
		}

		if exitInfo != "" {
			fmt.Print(exitInfo)
		}
	}()

	for _, repo := range repos {
		logger := log.WithField("repo", repo.FullName())
		logger.Info("Cloning and running script")
		err := r.runSingleRepo(repo.URL)
		if exitErr, ok := err.(*exec.ExitError); ok {
			logger.Infof("Got exit code %d", exitErr.ExitCode())
			exitCodeRepos[exitErr.ExitCode()] = append(exitCodeRepos[exitErr.ExitCode()], repo)
			continue
		} else if err == domain.NoChangeError {
			logger.Info("No change done on the repo by the script")
			noChangeRepos = append(noChangeRepos, repo)
			continue
		} else if err != nil {
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

		successRepos = append(successRepos, repo)
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

func (r Runner) runSingleRepo(url string) error {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "multi-git-changer-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	sourceController := git.Git{
		Directory: tmpDir,
		Repo:      url,
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

	writer := newLogger()
	defer writer.Close()
	cmd.Stdout = writer
	cmd.Stderr = writer

	err = cmd.Run()
	if err != nil {
		return err
	}

	err = sourceController.Commit(r.CommitMessage)
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

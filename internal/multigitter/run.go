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
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/lindell/multi-gitter/internal/domain"
	"github.com/lindell/multi-gitter/internal/git"
)

// VersionController fetches repositories
type VersionController interface {
	GetRepositories(ctx context.Context) ([]domain.Repository, error)
	CreatePullRequest(ctx context.Context, repo domain.Repository, newPR domain.NewPullRequest) error
	GetPullRequestStatuses(ctx context.Context, branchName string) ([]domain.PullRequest, error)
	MergePullRequest(ctx context.Context, pr domain.PullRequest) error
}

// Runner contains fields to be able to do the run
type Runner struct {
	VersionController VersionController

	ScriptPath    string // Must be absolute path
	Arguments     []string
	FeatureBranch string
	Token         string

	CommitMessage    string
	PullRequestTitle string
	PullRequestBody  string
	Reviewers        []string
	MaxReviewers     int // If set to zero, all reviewers will be used
	DryRun           bool
	CommitAuthor     *domain.CommitAuthor
}

// Run runs a script for multiple repositories and creates PRs with the changes made
func (r Runner) Run(ctx context.Context) error {
	repos, err := r.VersionController.GetRepositories(ctx)
	if err != nil {
		return err
	}

	rc := newRepoCounter()
	defer rc.printRepos()

	log.Infof("Running on %d repositories", len(repos))

	for _, repo := range repos {
		logger := log.WithField("repo", repo.FullName())
		err := r.runSingleRepo(ctx, repo)
		if err != nil {
			logger.Info(err)
			rc.addError(err, repo)
			continue
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
		Directory:    tmpDir,
		Repo:         repo.URL(r.Token),
		NewBranch:    r.FeatureBranch,
		CommitAuthor: r.CommitAuthor,
	}

	err = sourceController.Clone()
	if err != nil {
		return err
	}

	branchExist, err := sourceController.BranchExist()
	if err != nil {
		return err
	} else if branchExist {
		return domain.BranchExistError
	}

	err = sourceController.ChangeBranch()
	if err != nil {
		return err
	}

	// Run the command that might or might not change the content of the repo
	// If the command return a non zero exit code, abort.
	cmd := exec.Command(r.ScriptPath, r.Arguments...)
	cmd.Dir = tmpDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("REPOSITORY=%s", repo.FullName()),
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
		Base:      repo.DefaultBranch(),
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
	successRepos      []domain.Repository
	errorRepositories map[string][]domain.Repository
}

func newRepoCounter() *repoCounter {
	return &repoCounter{
		errorRepositories: map[string][]domain.Repository{},
	}
}

func (r *repoCounter) addError(err error, repo domain.Repository) {
	msg := err.Error()
	r.errorRepositories[msg] = append(r.errorRepositories[msg], repo)
}

func (r *repoCounter) printRepos() {
	var exitInfo string

	for errMsg := range r.errorRepositories {
		exitInfo += fmt.Sprintf("%s:\n", strings.ToUpper(errMsg[0:1])+errMsg[1:])
		for _, repo := range r.errorRepositories[errMsg] {
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

package multigitter

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/lindell/multi-gitter/internal/domain"
	"github.com/lindell/multi-gitter/internal/git"
	"github.com/lindell/multi-gitter/internal/multigitter/logger"
	"github.com/lindell/multi-gitter/internal/multigitter/repocounter"
)

// VersionController fetches repositories
type VersionController interface {
	GetRepositories(ctx context.Context) ([]domain.Repository, error)
	CreatePullRequest(ctx context.Context, repo domain.Repository, newPR domain.NewPullRequest) error
	GetPullRequestStatuses(ctx context.Context, branchName string) ([]domain.PullRequest, error)
	MergePullRequest(ctx context.Context, pr domain.PullRequest) error
	ClosePullRequest(ctx context.Context, pr domain.PullRequest) error
}

// Runner contains fields to be able to do the run
type Runner struct {
	VersionController VersionController

	ScriptPath    string // Must be absolute path
	Arguments     []string
	FeatureBranch string
	Token         string

	Output io.Writer

	CommitMessage    string
	PullRequestTitle string
	PullRequestBody  string
	Reviewers        []string
	MaxReviewers     int // If set to zero, all reviewers will be used
	DryRun           bool
	CommitAuthor     *domain.CommitAuthor
	BaseBranch       string // The base branch of the PR, use default branch if not set

	Concurrent int
}

var errAborted = errors.New("run was never started because of aborted execution")

// Run runs a script for multiple repositories and creates PRs with the changes made
func (r Runner) Run(ctx context.Context) error {
	repos, err := r.VersionController.GetRepositories(ctx)
	if err != nil {
		return err
	}

	rc := repocounter.NewCounter()
	defer func() {
		if info := rc.Info(); info != "" {
			fmt.Fprint(r.Output, info)
		}
	}()

	log.Infof("Running on %d repositories", len(repos))

	runInParallel(func(i int) {
		logger := log.WithField("repo", repos[i].FullName())
		err := r.runSingleRepo(ctx, repos[i])
		if err != nil {
			if err != errAborted {
				logger.Info(err)
			}
			rc.AddError(err, repos[i])
			return
		}

		rc.AddSuccess(repos[i])
	}, len(repos), r.Concurrent)

	return nil
}

func runInParallel(fun func(i int), total int, maxConcurrent int) {
	concurrentGoroutines := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	wg.Add(total)
	for i := 0; i < total; i++ {
		concurrentGoroutines <- struct{}{}
		go func(i int) {
			defer wg.Done()
			fun(i)
			<-concurrentGoroutines
		}(i)
	}
	wg.Wait()
}

func getReviewers(reviewers []string, maxReviewers int) []string {
	if maxReviewers == 0 || len(reviewers) <= maxReviewers {
		return reviewers
	}

	rand.Shuffle(len(reviewers), func(i, j int) { reviewers[i], reviewers[j] = reviewers[j], reviewers[i] })

	return reviewers[0:maxReviewers]
}

func (r Runner) runSingleRepo(ctx context.Context, repo domain.Repository) error {
	if ctx.Err() != nil {
		return errAborted
	}

	log := log.WithField("repo", repo.FullName())
	log.Info("Cloning and running script")

	tmpDir, err := ioutil.TempDir(os.TempDir(), "multi-git-changer-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	sourceController := &git.Git{
		Directory: tmpDir,
		Repo:      repo.URL(r.Token),
	}

	baseBranch := r.BaseBranch
	if baseBranch == "" {
		baseBranch = repo.DefaultBranch()
	}

	err = sourceController.Clone(baseBranch)
	if err != nil {
		return err
	}

	featureBranchExist, err := sourceController.BranchExist(r.FeatureBranch)
	if err != nil {
		return err
	} else if featureBranchExist {
		return domain.BranchExistError
	}

	err = sourceController.ChangeBranch(r.FeatureBranch)
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

	writer := logger.NewLogger(log)
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

	err = sourceController.Commit(r.CommitAuthor, r.CommitMessage)
	if err != nil {
		return err
	}

	if r.DryRun {
		log.Info("Skipping pushing changes because of dry run")
		return nil
	}

	err = sourceController.Push()
	if err != nil {
		return err
	}

	log.Info("Change done, creating pull request")
	err = r.VersionController.CreatePullRequest(ctx, repo, domain.NewPullRequest{
		Title:     r.PullRequestTitle,
		Body:      r.PullRequestBody,
		Head:      r.FeatureBranch,
		Base:      baseBranch,
		Reviewers: getReviewers(r.Reviewers, r.MaxReviewers),
	})
	if err != nil {
		return err
	}

	return nil
}

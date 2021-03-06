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
	CreatePullRequest(ctx context.Context, repo domain.Repository, newPR domain.NewPullRequest) (domain.PullRequest, error)
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

	FetchDepth      int // Limit fetching to the specified number of commits. Set to 0 for no limit
	Concurrent      int
	SkipPullRequest bool // If set, the script will run directly on the base-branch without creating any PR
}

var errAborted = errors.New("run was never started because of aborted execution")

type dryRunPullRequest struct {
	status     domain.PullRequestStatus
	Repository domain.Repository
}

func (pr dryRunPullRequest) Status() domain.PullRequestStatus {
	return pr.status
}

func (pr dryRunPullRequest) String() string {
	return fmt.Sprintf("%s #0", pr.Repository.FullName())
}

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
		pr, err := r.runSingleRepo(ctx, repos[i])
		if err != nil {
			if err != errAborted {
				logger.Info(err)
			}
			rc.AddError(err, repos[i])
			return
		}

		if pr != nil {
			rc.AddSuccessPullRequest(pr)
		} else {
			rc.AddSuccessRepositories(repos[i])
		}
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

func (r Runner) runSingleRepo(ctx context.Context, repo domain.Repository) (domain.PullRequest, error) {
	if ctx.Err() != nil {
		return nil, errAborted
	}

	log := log.WithField("repo", repo.FullName())
	log.Info("Cloning and running script")

	tmpDir, err := ioutil.TempDir(os.TempDir(), "multi-git-changer-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	sourceController := &git.Git{
		Directory:  tmpDir,
		Repo:       repo.URL(r.Token),
		FetchDepth: r.FetchDepth,
	}

	baseBranch := r.BaseBranch
	if baseBranch == "" {
		baseBranch = repo.DefaultBranch()
	}

	err = sourceController.Clone(baseBranch, r.FeatureBranch)
	if err != nil {
		return nil, err
	}

	// Change the branch to the feature branch
	if !r.SkipPullRequest {
		featureBranchExist, err := sourceController.BranchExist(r.FeatureBranch)
		if err != nil {
			return nil, errors.Wrap(err, "could not verify if branch already exist")
		} else if featureBranchExist {
			return nil, domain.BranchExistError
		}

		err = sourceController.ChangeBranch(r.FeatureBranch)
		if err != nil {
			return nil, err
		}
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
		return nil, err
	}

	if changed, err := sourceController.Changes(); err != nil {
		return nil, err
	} else if !changed {
		return nil, domain.NoChangeError
	}

	err = sourceController.Commit(r.CommitAuthor, r.CommitMessage)
	if err != nil {
		return nil, err
	}

	if r.DryRun {
		log.Info("Skipping pushing changes because of dry run")
		return dryRunPullRequest{
			Repository: repo,
		}, nil
	}

	err = sourceController.Push()
	if err != nil {
		return nil, errors.Wrap(err, "could not push changes")
	}

	if r.SkipPullRequest {
		return nil, nil
	}

	log.Info("Change done, creating pull request")
	pr, err := r.VersionController.CreatePullRequest(ctx, repo, domain.NewPullRequest{
		Title:     r.PullRequestTitle,
		Body:      r.PullRequestBody,
		Head:      r.FeatureBranch,
		Base:      baseBranch,
		Reviewers: getReviewers(r.Reviewers, r.MaxReviewers),
	})
	if err != nil {
		return nil, err
	}

	return pr, nil
}

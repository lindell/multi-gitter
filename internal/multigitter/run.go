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
	"syscall"

	"github.com/eiannone/keyboard"
	"github.com/lindell/multi-gitter/internal/git"
	"github.com/lindell/multi-gitter/internal/gittererrors"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/lindell/multi-gitter/internal/multigitter/logger"
	"github.com/lindell/multi-gitter/internal/multigitter/repocounter"
	"github.com/lindell/multi-gitter/internal/multigitter/terminal"
)

// VersionController fetches repositories
type VersionController interface {
	GetRepositories(ctx context.Context) ([]git.Repository, error)
	CreatePullRequest(ctx context.Context, repo git.Repository, prRepo git.Repository, newPR git.NewPullRequest) (git.PullRequest, error)
	GetPullRequests(ctx context.Context, branchName string) ([]git.PullRequest, error)
	MergePullRequest(ctx context.Context, pr git.PullRequest) error
	ClosePullRequest(ctx context.Context, pr git.PullRequest) error
	ForkRepository(ctx context.Context, repo git.Repository, newOwner string) (git.Repository, error)
}

// Runner contains fields to be able to do the run
type Runner struct {
	VersionController VersionController

	ScriptPath    string // Must be absolute path
	Arguments     []string
	FeatureBranch string

	Output io.Writer

	CommitMessage    string
	PullRequestTitle string
	PullRequestBody  string
	Reviewers        []string
	MaxReviewers     int // If set to zero, all reviewers will be used
	DryRun           bool
	CommitAuthor     *git.CommitAuthor
	BaseBranch       string // The base branch of the PR, use default branch if not set

	Concurrent      int
	SkipPullRequest bool // If set, the script will run directly on the base-branch without creating any PR

	Fork      bool   // If set, create a fork and make the pull request from it
	ForkOwner string // The owner of the new fork. If empty, the fork should happen on the logged in user

	Interactive bool // If set, interactive mode is activated and the user will be asked to verify every change

	CreateGit func(dir string) Git
}

var errAborted = errors.New("run was never started because of aborted execution")
var errRejected = errors.New("changes were not included since they were manually rejected")

type dryRunPullRequest struct {
	status     git.PullRequestStatus
	Repository git.Repository
}

func (pr dryRunPullRequest) Status() git.PullRequestStatus {
	return pr.status
}

func (pr dryRunPullRequest) String() string {
	return fmt.Sprintf("%s #0", pr.Repository.FullName())
}

// Run runs a script for multiple repositories and creates PRs with the changes made
func (r *Runner) Run(ctx context.Context) error {
	// Fetch all repositories that are are going to be used in the run
	repos, err := r.VersionController.GetRepositories(ctx)
	if err != nil {
		return errors.Wrap(err, "could not fetch repositories")
	}

	// Setting up a "counter" that keeps track of successful and failed runs
	rc := repocounter.NewCounter()
	defer func() {
		if info := rc.Info(); info != "" {
			fmt.Fprint(r.Output, info)
		}
	}()

	log.Infof("Running on %d repositories", len(repos))

	runInParallel(func(i int) {
		logger := log.WithField("repo", repos[i].FullName())

		defer func() {
			if r := recover(); r != nil {
				log.Error(r)
				rc.AddError(errors.New("run paniced"), repos[i])
			}
		}()

		pr, err := r.runSingleRepo(ctx, repos[i])
		if err != nil {
			if err != errAborted {
				logger.Info(err)
			}
			rc.AddError(err, repos[i])

			if log.IsLevelEnabled(log.TraceLevel) {
				if stackTrace := getStackTrace(err); stackTrace != "" {
					log.Trace(stackTrace)
				}
			}

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

func (r *Runner) runSingleRepo(ctx context.Context, repo git.Repository) (git.PullRequest, error) {
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

	sourceController := r.CreateGit(tmpDir)

	baseBranch := r.BaseBranch
	if baseBranch == "" {
		baseBranch = repo.DefaultBranch()
	}

	err = sourceController.Clone(repo.CloneURL(), baseBranch)
	if err != nil {
		return nil, err
	}

	// Change the branch to the feature branch
	if !r.SkipPullRequest {
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

	// Setup logger that transfers stdout and stderr from the run to logs
	writer := logger.NewLogger(log)
	defer writer.Close()
	cmd.Stdout = writer
	cmd.Stderr = writer

	err = cmd.Run()
	if err != nil {
		return nil, transformExecError(err)
	}

	if changed, err := sourceController.Changes(); err != nil {
		return nil, err
	} else if !changed {
		return nil, gittererrors.NoChangeError
	}

	err = sourceController.Commit(r.CommitAuthor, r.CommitMessage)
	if err != nil {
		return nil, err
	}

	if r.Interactive {
		err = r.interactive(tmpDir, repo)
		if err != nil {
			return nil, err
		}
	}

	if r.DryRun {
		log.Info("Skipping pushing changes because of dry run")
		return dryRunPullRequest{
			Repository: repo,
		}, nil
	}

	remoteName := "origin"
	var prRepo = repo
	if r.Fork {
		log.Info("Forking repository")

		prRepo, err = r.VersionController.ForkRepository(ctx, repo, r.ForkOwner)
		if err != nil {
			return nil, errors.Wrap(err, "could not fork repository")
		}

		err = sourceController.AddRemote("fork", prRepo.CloneURL())
		if err != nil {
			return nil, err
		}
		remoteName = "fork"
	}

	if !r.SkipPullRequest {
		featureBranchExist, err := sourceController.BranchExist(remoteName, r.FeatureBranch)
		if err != nil {
			return nil, errors.Wrap(err, "could not verify if branch already exist")
		} else if featureBranchExist {
			return nil, gittererrors.BranchExistError
		}
	}

	log.Info("Pushing changes to remote")
	err = sourceController.Push(remoteName)
	if err != nil {
		return nil, errors.Wrap(err, "could not push changes")
	}

	if r.SkipPullRequest {
		return nil, nil
	}

	log.Info("Creating pull request")
	pr, err := r.VersionController.CreatePullRequest(ctx, repo, prRepo, git.NewPullRequest{
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

var interactiveInfo = `(V)iew changes. (A)ccept or (R)eject`

func (r *Runner) interactive(dir string, repo git.Repository) error {
	fmt.Printf("Changes were made to %s\n", terminal.Bold(repo.FullName()))
	fmt.Println(interactiveInfo)
	for {
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			return err
		}

		if key == keyboard.KeyCtrlC {
			proc, err := os.FindProcess(os.Getpid())
			if err != nil {
				return err
			}
			_ = proc.Signal(syscall.SIGTERM)

			return errRejected
		}

		switch char {
		case 'v':
			fmt.Println("Showing changes...")
			cmd := exec.Command("git", "diff", "HEAD~1")
			cmd.Dir = dir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err != nil {
				return err
			}
			err = cmd.Run()
			if err != nil {
				return err
			}
		case 'r':
			fmt.Println("Rejected, continuing...")
			return errRejected
		case 'a':
			fmt.Println("Accepted, proceeding...")
			return nil
		}
	}
}

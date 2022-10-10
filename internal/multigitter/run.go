package multigitter

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sync"

	"github.com/lindell/multi-gitter/internal/git"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	mgerrors "github.com/lindell/multi-gitter/internal/multigitter/errors"
	"github.com/lindell/multi-gitter/internal/multigitter/logger"
	"github.com/lindell/multi-gitter/internal/multigitter/repocounter"
)

// VersionController fetches repositories
type VersionController interface {
	GetRepositories(ctx context.Context) ([]scm.Repository, error)
	CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error)
	GetPullRequests(ctx context.Context, branchName string) ([]scm.PullRequest, error)
	GetOpenPullRequest(ctx context.Context, repo scm.Repository, branchName string) (scm.PullRequest, error)
	MergePullRequest(ctx context.Context, pr scm.PullRequest) error
	ClosePullRequest(ctx context.Context, pr scm.PullRequest) error
	ForkRepository(ctx context.Context, repo scm.Repository, newOwner string) (scm.Repository, error)
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
	Assignees        []string

	Concurrent      int
	SkipPullRequest bool     // If set, the script will run directly on the base-branch without creating any PR
	SkipRepository  []string // A list of repositories that run will skip

	Fork      bool   // If set, create a fork and make the pull request from it
	ForkOwner string // The owner of the new fork. If empty, the fork should happen on the logged in user

	ConflictStrategy ConflictStrategy // Defines what will happen if a branch does already exist

	Draft bool // If set, creates Pull Requests as draft

	Interactive bool // If set, interactive mode is activated and the user will be asked to verify every change

	CreateGit func(dir string) Git

	TTY bool // If set, the progress will be displayed using TTY. If set, it's important that nothing writes to stdout/stderr

	repocounter *repocounter.Counter
}

type dryRunPullRequest struct {
	status     scm.PullRequestStatus
	Repository scm.Repository
}

func (pr dryRunPullRequest) Status() scm.PullRequestStatus {
	return pr.status
}

func (pr dryRunPullRequest) String() string {
	return fmt.Sprintf("%s #0", pr.Repository.FullName())
}

// Run runs a script for multiple repositories and creates PRs with the changes made
func (r *Runner) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Fetch all repositories that are are going to be used in the run
	repos, err := r.VersionController.GetRepositories(ctx)
	if err != nil {
		return errors.Wrap(err, "could not fetch repositories")
	}

	repos = filterRepositories(repos, r.SkipRepository)

	if len(repos) == 0 {
		log.Infof("No repositories found. Please make sure the user of the token has the correct access to the repos you want to change.")
		return nil
	}

	// Setting up a "counter" that keeps track of successful and failed runs
	r.repocounter = repocounter.NewCounter(repos)
	defer func() {
		if info := r.repocounter.Info(); info != "" {
			fmt.Fprint(r.Output, info)
		}
	}()

	log.Infof("Running on %d repositories", len(repos))

	if r.TTY {
		err := r.repocounter.OpenTTY()
		if err != nil {
			return err
		}
		defer func() { _ = r.repocounter.CloseTTY() }()
	}

	runInParallel(func(i int) {
		logger := log.WithField("repo", repos[i].FullName())

		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
				r.repocounter.SetError(errors.New("run paniced"), repos[i])
			}
		}()

		pr, err := r.runSingleRepo(ctx, cancel, repos[i])
		if err != nil {
			if err != mgerrors.ErrAborted {
				logger.Info(err)
			}
			r.repocounter.SetError(err, repos[i])

			if log.IsLevelEnabled(log.TraceLevel) {
				if stackTrace := getStackTrace(err); stackTrace != "" {
					log.Trace(stackTrace)
				}
			}

			return
		}

		if pr != nil {
			r.repocounter.AddSuccessPullRequest(repos[i], pr)
			r.repocounter.SetRepoAction(repos[i], repocounter.ActionSuccess)
		} else {
			r.repocounter.SetRepoAction(repos[i], repocounter.ActionSuccess)
		}
	}, len(repos), r.Concurrent)

	return nil
}

func filterRepositories(repos []scm.Repository, skipRepositoryNames []string) []scm.Repository {
	skipReposMap := map[string]struct{}{}
	for _, skipRepo := range skipRepositoryNames {
		skipReposMap[skipRepo] = struct{}{}
	}

	filteredRepos := make([]scm.Repository, 0, len(repos))
	for _, r := range repos {
		if _, shouldSkip := skipReposMap[r.FullName()]; !shouldSkip {
			filteredRepos = append(filteredRepos, r)
		} else {
			log.Infof("Skipping %s", r.FullName())
		}
	}
	return filteredRepos
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

func (r *Runner) runSingleRepo(ctx context.Context, cancelAll func(), repo scm.Repository) (scm.PullRequest, error) {
	if ctx.Err() != nil {
		return nil, mgerrors.ErrAborted
	}

	log := log.WithField("repo", repo.FullName())
	log.Info("Cloning and running script")
	r.repocounter.SetRepoAction(repo, repocounter.ActionClone)

	tmpDir, err := os.MkdirTemp(os.TempDir(), "multi-git-changer-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	sourceController := r.CreateGit(tmpDir)

	baseBranch := r.BaseBranch
	if baseBranch == "" {
		baseBranch = repo.DefaultBranch()
	}

	if baseBranch == r.FeatureBranch {
		return nil, errors.Errorf("both the feature branch and base branch was named %s, if you intended to push directly into the base branch, please use the `skip-pr` option", baseBranch)
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

	r.repocounter.SetRepoAction(repo, repocounter.ActionRun)

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
		return nil, mgerrors.ErrNoChange
	}

	err = sourceController.Commit(r.CommitAuthor, r.CommitMessage)
	if err != nil {
		return nil, err
	}

	if r.Interactive {
		err = r.interactive(cancelAll, tmpDir, repo)
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

	// Determine if a branch already exist and (depending on the conflict strategy) skip making changes
	featureBranchExist := false
	if !r.SkipPullRequest {
		featureBranchExist, err = sourceController.BranchExist(remoteName, r.FeatureBranch)
		if err != nil {
			return nil, errors.Wrap(err, "could not verify if branch already exist")
		} else if featureBranchExist && r.ConflictStrategy == ConflictStrategySkip {
			return nil, mgerrors.ErrBranchExist
		}
	}

	log.Info("Pushing changes to remote")
	r.repocounter.SetRepoAction(repo, repocounter.ActionPush)
	forcePush := featureBranchExist && r.ConflictStrategy == ConflictStrategyReplace
	err = sourceController.Push(remoteName, forcePush)
	if err != nil {
		return nil, errors.Wrap(err, "could not push changes")
	}

	if r.SkipPullRequest {
		return nil, nil
	}

	r.repocounter.SetRepoAction(repo, repocounter.ActionCreatePR)

	// Fetching any potentially existing pull request
	var existingPullRequest scm.PullRequest
	if featureBranchExist {
		pr, err := r.VersionController.GetOpenPullRequest(ctx, repo, r.FeatureBranch)
		if err != nil {
			return nil, err
		}
		existingPullRequest = pr
	}

	var pr scm.PullRequest
	if existingPullRequest != nil {
		log.Info("Skip creating pull requests since one is already open")
		pr = existingPullRequest
	} else {
		log.Info("Creating pull request")
		pr, err = r.VersionController.CreatePullRequest(ctx, repo, prRepo, scm.NewPullRequest{
			Title:     r.PullRequestTitle,
			Body:      r.PullRequestBody,
			Head:      r.FeatureBranch,
			Base:      baseBranch,
			Reviewers: getReviewers(r.Reviewers, r.MaxReviewers),
			Assignees: r.Assignees,
			Draft:     r.Draft,
		})
		if err != nil {
			return nil, err
		}
	}

	return pr, nil
}

func (r *Runner) interactive(cancelAll func(), dir string, repo scm.Repository) error {
	r.repocounter.QuestionLock()
	defer r.repocounter.QuestionUnlock()

	for {
		index := r.repocounter.AskQuestion(fmt.Sprintf("Changes were made to %s", repo.FullName()),
			[]repocounter.QuestionOption{
				{Text: "View changes", Shortcut: 'v'},
				{Text: "Accept", Shortcut: 'a'},
				{Text: "Reject", Shortcut: 'r'},
				{Text: "Cancel All", Shortcut: 'c'},
			})

		switch index {
		case -1: // Ctrl + C
			cancelAll()
			return mgerrors.ErrRejected
		case 0:
			cmd := exec.Command("git", "diff", "HEAD~1")
			cmd.Dir = dir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			r.repocounter.SuspendTTY()
			err := cmd.Run()
			r.repocounter.ResumeTTY()
			if err != nil {
				return err
			}
		case 1:
			return nil
		case 2:
			return mgerrors.ErrRejected
		case 3:
			cancelAll()
			return mgerrors.ErrRejected
		default:
			panic("should never happen")
		}
	}
}

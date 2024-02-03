package multigitter

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"syscall"

	"github.com/eiannone/keyboard"
	"github.com/lindell/multi-gitter/internal/git"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/lindell/multi-gitter/internal/multigitter/logger"
	"github.com/lindell/multi-gitter/internal/multigitter/repocounter"
	"github.com/lindell/multi-gitter/internal/multigitter/terminal"
)

// VersionController fetches repositories
type VersionController interface {
	GetRepositories(ctx context.Context) ([]scm.Repository, error)
	CreatePullRequest(ctx context.Context, repo scm.Repository, prRepo scm.Repository, newPR scm.NewPullRequest) (scm.PullRequest, error)
	UpdatePullRequest(ctx context.Context, repo scm.Repository, pullReq scm.PullRequest, updatedPR scm.NewPullRequest) (scm.PullRequest, error)
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
	TeamReviewers    []string
	MaxReviewers     int // If set to zero, all reviewers will be use
	MaxTeamReviewers int // If set to zero, all team-reviewers will be used
	DryRun           bool
	CommitAuthor     *git.CommitAuthor
	BaseBranch       string // The base branch of the PR, use default branch if not set
	Assignees        []string

	Concurrent             int
	SkipPullRequest        bool     // If set, the script will run directly on the base-branch without creating any PR
	SkipRepository         []string // A list of repositories that run will skip
	RegExIncludeRepository *regexp.Regexp
	RegExExcludeRepository *regexp.Regexp

	Fork      bool   // If set, create a fork and make the pull request from it
	ForkOwner string // The owner of the new fork. If empty, the fork should happen on the logged in user

	ConflictStrategy ConflictStrategy // Defines what will happen if a branch already exists

	Draft bool // If set, creates Pull Requests as draft

	Labels []string // Labels to be added to the pull request

	Interactive bool // If set, interactive mode is activated and the user will be asked to verify every change

	CreateGit func(dir string) Git
}

var (
	errAborted     = errors.New("run was never started because of aborted execution")
	errRejected    = errors.New("changes were not included since they were manually rejected")
	errNoChange    = errors.New("no data was changed")
	errBranchExist = errors.New("the new branch already exists")
)

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
	// Fetch all repositories that are are going to be used in the run
	repos, err := r.VersionController.GetRepositories(ctx)
	if err != nil {
		return errors.Wrap(err, "could not fetch repositories")
	}

	repos = filterRepositories(repos, r.SkipRepository, r.RegExIncludeRepository, r.RegExExcludeRepository)

	if len(repos) == 0 {
		log.Infof("No repositories found. Please make sure the user of the token has the correct access to the repos you want to change.")
		return nil
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
				rc.AddError(errors.New("run panicked"), repos[i], nil)
			}
		}()

		pr, err := r.runSingleRepo(ctx, repos[i])
		if err != nil {
			if err != errAborted {
				logger.Info(err)
			}
			rc.AddError(err, repos[i], pr)

			if log.IsLevelEnabled(log.TraceLevel) {
				if stackTrace := getStackTrace(err); stackTrace != "" {
					log.Trace(stackTrace)
				}
			}

			return
		}

		if pr != nil {
			rc.AddSuccessPullRequest(repos[i], pr)
		} else {
			rc.AddSuccessRepositories(repos[i])
		}
	}, len(repos), r.Concurrent)

	return nil
}

// Determines if Repository should be excluded based on provided Regular Expression
func excludeRepositoryFilter(repoName string, regExp *regexp.Regexp) bool {
	if regExp == nil {
		return false
	}
	return regExp.MatchString(repoName)
}

// Determines if Repository should be included based on provided Regular Expression
func matchesRepositoryFilter(repoName string, regExp *regexp.Regexp) bool {
	if regExp == nil {
		return true
	}
	return regExp.MatchString(repoName)
}

func filterRepositories(repos []scm.Repository, skipRepositoryNames []string, regExIncludeRepository *regexp.Regexp,
	regExExcludeRepository *regexp.Regexp,
) []scm.Repository {
	skipReposMap := map[string]struct{}{}
	for _, skipRepo := range skipRepositoryNames {
		skipReposMap[skipRepo] = struct{}{}
	}

	filteredRepos := make([]scm.Repository, 0, len(repos))
	for _, r := range repos {
		if _, shouldSkip := skipReposMap[r.FullName()]; shouldSkip {
			log.Infof("Skipping %s since it is in exclusion list", r.FullName())
		} else if !matchesRepositoryFilter(r.FullName(), regExIncludeRepository) {
			log.Infof("Skipping %s since it does not match the inclusion regexp", r.FullName())
		} else if excludeRepositoryFilter(r.FullName(), regExExcludeRepository) {
			log.Infof("Skipping %s since it match the exclusion regexp", r.FullName())
		} else {
			filteredRepos = append(filteredRepos, r)
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

func (r *Runner) runSingleRepo(ctx context.Context, repo scm.Repository) (scm.PullRequest, error) {
	if ctx.Err() != nil {
		return nil, errAborted
	}

	log := log.WithField("repo", repo.FullName())
	log.Info("Cloning and running script")

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

	err = sourceController.Clone(ctx, repo.CloneURL(), baseBranch)
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

	cmd := prepareScriptCommand(ctx, repo, tmpDir, r.ScriptPath, r.Arguments)
	if r.DryRun {
		cmd.Env = append(cmd.Env, "DRY_RUN=true")
	}

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
		return nil, errNoChange
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
	prRepo := repo
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

	// Determine if a branch already exists and (depending on the conflict strategy) skip making changes
	featureBranchExist := false
	if !r.SkipPullRequest {
		featureBranchExist, err = sourceController.BranchExist(remoteName, r.FeatureBranch)
		if err != nil {
			return nil, errors.Wrap(err, "could not verify if branch already exists")
		} else if featureBranchExist && r.ConflictStrategy == ConflictStrategySkip {
			pr, err := r.ensurePullRequestExists(ctx, log, repo, prRepo, baseBranch, featureBranchExist)
			if err != nil {
				return nil, err
			}

			return pr, errBranchExist
		}
	}

	log.Info("Pushing changes to remote")
	forcePush := featureBranchExist && r.ConflictStrategy == ConflictStrategyReplace
	err = sourceController.Push(ctx, remoteName, forcePush)
	if err != nil {
		return nil, errors.Wrap(err, "could not push changes")
	}

	return r.ensurePullRequestExists(ctx, log, repo, prRepo, baseBranch, featureBranchExist)
}

func (r *Runner) ensurePullRequestExists(ctx context.Context, log log.FieldLogger, repo scm.Repository, prRepo scm.Repository, baseBranch string, featureBranchExist bool) (scm.PullRequest, error) {
	if r.SkipPullRequest {
		return nil, nil
	}

	// Fetching any potentially existing pull request
	var existingPullRequest scm.PullRequest
	if featureBranchExist {
		pr, err := r.VersionController.GetOpenPullRequest(ctx, repo, r.FeatureBranch)
		if err != nil {
			return nil, err
		}
		existingPullRequest = pr
	}

	if existingPullRequest != nil {
		if r.ConflictStrategy == ConflictStrategyReplace {
			log.Info("Updating pull request since one is already open")
			return r.VersionController.UpdatePullRequest(ctx, repo, existingPullRequest, scm.NewPullRequest{
				Title:         r.PullRequestTitle,
				Body:          r.PullRequestBody,
				Head:          r.FeatureBranch,
				Base:          baseBranch,
				Reviewers:     getReviewers(r.Reviewers, r.MaxReviewers),
				TeamReviewers: getReviewers(r.TeamReviewers, r.MaxTeamReviewers),
				Assignees:     r.Assignees,
				Draft:         r.Draft,
				Labels:        r.Labels,
			})
		}
		log.Info("Skip creating pull requests since one is already open")
		return existingPullRequest, nil
	}

	log.Info("Creating pull request")
	return r.VersionController.CreatePullRequest(ctx, repo, prRepo, scm.NewPullRequest{
		Title:         r.PullRequestTitle,
		Body:          r.PullRequestBody,
		Head:          r.FeatureBranch,
		Base:          baseBranch,
		Reviewers:     getReviewers(r.Reviewers, r.MaxReviewers),
		TeamReviewers: getReviewers(r.TeamReviewers, r.MaxTeamReviewers),
		Assignees:     r.Assignees,
		Draft:         r.Draft,
		Labels:        r.Labels,
	})
}

var interactiveInfo = `(V)iew changes. (A)ccept or (R)eject`

func (r *Runner) interactive(dir string, repo scm.Repository) error {
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

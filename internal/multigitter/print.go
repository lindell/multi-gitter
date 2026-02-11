package multigitter

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lindell/multi-gitter/internal/multigitter/repocounter"
	"github.com/lindell/multi-gitter/internal/scm"
	log "github.com/sirupsen/logrus"
)

// Printer contains fields to be able to do the print command
type Printer struct {
	VersionController VersionController

	ScriptPath string // Must be absolute path
	Arguments  []string

	Stdout io.Writer
	Stderr io.Writer

	// RepoFilters contains repository filtering options
	RepoFilters RepoFilters

	Concurrent int
	CloneDir   string

	Keep bool // If set, skip deletion of cloned repos and reuse them if already present

	CreateGit func(dir string) Git
}

// Print runs a script for multiple repositories and print the output of each run
func (r Printer) Print(ctx context.Context) error {
	repos, err := r.VersionController.GetRepositories(ctx)
	if err != nil {
		return err
	}

	repos = filterRepositories(repos, r.RepoFilters)

	if len(repos) == 0 {
		log.Infof("No repositories found. Please make sure the user of the token has the correct access to the repos you want print to run on.")
		return nil
	}

	rc := repocounter.NewCounter()
	defer func() {
		if info := rc.Info(); info != "" {
			fmt.Fprint(log.StandardLogger().Out, info)
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
			rc.AddError(err, repos[i], nil)
			return
		}

		rc.AddSuccessRepositories(repos[i])
	}, len(repos), r.Concurrent)

	return nil
}

func (r Printer) runSingleRepo(ctx context.Context, repo scm.Repository) error {
	if ctx.Err() != nil {
		return errAborted
	}

	log := log.WithField("repo", repo.FullName())
	log.Info("Cloning and running script")

	var tmpDir string
	var err error
	if r.Keep {
		tmpDir, err = keepDir(r.CloneDir, repo.FullName())
	} else {
		tmpDir, err = createTempDir(r.CloneDir)
	}
	if err != nil {
		return err
	}

	if !r.Keep {
		defer os.RemoveAll(tmpDir)
	}

	sourceController := r.CreateGit(tmpDir)

	baseBranch := repo.DefaultBranch()

	// If keep mode is enabled and the directory already exists, reuse it with a hard reset
	if r.Keep {
		if _, statErr := os.Stat(filepath.Join(tmpDir, ".git")); statErr == nil {
			log.Info("Reusing existing clone, resetting to base branch")
			err = sourceController.FetchAndResetToDefault(ctx, baseBranch)
			if err != nil {
				// If reset fails, remove and re-clone
				log.WithError(err).Info("Reset failed, re-cloning")
				os.RemoveAll(tmpDir)
				err = os.MkdirAll(tmpDir, 0755)
				if err != nil {
					return err
				}
				err = sourceController.Clone(ctx, repo.CloneURL(), baseBranch)
				if err != nil {
					return err
				}
			}
		} else {
			err = os.MkdirAll(tmpDir, 0755)
			if err != nil {
				return err
			}
			err = sourceController.Clone(ctx, repo.CloneURL(), baseBranch)
			if err != nil {
				return err
			}
		}
	} else {
		err = sourceController.Clone(ctx, repo.CloneURL(), baseBranch)
		if err != nil {
			return err
		}
	}

	cmd := prepareScriptCommand(ctx, repo, tmpDir, r.ScriptPath, r.Arguments)

	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr

	err = cmd.Run()
	if err != nil {
		return transformExecError(err)
	}

	return nil
}

package multigitter

import (
	"context"
	"fmt"
	"github.com/lindell/multi-gitter/internal/multigitter/repocounter"
	"github.com/lindell/multi-gitter/internal/scm"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

// Printer contains fields to be able to do the print command
type Printer struct {
	VersionController VersionController

	ScriptPath string // Must be absolute path
	Arguments  []string

	Stdout io.Writer
	Stderr io.Writer

	Concurrent int
	CloneDir   string

	CreateGit func(dir string) Git
}

// Print runs a script for multiple repositories and print the output of each run
func (r Printer) Print(ctx context.Context) error {
	repos, err := r.VersionController.GetRepositories(ctx)
	if err != nil {
		return err
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
			rc.AddError(err, repos[i])
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
	tmpDir, err := createTempDir(r.CloneDir)

	defer os.RemoveAll(tmpDir)
	if err != nil {
		return err
	}

	sourceController := r.CreateGit(tmpDir)

	err = sourceController.Clone(ctx, repo.CloneURL(), repo.DefaultBranch())
	if err != nil {
		return err
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

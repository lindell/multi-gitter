package multigitter

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	mgerrors "github.com/lindell/multi-gitter/internal/multigitter/errors"
	"github.com/lindell/multi-gitter/internal/multigitter/repocounter"
	"github.com/lindell/multi-gitter/internal/scm"
)

// Printer contains fields to be able to do the print command
type Printer struct {
	VersionController VersionController

	ScriptPath string // Must be absolute path
	Arguments  []string

	Stdout io.Writer
	Stderr io.Writer

	Concurrent int

	CreateGit func(dir string) Git
}

// Print runs a script for multiple repositories and print the output of each run
func (r Printer) Print(ctx context.Context) error {
	repos, err := r.VersionController.GetRepositories(ctx)
	if err != nil {
		return err
	}

	rc := repocounter.NewCounter(repos)
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
			if err != mgerrors.ErrAborted {
				logger.Info(err)
			}
			rc.SetError(err, repos[i])
			return
		}

		rc.SetRepoAction(repos[i], repocounter.ActionSuccess)
	}, len(repos), r.Concurrent)

	return nil
}

func (r Printer) runSingleRepo(ctx context.Context, repo scm.Repository) error {
	if ctx.Err() != nil {
		return mgerrors.ErrAborted
	}

	log := log.WithField("repo", repo.FullName())
	log.Info("Cloning and running script")

	tmpDir, err := os.MkdirTemp(os.TempDir(), "multi-git-changer-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	sourceController := r.CreateGit(tmpDir)

	err = sourceController.Clone(repo.CloneURL(), repo.DefaultBranch())
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

	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr

	err = cmd.Run()
	if err != nil {
		return transformExecError(err)
	}

	return nil
}

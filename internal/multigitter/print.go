package multigitter

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/lindell/multi-gitter/internal/git"
	log "github.com/sirupsen/logrus"

	"github.com/lindell/multi-gitter/internal/multigitter/repocounter"
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

func (r Printer) runSingleRepo(ctx context.Context, repo git.Repository) error {
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

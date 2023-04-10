package multigitter

import (
	"context"
	"fmt"
	"github.com/lindell/multi-gitter/internal/scm"
	"os"
	"os/exec"
)

func prepareScriptCommand(
	ctx context.Context,
	repo scm.Repository,
	workDir string,
	scriptPath string,
	arguments []string,
) (cmd *exec.Cmd) {
	// Run the command that might or might not change the content of the repo
	// If the command return a non-zero exit code, abort.
	cmd = exec.CommandContext(ctx, scriptPath, arguments...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("REPOSITORY=%s", repo.FullName()),
	)
	return cmd
}

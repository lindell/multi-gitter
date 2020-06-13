package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

// this file contains a helper wrapper of exec.Cmd to get better error messages

func (g Git) command(name string, args ...string) *cmd {
	c := &cmd{
		Cmd: exec.Command(name, args...),
	}
	c.Dir = g.Directory
	c.Stderr = &bytes.Buffer{}
	return c
}

type cmd struct {
	*exec.Cmd
}

func (c *cmd) Run() error {
	err := c.Cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("command exit with error code %d:\n%s", exitErr.ExitCode(), c.Stderr)
	}
	return nil
}

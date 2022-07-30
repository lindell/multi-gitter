package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	"github.com/lindell/multi-gitter/internal/multigitter"
	"github.com/spf13/cobra"
)

const printHelp = `
This command will clone down multiple repositories. For each of those repositories, the script will be run in the context of that repository. The output of each script run in each repo will be printed, by default to stdout and stderr, but it can be configured to files as well.

The environment variable REPOSITORY will be set to the name of the repository currently being executed by the script.
`

// PrintCmd is the main command that runs a script for multiple repositories and print the output of each run
func PrintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "print [script path]",
		Short:   "Clones multiple repositories, run a script in that directory, and prints the output of each run.",
		Long:    printHelp,
		Args:    cobra.ExactArgs(1),
		PreRunE: logFlagInit,
		RunE:    print,
	}

	cmd.Flags().IntP("concurrent", "C", 1, "The maximum number of concurrent runs.")
	cmd.Flags().StringP("error-output", "E", "-", `The file that the output of the script should be outputted to. "-" means stderr.`)
	configureGit(cmd)
	configurePlatform(cmd)
	configureLogging(cmd, "")
	configureConfig(cmd)
	cmd.Flags().AddFlagSet(outputFlag())

	return cmd
}

func print(cmd *cobra.Command, args []string) error {
	flag := cmd.Flags()

	concurrent, _ := flag.GetInt("concurrent")
	strOutput, _ := flag.GetString("output")
	strErrOutput, _ := flag.GetString("error-output")

	_ = flag.Set("readOnly", "true")

	if concurrent < 1 {
		return errors.New("concurrent runs can't be less than one")
	}

	output, err := fileOutput(strOutput, os.Stdout)
	if err != nil {
		return err
	}

	errOutput, err := fileOutput(strErrOutput, os.Stderr)
	if err != nil {
		return err
	}

	vc, err := getVersionController(flag, true)
	if err != nil {
		return err
	}

	gitCreator, err := getGitCreator(flag)
	if err != nil {
		return err
	}

	executablePath, arguments, err := parseCommand(flag.Arg(0))
	if err != nil {
		return err
	}

	// Set up signal listening to cancel the context and let started runs finish gracefully
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Finishing up ongoing runs. Press CTRL+C again to abort now.")
		cancel()
		<-c
		os.Exit(1)
	}()

	printer := multigitter.Printer{
		ScriptPath: executablePath,
		Arguments:  arguments,

		VersionController: vc,

		Stdout: output,
		Stderr: errOutput,

		Concurrent: concurrent,

		CreateGit: gitCreator,
	}

	err = printer.Print(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return nil
}

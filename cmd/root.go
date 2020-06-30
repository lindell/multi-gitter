package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/lindell/multi-gitter/internal/github"
	"github.com/lindell/multi-gitter/internal/multigitter"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// RootCmd is the root command containing all subcommands
var RootCmd = &cobra.Command{
	Use:               "multi-gitter",
	Short:             "Multi gitter is a tool for making changes into multiple git repositories",
	PersistentPreRunE: persistentPreRun,
}

func init() {
	RootCmd.PersistentFlags().StringP("gh-base-url", "g", "", "Base URL of the (v3) GitHub API, needs to be changed if GitHub enterprise is used.")
	RootCmd.PersistentFlags().StringP("token", "T", "", "The GitHub personal access token. Can also be set using the GITHUB_TOKEN environment variable.")
	RootCmd.PersistentFlags().StringP("log-level", "L", "info", "The level of logging that should be made. Available values: trace, debug, info, error")

	RootCmd.AddCommand(RunCmd)
	RootCmd.AddCommand(StatusCmd)
	RootCmd.AddCommand(MergeCmd)
	RootCmd.AddCommand(VersionCmd)

	rand.Seed(time.Now().UTC().UnixNano())
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	strLevel, _ := cmd.Flags().GetString("log-level")
	logLevel, err := log.ParseLevel(strLevel)
	if err != nil {
		return fmt.Errorf("invalid log-level: %s", strLevel)
	}
	log.SetLevel(logLevel)

	return nil
}

func getVersionController(flag *flag.FlagSet) (multigitter.VersionController, error) {
	ghBaseURL, _ := flag.GetString("gh-base-url")

	token, err := getToken(flag)
	if err != nil {
		return nil, err
	}

	vc, err := github.New(token, ghBaseURL)
	if err != nil {
		return nil, err
	}

	return vc, nil
}

func getToken(flag *flag.FlagSet) (string, error) {
	token, _ := flag.GetString("token")

	if token == "" {
		if ght := os.Getenv("GITHUB_TOKEN"); ght != "" {
			token = ght
		}
	}

	if token == "" {
		return "", errors.New("either the --token flag or the GITHUB_TOKEN environment variable has to be set")
	}

	return token, nil
}

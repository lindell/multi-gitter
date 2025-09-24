package cmd

import (
	"io"
	"os"

	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

func outputFlag() *flag.FlagSet {
	flags := flag.NewFlagSet("output", flag.ExitOnError)

	flags.StringP("output", "o", "-", `The file that the output of the script should be outputted to. "-" means stdout.`)

	return flags
}

func configureMergeType(cmd *cobra.Command, includeAutoMergeText bool) {
	description := "The type of merge that should be done (GitHub/Gitea). Multiple types can be used as backup strategies if the first one is not allowed."
	if includeAutoMergeText {
		description += " The first type is used for auto-merge."
	}
	
	cmd.Flags().StringSliceP("merge-type", "", []string{"merge", "squash", "rebase"}, description)
	_ = cmd.RegisterFlagCompletionFunc("merge-type", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"merge", "squash", "rebase"}, cobra.ShellCompDirectiveNoFileComp
	})
}

func getToken(flag *flag.FlagSet) (string, error) {
	if OverrideVersionController != nil {
		return "", nil
	}

	token, _ := flag.GetString("token")

	if token == "" {
		if ght := os.Getenv("GITHUB_TOKEN"); ght != "" {
			token = ght
		} else if ght := os.Getenv("GITLAB_TOKEN"); ght != "" {
			token = ght
		} else if ght := os.Getenv("GITEA_TOKEN"); ght != "" {
			token = ght
		} else if ght := os.Getenv("BITBUCKET_SERVER_TOKEN"); ght != "" {
			token = ght
		} else if ght := os.Getenv("BITBUCKET_CLOUD_APP_PASSWORD"); ght != "" {
			token = ght
		} else if ght := os.Getenv("BITBUCKET_CLOUD_WORKSPACE_TOKEN"); ght != "" {
			token = ght
		} else if ght := os.Getenv("GERRIT_TOKEN"); ght != "" {
			token = ght
		}
	}

	if token == "" {
		return "", errors.New("either the --token flag or the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN/BITBUCKET_CLOUD_APP_PASSWORD/BITBUCKET_CLOUD_WORKSPACE_TOKEN/GERRIT_TOKEN environment variable has to be set")
	}

	return token, nil
}

func getMergeTypes(flag *flag.FlagSet) ([]scm.MergeType, error) {
	mergeTypeStrs, _ := flag.GetStringSlice("merge-type") // Only used for the merge command

	// Convert all defined merge types (if any)
	var err error
	mergeTypes := make([]scm.MergeType, len(mergeTypeStrs))
	for i, mt := range mergeTypeStrs {
		mergeTypes[i], err = scm.ParseMergeType(mt)
		if err != nil {
			return nil, err
		}
	}

	return mergeTypes, nil
}

// nopWriter is a writer that does nothing
type nopWriter struct{}

func (nw nopWriter) Write(bb []byte) (int, error) {
	return len(bb), nil
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func fileOutput(value string, std io.Writer) (io.WriteCloser, error) {
	if value != "-" {
		file, err := os.Create(value)
		if err != nil {
			return nil, errors.Wrapf(err, "could not open file %s", value)
		}
		return file, nil
	}
	return nopCloser{std}, nil
}

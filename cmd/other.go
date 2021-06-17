package cmd

import (
	"io"
	"os"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"

	"github.com/lindell/multi-gitter/internal/domain"
)

func outputFlag() *flag.FlagSet {
	flags := flag.NewFlagSet("output", flag.ExitOnError)

	flags.StringP("output", "o", "-", `The file that the output of the script should be outputted to. "-" means stdout`)

	return flags
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
		}
	}

	if token == "" {
		return "", errors.New("either the --token flag or the GITHUB_TOKEN environment variable has to be set")
	}

	return token, nil
}

func getMergeTypes(flag *flag.FlagSet) ([]domain.MergeType, error) {
	mergeTypeStrs, _ := flag.GetStringSlice("merge-type") // Only used for the merge command

	// Convert all defined merge types (if any)
	var err error
	mergeTypes := make([]domain.MergeType, len(mergeTypeStrs))
	for i, mt := range mergeTypeStrs {
		mergeTypes[i], err = domain.ParseMergeType(mt)
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

package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	internallog "github.com/lindell/multi-gitter/internal/log"
	"github.com/lindell/multi-gitter/internal/multigitter/terminal"
)

func configureLogging(cmd *cobra.Command, logFile string) {
	flags := cmd.Flags()

	flags.StringP("log-level", "L", "info", "The level of logging that should be made. Available values: trace, debug, info, error.")
	_ = cmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"trace", "debug", "info", "error"}, cobra.ShellCompDirectiveDefault
	})

	flags.StringP("log-format", "", "text", `The formatting of the logs. Available values: text, json, json-pretty.`)
	_ = cmd.RegisterFlagCompletionFunc("log-format", func(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"text", "json", "json-pretty"}, cobra.ShellCompDirectiveDefault
	})

	flags.StringP("log-file", "", logFile, `The file where all logs should be printed to. "-" means stdout.`)

	flags.BoolP("plain-output", "", false, `Don't use any terminal formatting when printing the output.`)
}

func logFlagInit(cmd *cobra.Command, _ []string) error {
	// Parse and set log level
	strLevel, _ := cmd.Flags().GetString("log-level")
	logLevel, err := log.ParseLevel(strLevel)
	if err != nil {
		return fmt.Errorf("invalid log-level: %s", strLevel)
	}
	log.SetLevel(logLevel)

	// Set how custom terminal formatting is handled
	plainOutput, _ := cmd.Flags().GetBool("plain-output")
	terminal.DefaultPrinter.Plain = plainOutput

	// Parse and set the log format
	strFormat, _ := cmd.Flags().GetString("log-format")

	var formatter log.Formatter
	switch strFormat {
	case "text":
		formatter = &log.TextFormatter{
			DisableColors: plainOutput,
		}
	case "json":
		formatter = &log.JSONFormatter{}
	case "json-pretty":
		if plainOutput {
			return errors.New("can't use json-pretty logs with with plain-output")
		}
		formatter = &log.JSONFormatter{
			PrettyPrint: true,
		}
	default:
		return fmt.Errorf(`unknown log-format "%s"`, strFormat)
	}

	// Make sure sensitive data is censored before logging them
	var censorItems []internallog.CensorItem
	if token, err := getToken(cmd.Flags()); err == nil && token != "" {
		censorItems = append(censorItems, internallog.CensorItem{
			Sensitive:   token,
			Replacement: "<TOKEN>",
		})
	}

	log.SetFormatter(internallog.NewCensorFormatter(formatter, censorItems...))

	// Set the output (file)
	strFile, _ := cmd.Flags().GetString("log-file")
	if strFile == "" {
		log.SetOutput(nopWriter{})
	} else if strFile != "-" {
		file, err := os.Create(strFile)
		if err != nil {
			return errors.Wrapf(err, "could not open log-file %s", strFile)
		}
		log.SetOutput(file)
	}

	return nil
}

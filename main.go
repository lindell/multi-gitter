package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lindell/multi-gitter/cmd"
)

var version = "development"
var date = "now"
var commit = "unknown"

func main() {
	cmd.Version = version
	cmd.BuildDate, _ = time.ParseInLocation(time.RFC3339, date, time.UTC)
	cmd.Commit = commit
	if err := cmd.RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

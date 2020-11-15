package main

import (
	"fmt"
	"os"

	"github.com/lindell/multi-gitter/cmd"
)

var version = "development"

func main() {
	cmd.Version = version
	if err := cmd.RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

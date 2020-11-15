package main

import (
	"log"
	"os"

	"github.com/lindell/multi-gitter/cmd"
	"github.com/spf13/cobra/doc"
)

const genDir = "./tmp-docs"

func main() {
	os.RemoveAll(genDir)
	err := os.MkdirAll(genDir, os.ModeDir|0700)
	if err != nil {
		log.Fatal(err)
	}

	rootCmd := cmd.RootCmd()
	rootCmd.DisableAutoGenTag = true
	err = doc.GenMarkdownTree(rootCmd, genDir)
	if err != nil {
		log.Fatal(err)
	}
}

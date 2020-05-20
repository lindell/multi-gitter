package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"github.com/lindell/multi-gitter/cmd"
)

const templatePath = "./docs/README.template.md"
const resultingPath = "./README.md"

type templateData struct {
	Usage string
}

func main() {
	usageBuf := &bytes.Buffer{}
	cmd.RunCmd.SetOutput(usageBuf)
	err := cmd.RunCmd.Usage()
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal(err)
	}

	tmplBuf := &bytes.Buffer{}
	err = tmpl.Execute(tmplBuf, templateData{
		Usage: strings.TrimSpace(usageBuf.String()),
	})
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(resultingPath, tmplBuf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

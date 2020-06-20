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
	MainUsage string
	Commands  []command
}

type command struct {
	Name        string
	Description string
	Usage       string
}

func main() {
	data := templateData{}

	data.MainUsage = strings.TrimSpace(cmd.RootCmd.UsageString())

	for _, c := range cmd.RootCmd.Commands() {
		data.Commands = append(data.Commands, command{
			Name:        c.Name(),
			Description: c.Long,
			Usage:       strings.TrimSpace(c.UsageString()),
		})
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal(err)
	}

	tmplBuf := &bytes.Buffer{}
	err = tmpl.Execute(tmplBuf, data)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(resultingPath, tmplBuf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

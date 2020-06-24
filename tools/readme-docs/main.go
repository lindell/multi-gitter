package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/lindell/multi-gitter/cmd"
)

const templatePath = "./docs/README.template.md"
const resultingPath = "./README.md"

type templateData struct {
	MainUsage string
	Commands  []command
}

type command struct {
	ImageIcon string
	Name      string
	Long      string
	Short     string
	Usage     string
}

func main() {
	data := templateData{}

	data.MainUsage = strings.TrimSpace(cmd.RootCmd.UsageString())

	cmds := []struct {
		imgIcon string
		cmd     *cobra.Command
	}{
		{
			imgIcon: "docs/img/fa/rabbit-fast.svg",
			cmd:     cmd.RunCmd,
		},
		{
			imgIcon: "docs/img/fa/code-merge.svg",
			cmd:     cmd.MergeCmd,
		},
		{
			imgIcon: "docs/img/fa/tasks.svg",
			cmd:     cmd.StatusCmd,
		},
	}
	for _, c := range cmds {
		data.Commands = append(data.Commands, command{
			Name:      c.cmd.Name(),
			ImageIcon: c.imgIcon,
			Long:      c.cmd.Long,
			Short:     c.cmd.Short,
			Usage:     strings.TrimSpace(c.cmd.UsageString()),
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

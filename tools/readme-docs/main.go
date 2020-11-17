package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/lindell/multi-gitter/cmd"
)

const templatePath = "./docs/README.template.md"
const resultingPath = "./README.md"

type templateData struct {
	MainUsage         string
	Commands          []command
	ExampleCategories []exampleCategory
}

type command struct {
	ImageIcon string
	Name      string
	Long      string
	Short     string
	Usage     string
}

type exampleCategory struct {
	Name     string
	Examples []example
}

type example struct {
	Title string
	Body  string
}

func main() {
	data := templateData{}

	// Main usage
	data.MainUsage = strings.TrimSpace(cmd.RootCmd().UsageString())

	subCommands := cmd.RootCmd().Commands()

	// All commands
	cmds := []struct {
		imgIcon string
		cmd     *cobra.Command
	}{
		{
			imgIcon: "docs/img/fa/rabbit-fast.svg",
			cmd:     commandByName(subCommands, "run"),
		},
		{
			imgIcon: "docs/img/fa/code-merge.svg",
			cmd:     commandByName(subCommands, "merge"),
		},
		{
			imgIcon: "docs/img/fa/tasks.svg",
			cmd:     commandByName(subCommands, "status"),
		},
		{
			imgIcon: "docs/img/fa/times-hexagon.svg",
			cmd:     commandByName(subCommands, "close"),
		},
		{
			imgIcon: "docs/img/fa/print.svg",
			cmd:     commandByName(subCommands, "print"),
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

	var err error
	data.ExampleCategories, err = readExamples()
	if err != nil {
		log.Fatal(err)
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

func commandByName(cmds []*cobra.Command, name string) *cobra.Command {
	for _, command := range cmds {
		if command.Name() == name {
			return command
		}
	}
	panic(fmt.Sprintf(`could not find command "%s"`, name))
}

var titleRegex = regexp.MustCompile("# ?Title: ([^\n]+)[\n\r]+")

func readExamples() ([]exampleCategory, error) {
	categories := []exampleCategory{}

	examplesDir := "./examples"
	files, err := ioutil.ReadDir(examplesDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		var examples []example
		categoryDir := path.Join(examplesDir, f.Name())
		exampleFiles, err := ioutil.ReadDir(categoryDir)
		if err != nil {
			return nil, err
		}
		for _, e := range exampleFiles {
			b, err := ioutil.ReadFile(path.Join(categoryDir, e.Name()))
			if err != nil {
				return nil, err
			}

			matches := titleRegex.FindSubmatch(b)
			if matches == nil {
				return nil, errors.New("could not find title")
			}

			examples = append(examples, example{
				Title: string(matches[1]),
				Body:  strings.TrimSpace(string(titleRegex.ReplaceAll(b, nil))),
			})
		}

		category := &exampleCategory{
			Name:     f.Name(),
			Examples: examples,
		}
		categories = append(categories, *category)
	}

	return categories, nil
}

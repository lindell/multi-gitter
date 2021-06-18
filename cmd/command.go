package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

func parseCommand(command string) (executablePath string, arguments []string, err error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", nil, errors.New("could not get the working directory")
	}

	parsedCommand, err := parseCommandLine(command)
	if err != nil {
		return "", nil, errors.Errorf("could not parse command: %s", err)
	}
	executablePath, err = exec.LookPath(parsedCommand[0])
	if err != nil {
		if _, err := os.Stat(parsedCommand[0]); os.IsNotExist(err) {
			return "", nil, errors.Errorf("could not find executable %s", parsedCommand[0])
		}
		return "", nil, errors.Errorf("could not find executable %s, does it have executable privileges?", parsedCommand[0])
	}
	// Executable needs to be defined with an absolute path since it will be run within the context of repositories
	if !filepath.IsAbs(executablePath) {
		executablePath = filepath.Join(workingDir, executablePath)
	}

	return executablePath, parsedCommand[1:], nil
}

// https://stackoverflow.com/a/46973603
func parseCommandLine(command string) ([]string, error) {
	type state int

	const (
		stateStart state = iota
		stateQuotes
		stateArg
	)

	var args []string
	currentState := stateStart
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if currentState == stateQuotes {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				currentState = stateStart
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			currentState = stateQuotes
			quote = string(c)
			continue
		}

		if currentState == stateArg {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				currentState = stateStart
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			currentState = stateArg
			current += string(c)
		}
	}

	if currentState == stateQuotes {
		return []string{}, fmt.Errorf("unclosed quote in command line: %s", command)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

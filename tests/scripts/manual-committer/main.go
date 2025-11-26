package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const fileName = "test.txt"

func main() {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	data = bytes.ReplaceAll(data, []byte("apple"), []byte("banana"))

	err = os.WriteFile(fileName, data, 0600)
	if err != nil {
		panic(err)
	}

	// Manually commit the changes
	cmd := exec.Command("git", "add", fileName)
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	cmd = exec.Command("git", "commit", "-m", "Manual commit message 1", "-m", "With a body")
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Author Name", "GIT_AUTHOR_EMAIL=email@address.com")
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("STDOUT:", stdout.String())
		fmt.Println("STDERR:", stderr.String())
		log.Fatal(err)
		panic(err)
	}

	data = bytes.ReplaceAll(data, []byte("banana"), []byte("pineapple"))

	err = os.WriteFile(fileName, data, 0600)
	if err != nil {
		panic(err)
	}

	// Manually commit the changes again
	cmd = exec.Command("git", "add", fileName)
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	cmd = exec.Command("git", "commit", "-m", "Manual commit message 2")
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Author Name", "GIT_AUTHOR_EMAIL=email@address.com")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

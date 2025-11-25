package main

import (
	"bytes"
	"fmt"
	"io"
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

	cmd = exec.Command("git", "commit", "-m", "Manual commit message 1", "-m", "With a body", "--author", "Author Name <email@address.com>")
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	if cmd.Run(); err != nil {
		fmt.Println(stdout.String())
		fmt.Println(stdout.String())
		panic(err)
	}

	if err := cmd.Run(); err != nil {
		slurp, _ := io.ReadAll(stderr)
		fmt.Println(string(slurp))
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

	cmd = exec.Command("git", "commit", "-m", "Manual commit message 2", "--author", "Author Name <email@address.com>")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

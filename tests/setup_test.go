package tests

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	command := exec.Command("go", "build", "-o", "scripts/changer/main", "scripts/changer/main.go")
	if err := command.Run(); err != nil {
		panic(err)
	}

	command = exec.Command("go", "build", "-o", "scripts/printer/main", "scripts/printer/main.go")
	if err := command.Run(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

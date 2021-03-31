package tests

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
)

var changerBinaryPath string
var printerBinaryPath string

func TestMain(m *testing.M) {
	switch runtime.GOOS {
	case "windows":
		changerBinaryPath = "scripts/changer/main.exe"
		printerBinaryPath = "scripts/printer/main.exe"
	default:
		changerBinaryPath = "scripts/changer/main"
		printerBinaryPath = "scripts/printer/main"
	}

	command := exec.Command("go", "build", "-o", changerBinaryPath, "scripts/changer/main.go")
	if err := command.Run(); err != nil {
		panic(err)
	}

	command = exec.Command("go", "build", "-o", printerBinaryPath, "scripts/printer/main.go")
	if err := command.Run(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

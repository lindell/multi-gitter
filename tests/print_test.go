package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/lindell/multi-gitter/cmd"
	"github.com/lindell/multi-gitter/tests/vcmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrint(t *testing.T) {
	vcMock := &vcmock.VersionController{}
	cmd.OverrideVersionController = vcMock

	tmpDir, err := ioutil.TempDir(os.TempDir(), "multi-git-test-run-")
	assert.NoError(t, err)

	workingDir, err := os.Getwd()
	assert.NoError(t, err)

	changeRepo := createRepo(t, "test-1", "i like apples")
	changeRepo2 := createRepo(t, "test-2", "i like my apple")
	noChangeRepo := createRepo(t, "test-3", "i like oranges")
	vcMock.AddRepository(changeRepo)
	vcMock.AddRepository(changeRepo2)
	vcMock.AddRepository(noChangeRepo)

	runLogFile := path.Join(tmpDir, "print-log.txt")
	outFile := path.Join(tmpDir, "out.txt")

	command := cmd.RootCmd()
	command.SetArgs([]string{"print",
		"--log-file", runLogFile,
		"--output", outFile,
		fmt.Sprintf(`go run %s`, path.Join(workingDir, "scripts/printer/main.go")),
	})
	err = command.Execute()
	assert.NoError(t, err)

	// Verify that the output was correct
	outData, err := ioutil.ReadFile(outFile)
	require.NoError(t, err)
	assert.Equal(t, "i like apples\ni like my apple\ni like oranges\n", string(outData))
}

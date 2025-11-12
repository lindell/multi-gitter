package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/lindell/multi-gitter/cmd"
	"github.com/lindell/multi-gitter/tests/vcmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrint(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedStdout string
		expectedStderr string
	}{
		{
			name:           "all",
			args:           []string{},
			expectedStdout: "i like apples\ni like my apple\ni like oranges\n",
			expectedStderr: "I LIKE APPLES\nI LIKE MY APPLE\nI LIKE ORANGES\n",
		},

		{
			name:           "skip repo",
			args:           []string{"--skip-repo", "owner/test-no-change"},
			expectedStdout: "i like apples\ni like my apple\n",
			expectedStderr: "I LIKE APPLES\nI LIKE MY APPLE\n",
		},
		{
			name:           "include repo",
			args:           []string{"--repo-include", "\\d$"},
			expectedStdout: "i like apples\ni like my apple\n",
			expectedStderr: "I LIKE APPLES\nI LIKE MY APPLE\n",
		},
		{
			name:           "exclude repo",
			args:           []string{"--repo-exclude", "owner/test-[a-z]+"},
			expectedStdout: "i like apples\ni like my apple\n",
			expectedStderr: "I LIKE APPLES\nI LIKE MY APPLE\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vcMock := &vcmock.VersionController{}
			defer vcMock.Clean()
			cmd.OverrideVersionController = vcMock

			tmpDir, err := os.MkdirTemp(os.TempDir(), "multi-git-test-run-")
			assert.NoError(t, err)

			workingDir, err := os.Getwd()
			assert.NoError(t, err)

			changeRepo := createRepo(t, "owner", "test-1", "i like apples")
			changeRepo2 := createRepo(t, "owner", "test-2", "i like my apple")
			noChangeRepo := createRepo(t, "owner", "test-no-change", "i like oranges")
			vcMock.AddRepository(changeRepo)
			vcMock.AddRepository(changeRepo2)
			vcMock.AddRepository(noChangeRepo)

			runLogFile := filepath.Join(tmpDir, "print-log.txt")
			outFile := filepath.Join(tmpDir, "out.txt")
			errOutFile := filepath.Join(tmpDir, "err-out.txt")

			baseArgs := []string{
				"print",
				"--log-file", filepath.ToSlash(runLogFile),
				"--output", filepath.ToSlash(outFile),
				"--error-output", filepath.ToSlash(errOutFile),
			}
			baseArgs = append(baseArgs, tt.args...)
			baseArgs = append(baseArgs,
				fmt.Sprintf(`go run %s`, normalizePath(filepath.Join(workingDir, "scripts/printer/main.go"))),
			)
			command := cmd.RootCmd()
			command.SetArgs(baseArgs)
			err = command.Execute()
			assert.NoError(t, err)

			// Verify that the output was correct
			outData, err := os.ReadFile(outFile)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStdout, string(outData))

			// Verify that the error output was correct
			errOutData, err := os.ReadFile(errOutFile)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStderr, string(errOutData))
		})
	}
}

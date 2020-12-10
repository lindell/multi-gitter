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

type outputs struct {
	out    string
	logOut string
}

func TestTable(t *testing.T) {
	workingDir, err := os.Getwd()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		vc       *vcmock.VersionController
		vcCreate func(t *testing.T) *vcmock.VersionController // Can be used if advanced setup is needed for the vc

		args   []string
		verify func(t *testing.T, vcMock *vcmock.VersionController, outputs outputs)

		expectErr bool
	}{
		{
			name: "simple",
			vc: &vcmock.VersionController{
				Repositories: []vcmock.Repository{
					createRepo(t, "should-change", "i like apples"),
				},
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "custom message",
				fmt.Sprintf(`go run %s`, path.Join(workingDir, "scripts/changer/main.go")),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, outputs outputs) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[0].Head)
				assert.Equal(t, "custom message", vcMock.PullRequests[0].Title)

				assert.Contains(t, outputs.logOut, "Running on 1 repositories")
				assert.Contains(t, outputs.logOut, "Cloning and running script")
				assert.Contains(t, outputs.logOut, "Change done, creating pull request")

				assert.Equal(t, `Repositories with a successful run:
  should-change
`, outputs.out)
			},
		},

		{
			name: "failing base-branch",
			vc: &vcmock.VersionController{
				Repositories: []vcmock.Repository{
					createRepo(t, "should-change", "i like apples"),
				},
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"--base-branch", "custom-base-branch",
				"-m", "custom message",
				fmt.Sprintf(`go run %s`, path.Join(workingDir, "scripts/changer/main.go")),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, outputs outputs) {
				require.Len(t, vcMock.PullRequests, 0)
				assert.Contains(t, outputs.logOut, `msg="couldn't find remote ref \"refs/heads/custom-base-branch\""`)
			},
		},

		{
			name: "success base-branch",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				repo := createRepo(t, "should-change", "i like apples")
				changeBranch(t, repo.Path, "custom-base-branch", true)
				changeTestFile(t, repo.Path, "i like apple", "test change")
				changeBranch(t, repo.Path, "master", false)
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						repo,
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"--base-branch", "custom-base-branch",
				"-m", "custom message",
				fmt.Sprintf(`go run %s`, path.Join(workingDir, "scripts/changer/main.go")),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, outputs outputs) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Equal(t, "custom-base-branch", vcMock.PullRequests[0].Base)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[0].Head)
				assert.Equal(t, "custom message", vcMock.PullRequests[0].Title)

				changeBranch(t, vcMock.Repositories[0].Path, "custom-branch-name", false)
				assert.Equal(t, "i like banana", readTestFile(t, vcMock.Repositories[0].Path))
			},
		},

		{
			name: "reviewers",
			vc: &vcmock.VersionController{
				Repositories: []vcmock.Repository{
					createRepo(t, "should-change", "i like apples"),
				},
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-m", "custom message",
				"-r", "reviewer1,reviewer2",
				fmt.Sprintf(`go run %s`, path.Join(workingDir, "scripts/changer/main.go")),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, outputs outputs) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Len(t, vcMock.PullRequests[0].Reviewers, 2)
				assert.Contains(t, vcMock.PullRequests[0].Reviewers, "reviewer1")
				assert.Contains(t, vcMock.PullRequests[0].Reviewers, "reviewer2")
			},
		},

		{
			name: "random reviewers",
			vc: &vcmock.VersionController{
				Repositories: []vcmock.Repository{
					createRepo(t, "should-change", "i like apples"),
				},
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-m", "custom message",
				"-r", "reviewer1,reviewer2,reviewer3",
				"--max-reviewers", "2",
				fmt.Sprintf(`go run %s`, path.Join(workingDir, "scripts/changer/main.go")),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, outputs outputs) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Len(t, vcMock.PullRequests[0].Reviewers, 2)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logFile, err := ioutil.TempFile(os.TempDir(), "multi-gitter-test-log")
			require.NoError(t, err)
			defer os.Remove(logFile.Name())

			outFile, err := ioutil.TempFile(os.TempDir(), "multi-gitter-test-output")
			require.NoError(t, err)
			defer os.Remove(outFile.Name())

			var vc *vcmock.VersionController
			if test.vcCreate != nil {
				vc = test.vcCreate(t)
			} else {
				vc = test.vc
			}
			cmd.OverrideVersionController = vc

			command := cmd.RootCmd()
			command.SetArgs(append(
				test.args,
				"--log-file", logFile.Name(),
				"--output", outFile.Name(),
			))
			err = command.Execute()
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			logData, err := ioutil.ReadAll(logFile)
			assert.NoError(t, err)

			outData, err := ioutil.ReadAll(outFile)
			assert.NoError(t, err)

			test.verify(t, vc, outputs{
				logOut: string(logData),
				out:    string(outData),
			})
		})
	}
}

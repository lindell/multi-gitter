package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/lindell/multi-gitter/cmd"
	"github.com/lindell/multi-gitter/tests/vcmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type runData struct {
	out    string
	logOut string
	took   time.Duration
}

func TestTable(t *testing.T) {
	workingDir, err := os.Getwd()
	assert.NoError(t, err)

	tests := []struct {
		name     string
		vc       *vcmock.VersionController
		vcCreate func(t *testing.T) *vcmock.VersionController // Can be used if advanced setup is needed for the vc

		args   []string
		verify func(t *testing.T, vcMock *vcmock.VersionController, runData runData)

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
				path.Join(workingDir, "scripts/changer/main"),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[0].Head)
				assert.Equal(t, "custom message", vcMock.PullRequests[0].Title)

				assert.Contains(t, runData.logOut, "Running on 1 repositories")
				assert.Contains(t, runData.logOut, "Cloning and running script")
				assert.Contains(t, runData.logOut, "Change done, creating pull request")

				assert.Equal(t, `Repositories with a successful run:
  should-change
`, runData.out)
			},
		},

		{
			name: "with go run",
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
				fmt.Sprintf("go run %s", path.Join(workingDir, "scripts/changer/main.go")),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[0].Head)
				assert.Equal(t, "custom message", vcMock.PullRequests[0].Title)

				assert.Contains(t, runData.logOut, "Running on 1 repositories")
				assert.Contains(t, runData.logOut, "Cloning and running script")
				assert.Contains(t, runData.logOut, "Change done, creating pull request")

				assert.Equal(t, `Repositories with a successful run:
  should-change
`, runData.out)
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
				path.Join(workingDir, "scripts/changer/main"),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 0)
				assert.Contains(t, runData.logOut, `msg="couldn't find remote ref \"refs/heads/custom-base-branch\""`)
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
				path.Join(workingDir, "scripts/changer/main"),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
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
				path.Join(workingDir, "scripts/changer/main"),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
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
				path.Join(workingDir, "scripts/changer/main"),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Len(t, vcMock.PullRequests[0].Reviewers, 2)
			},
		},

		{
			name: "dry run",
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
				"-B", "custom-branch-name",
				"--dry-run",
				path.Join(workingDir, "scripts/changer/main"),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 0)
				assert.True(t, branchExist(t, vcMock.Repositories[0].Path, "master"))
				assert.False(t, branchExist(t, vcMock.Repositories[0].Path, "custom-branch-name"))
			},
		},

		{
			name: "parallel",
			vc: &vcmock.VersionController{
				Repositories: []vcmock.Repository{
					createRepo(t, "should-change-1", "i like apples"),
					createRepo(t, "should-change-2", "i like apples"),
					createRepo(t, "should-change-3", "i like apples"),
					createRepo(t, "should-change-4", "i like apples"),
					createRepo(t, "should-change-5", "i like apples"),
					createRepo(t, "should-change-6", "i like apples"),
				},
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-m", "custom message",
				"-B", "custom-branch-name",
				"-C", "3",
				fmt.Sprintf("%s -sleep 100ms", path.Join(workingDir, "scripts/changer/main")),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 6)
				require.Less(t, runData.took.Milliseconds(), int64(600))
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
			before := time.Now()
			err = command.Execute()
			took := time.Since(before)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			logData, err := ioutil.ReadAll(logFile)
			assert.NoError(t, err)

			outData, err := ioutil.ReadAll(outFile)
			assert.NoError(t, err)

			test.verify(t, vc, runData{
				logOut: string(logData),
				out:    string(outData),
				took:   took,
			})
		})
	}
}

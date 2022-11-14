package tests

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/lindell/multi-gitter/cmd"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/lindell/multi-gitter/tests/vcmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type runData struct {
	out    string
	logOut string
	cmdOut string
	took   time.Duration
}

type gitBackend string

const (
	gitBackendGo  gitBackend = "go"
	gitBackendCmd gitBackend = "cmd"
)

type skipType string

const (
	skipTypeTimeDependent skipType = "time-dependent"
)

// skipTypes is a list of types that can be skipped. These can be set with build tags
var skipTypes = []skipType{}

func init() {
	for _, t := range strings.Split(os.Getenv("SKIP_TYPES"), ",") {
		skipTypes = append(skipTypes, skipType(t))
	}
}

func skipOverlap(tt1, tt2 []skipType) bool {
	for _, t1 := range tt1 {
		for _, t2 := range tt2 {
			if t1 == t2 {
				return true
			}
		}
	}
	return false
}

var gitBackends = []gitBackend{gitBackendGo, gitBackendCmd}

func containsGitBackend(gitBackends []gitBackend, gitBackend gitBackend) bool {
	for _, gb := range gitBackends {
		if gb == gitBackend {
			return true
		}
	}
	return false
}

func TestTable(t *testing.T) {
	workingDir, err := os.Getwd()
	assert.NoError(t, err)

	changerBinaryPath := normalizePath(filepath.Join(workingDir, changerBinaryPath))

	tests := []struct {
		name        string
		gitBackends []gitBackend                                 // If set, use only the specified git backends, otherwise use all
		vcCreate    func(t *testing.T) *vcmock.VersionController // Can be used if advanced setup is needed for the vc

		args   []string
		verify func(t *testing.T, vcMock *vcmock.VersionController, runData runData)

		skipTypes []skipType // Defined if a test should be skipped in some cases

		expectErr bool
	}{
		{
			name: "simple",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "custom message",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[0].Head)
				assert.Equal(t, "master", vcMock.PullRequests[0].Base)
				assert.Equal(t, "custom message", vcMock.PullRequests[0].Title)

				assert.Contains(t, runData.logOut, "Running on 1 repositories")
				assert.Contains(t, runData.logOut, "Cloning and running script")
				assert.Contains(t, runData.logOut, "Pushing changes to remote")
				assert.Contains(t, runData.logOut, "Creating pull request")

				assert.Equal(t, `Repositories with a successful run:
  owner/should-change #1
`, runData.out)
			},
		},

		{
			name: "with go run",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "custom message",
				fmt.Sprintf("go run %s", normalizePath(filepath.Join(workingDir, "scripts/changer/main.go"))),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[0].Head)
				assert.Equal(t, "custom message", vcMock.PullRequests[0].Title)

				assert.Contains(t, runData.logOut, "Running on 1 repositories")
				assert.Contains(t, runData.logOut, "Cloning and running script")
				assert.Contains(t, runData.logOut, "Pushing changes to remote")
				assert.Contains(t, runData.logOut, "Creating pull request")

				assert.Equal(t, `Repositories with a successful run:
  owner/should-change #1
`, runData.out)
			},
		},

		{
			name:        "failing base-branch",
			gitBackends: []gitBackend{gitBackendGo},
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
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
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 0)
				assert.Contains(t, runData.logOut, `msg="could not clone from the remote: couldn't find remote ref \"refs/heads/custom-base-branch\""`)
			},
		},

		{
			name:        "failing base-branch",
			gitBackends: []gitBackend{gitBackendCmd},
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
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
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 0)
				assert.Contains(t, runData.logOut, `msg="Remote branch custom-base-branch not found in upstream origin"`)
			},
		},

		{
			name: "success base-branch",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				repo := createRepo(t, "owner", "should-change", "i like apples")
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
				changerBinaryPath,
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
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-m", "custom message",
				"-r", "reviewer1,reviewer2",
				changerBinaryPath,
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
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-m", "custom message",
				"-r", "reviewer1,reviewer2,reviewer3",
				"--max-reviewers", "2",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Len(t, vcMock.PullRequests[0].Reviewers, 2)
			},
		},

		{
			name: "dry run",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-m", "custom message",
				"-B", "custom-branch-name",
				"--dry-run",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 0)
				assert.True(t, branchExist(t, vcMock.Repositories[0].Path, "master"))
				assert.False(t, branchExist(t, vcMock.Repositories[0].Path, "custom-branch-name"))
			},
		},

		{
			name:      "parallel",
			skipTypes: []skipType{skipTypeTimeDependent}, // This test is time dependent, don't run it in CI since some runs might be to slow
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change-1", "i like apples"),
						createRepo(t, "owner", "should-change-2", "i like apples"),
						createRepo(t, "owner", "should-change-3", "i like apples"),
						createRepo(t, "owner", "should-change-4", "i like apples"),
						createRepo(t, "owner", "should-change-5", "i like apples"),
						createRepo(t, "owner", "should-change-6", "i like apples"),
						createRepo(t, "owner", "should-change-7", "i like apples"),
						createRepo(t, "owner", "should-change-8", "i like apples"),
						createRepo(t, "owner", "should-change-9", "i like apples"),
						createRepo(t, "owner", "should-change-10", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-m", "custom message",
				"-B", "custom-branch-name",
				"-C", "7",
				fmt.Sprintf("%s -sleep 500ms", changerBinaryPath),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 10)
				require.Less(t, runData.took.Milliseconds(), int64(5000))
			},
		},

		{
			name: "existing head branch",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				repo := createRepo(t, "owner", "already-existing-branch", "i like apples")
				changeBranch(t, repo.Path, "custom-branch-name", true)
				changeTestFile(t, repo.Path, "i like apple", "test change")
				changeBranch(t, repo.Path, "master", false)
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						repo,
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "custom message",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Contains(t, runData.logOut, "Running on 2 repositories")
				assert.Contains(t, runData.logOut, "Cloning and running script")

				assert.Equal(t, `The new branch does already exist:
  owner/already-existing-branch
Repositories with a successful run:
  owner/should-change #1
`, runData.out)
			},
		},

		{
			name: "skip-repo",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-skip", "i like my oranges"),
						createRepo(t, "owner", "should-not-skip", "i like my apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "test",
				"--skip-repo", "owner/should-skip",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)

				assert.Contains(t, runData.logOut, "Skipping owner/should-skip")
				assert.Equal(t, `Repositories with a successful run:
  owner/should-not-skip #1
`, runData.out)
				assert.Equal(t, "i like my oranges", readTestFile(t, vcMock.Repositories[0].Path))
			},
		},

		{
			name: "skip-pr",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				repo := createRepo(t, "owner", "should-change", "i like apples")

				// Change branch so that it's not the one we are expected to push to.
				// If this can be avoided, it would be good.
				changeBranch(t, repo.Path, "test", true)

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
				"-m", "custom message",
				"--skip-pr",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 0)

				assert.Contains(t, runData.logOut, "Running on 1 repositories")
				assert.Contains(t, runData.logOut, "Cloning and running script")

				assert.Equal(t, `Repositories with a successful run:
  owner/should-change
`, runData.out)

				changeBranch(t, vcMock.Repositories[0].Path, "master", false)

				assert.False(t, branchExist(t, vcMock.Repositories[0].Path, "custom-branch-name"))
				assert.False(t, branchExist(t, vcMock.Repositories[0].Path, "multi-gitter-branch"))
				assert.Equal(t, "i like bananas", readTestFile(t, vcMock.Repositories[0].Path))
			},
		},

		{
			name: "autocomplete org",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{}
			},
			args: []string{
				"__complete", "run",
				"--org", "dynamic-org",
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				assert.Equal(t, "static-org\ndynamic-org\n:4\nCompletion ended with directive: ShellCompDirectiveNoFileComp\n", runData.cmdOut)
			},
		},

		{
			name: "autocomplete user",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{}
			},
			args: []string{
				"__complete", "run",
				"--user", "dynamic-user",
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				assert.Equal(t, "static-user\ndynamic-user\n:4\nCompletion ended with directive: ShellCompDirectiveNoFileComp\n", runData.cmdOut)
			},
		},

		{
			name: "autocomplete repo",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{}
			},
			args: []string{
				"__complete", "run",
				"--repo", "dynamic-repo",
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				assert.Equal(t, "static-repo\ndynamic-repo\n:4\nCompletion ended with directive: ShellCompDirectiveNoFileComp\n", runData.cmdOut)
			},
		},

		{
			name: "debug log",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "custom message",
				"--log-level", "debug",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[0].Head)
				assert.Equal(t, "master", vcMock.PullRequests[0].Base)
				assert.Equal(t, "custom message", vcMock.PullRequests[0].Title)

				assert.Contains(t, runData.logOut, "Running on 1 repositories")
				assert.Contains(t, runData.logOut, "Cloning and running script")
				assert.Contains(t, runData.logOut, "Pushing changes to remote")
				assert.Contains(t, runData.logOut, "Creating pull request")

				assert.Equal(t, `Repositories with a successful run:
  owner/should-change #1
`, runData.out)
				assert.Contains(t, runData.logOut, `--- a/test.txt\n+++ b/test.txt\n@@ -1 +1 @@\n-i like apples\n\\ No newline at end of file\n+i like bananas\n\\ No newline at end of file\n`)
			},
		},

		{
			name: "gitignore",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				repo := createRepo(t, "owner", "should-change", "i like apples")
				addFile(t, repo.Path, ".gitignore", "node_modules", "added .gitignore")
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
				"-m", "custom message",
				fmt.Sprintf("go run %s -filenames node_modules/react/README.md,src/index.js -data test", normalizePath(filepath.Join(workingDir, "scripts/adder/main.go"))),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)

				changeBranch(t, vcMock.Repositories[0].Path, "custom-branch-name", false)
				assert.Equal(t, "test", readFile(t, vcMock.Repositories[0].Path, "src/index.js"))
				assert.False(t, fileExist(t, vcMock.Repositories[0].Path, "node_modules/react/README.md"))
			},
		},

		{
			name: "no gitignore",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				repo := createRepo(t, "owner", "should-change", "i like apples")
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
				"-m", "custom message",
				fmt.Sprintf("go run %s -filenames node_modules/react/README.md,src/index.js -data test", normalizePath(filepath.Join(workingDir, "scripts/adder/main.go"))),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)

				changeBranch(t, vcMock.Repositories[0].Path, "custom-branch-name", false)
				assert.Equal(t, "test", readFile(t, vcMock.Repositories[0].Path, "src/index.js"))
				assert.True(t, fileExist(t, vcMock.Repositories[0].Path, "node_modules/react/README.md"))
			},
		},

		{
			name: "fork mode",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "custom message",
				"--fork",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[0].Head)
				assert.Equal(t, "master", vcMock.PullRequests[0].Base)
				assert.Equal(t, "custom message", vcMock.PullRequests[0].Title)

				assert.Contains(t, runData.logOut, "Running on 1 repositories")
				assert.Contains(t, runData.logOut, "Cloning and running script")
				assert.Contains(t, runData.logOut, "Forking repository")
				assert.Contains(t, runData.logOut, "Pushing changes to remote")
				assert.Contains(t, runData.logOut, "Creating pull request")

				assert.Equal(t, `Repositories with a successful run:
  owner/should-change #1
`, runData.out)

				assert.False(t, branchExist(t, vcMock.Repositories[0].Path, "custom-branch-name"))
				changeBranch(t, vcMock.Repositories[0].Path+"-forked-default-owner", "custom-branch-name", false)
				assert.Equal(t, "i like bananas", readTestFile(t, vcMock.Repositories[0].Path+"-forked-default-owner"))
			},
		},

		{
			name: "fork mode with specified owner",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "custom message",
				"--fork",
				"--fork-owner", "custom-org",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[0].Head)
				assert.Equal(t, "master", vcMock.PullRequests[0].Base)
				assert.Equal(t, "custom message", vcMock.PullRequests[0].Title)

				assert.Contains(t, runData.logOut, "Running on 1 repositories")
				assert.Contains(t, runData.logOut, "Cloning and running script")
				assert.Contains(t, runData.logOut, "Forking repository")
				assert.Contains(t, runData.logOut, "Pushing changes to remote")
				assert.Contains(t, runData.logOut, "Creating pull request")

				assert.Equal(t, `Repositories with a successful run:
  owner/should-change #1
`, runData.out)

				assert.False(t, branchExist(t, vcMock.Repositories[0].Path, "custom-branch-name"))
				changeBranch(t, vcMock.Repositories[0].Path+"-forked-custom-org", "custom-branch-name", false)
				assert.Equal(t, "i like bananas", readTestFile(t, vcMock.Repositories[0].Path+"-forked-custom-org"))
			},
		},

		{
			name: "config file",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"--config", "test-config.yaml",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[0].Head)
				assert.Equal(t, "master", vcMock.PullRequests[0].Base)
				assert.Equal(t, "config-message", vcMock.PullRequests[0].Title)

				assert.Equal(t, `Repositories with a successful run:
  owner/should-change #1
`, runData.out)
			},
		},

		{
			name: "assignees",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-m", "custom message",
				"-a", "assignee1,assignee2",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				require.Len(t, vcMock.PullRequests[0].Assignees, 2)
				assert.Contains(t, vcMock.PullRequests[0].Assignees, "assignee1")
				assert.Contains(t, vcMock.PullRequests[0].Assignees, "assignee2")
			},
		},

		{
			name: "replace conflict strategy",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				// Create and run mg with these different repositories:
				// 1. No branch exist when running, change is expected
				// 2. A branch exist, but no PR
				// 3. A branch and pull request exist
				// 4. No branch exist when running, no change should be made
				// 5. A branch exist, no change compared to the base branch should be made
				repo1 := createRepo(t, "owner", "no-existing-branch", "i like apples")
				repo2 := createRepo(t, "owner", "existing-branch", "i like apples")
				changeBranch(t, repo2.Path, "custom-branch-name", true)
				changeTestFile(t, repo2.Path, "i like apple", "test change")
				changeBranch(t, repo2.Path, "master", false)
				repo3 := createRepo(t, "owner", "existing-pr", "i like apples")
				changeBranch(t, repo3.Path, "custom-branch-name", true)
				changeTestFile(t, repo3.Path, "i like apple", "test change")
				changeBranch(t, repo3.Path, "master", false)
				repo4 := createRepo(t, "owner", "no-change-1", "i like oranges")
				repo5 := createRepo(t, "owner", "no-change-2", "i like oranges")
				changeBranch(t, repo5.Path, "custom-branch-name", true)
				changeTestFile(t, repo5.Path, "i like apples", "test change")
				changeBranch(t, repo5.Path, "master", false)

				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						repo1,
						repo2,
						repo3,
						repo4,
						repo5,
					},
					PullRequests: []vcmock.PullRequest{
						{
							PRStatus:   scm.PullRequestStatusSuccess,
							PRNumber:   42,
							Repository: repo3,
							NewPullRequest: scm.NewPullRequest{
								Head: "custom-branch-name",
							},
						},
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "custom message",
				"--conflict-strategy", "replace",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				fmt.Println(runData.logOut)

				require.Len(t, vcMock.PullRequests, 3)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[1].Head)
				assert.Equal(t, "master", vcMock.PullRequests[1].Base)
				assert.Equal(t, "custom message", vcMock.PullRequests[1].Title)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[2].Head)
				assert.Equal(t, "master", vcMock.PullRequests[2].Base)
				assert.Equal(t, "custom message", vcMock.PullRequests[2].Title)

				assert.Contains(t, runData.logOut, "Running on 5 repositories")
				assert.Equal(t, 5, strings.Count(runData.logOut, "Cloning and running script"))
				assert.Equal(t, 3, strings.Count(runData.logOut, "Pushing changes to remote"))
				assert.Equal(t, 2, strings.Count(runData.logOut, "Creating pull request"))
				assert.Equal(t, 1, strings.Count(runData.logOut, "Skip creating pull requests since one is already open"))

				assert.Equal(t, `No data was changed:
  owner/no-change-1
  owner/no-change-2
Repositories with a successful run:
  owner/no-existing-branch #1
  owner/existing-branch #2
  owner/existing-pr #42
`, runData.out)

				changeBranch(t, vcMock.Repositories[0].Path, "custom-branch-name", false)
				assert.Equal(t, "i like bananas", readTestFile(t, vcMock.Repositories[0].Path))
				changeBranch(t, vcMock.Repositories[1].Path, "custom-branch-name", false)
				assert.Equal(t, "i like bananas", readTestFile(t, vcMock.Repositories[1].Path))
				changeBranch(t, vcMock.Repositories[2].Path, "custom-branch-name", false)
				assert.Equal(t, "i like bananas", readTestFile(t, vcMock.Repositories[2].Path))
			},
		},

		{
			name: "multi line message body",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "custom message\nwith more info\nand even more",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Equal(t, "custom-branch-name", vcMock.PullRequests[0].Head)
				assert.Equal(t, "master", vcMock.PullRequests[0].Base)
				assert.Equal(t, "custom message", vcMock.PullRequests[0].Title)
				assert.Equal(t, "with more info\nand even more", vcMock.PullRequests[0].Body)
			},
		},

		{
			name: "no repositories found",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "custom message",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 0)
				assert.NotContains(t, runData.logOut, "Running on")
				assert.Contains(t, runData.logOut, "No repositories found")
			},
		},

		{
			name: "PRs as draft",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples as draft"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "custom message",
				"--draft",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.True(t, vcMock.PullRequests[0].Draft)
			},
		},

		{
			name: "labels",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "custom-branch-name",
				"-m", "custom message",
				"--labels", "label1,label2",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)
				assert.Equal(t, []string{"label1", "label2"}, vcMock.PullRequests[0].Labels)
			},
		},

		{
			name: "remove files",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				repo := createRepo(t, "owner", "should-delete", "i like apples")
				addFile(t, repo.Path, "test_file", "some content", "added test_file")
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
				"-m", "custom message",
				fmt.Sprintf("go run %s", normalizePath(filepath.Join(workingDir, "scripts/remover/main.go"))),
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 1)

				changeBranch(t, vcMock.Repositories[0].Path, "custom-branch-name", false)

				assert.False(t, fileExist(t, vcMock.Repositories[0].Path, "test_file"))
			},
		},

		{
			name: "same feature and base branch name",
			vcCreate: func(t *testing.T) *vcmock.VersionController {
				return &vcmock.VersionController{
					Repositories: []vcmock.Repository{
						createRepo(t, "owner", "should-not-change", "i like apples"),
					},
				}
			},
			args: []string{
				"run",
				"--author-name", "Test Author",
				"--author-email", "test@example.com",
				"-B", "master",
				"-m", "custom message",
				changerBinaryPath,
			},
			verify: func(t *testing.T, vcMock *vcmock.VersionController, runData runData) {
				require.Len(t, vcMock.PullRequests, 0)
				assert.Equal(t, "Both the feature branch and base branch was named master, if you intended to push directly into the base branch, please use the `skip-pr` option:\n  owner/should-not-change\n", runData.out)
			},
		},
	}

	for _, gitBackend := range gitBackends {
		for _, test := range tests {
			// Some tests should only be run with specific git backends
			if test.gitBackends != nil && !containsGitBackend(test.gitBackends, gitBackend) {
				continue
			}

			t.Run(fmt.Sprintf("%s_%s", gitBackend, test.name), func(t *testing.T) {
				// Skip some tests depending on the values in skipTypes
				if skipOverlap(skipTypes, test.skipTypes) {
					t.SkipNow()
				}

				logFile, err := os.CreateTemp(os.TempDir(), "multi-gitter-test-log")
				require.NoError(t, err)
				defer os.Remove(logFile.Name())

				outFile, err := os.CreateTemp(os.TempDir(), "multi-gitter-test-output")
				require.NoError(t, err)
				// defer os.Remove(outFile.Name())

				vc := test.vcCreate(t)
				defer vc.Clean()

				cmd.OverrideVersionController = vc

				cobraBuf := &bytes.Buffer{}

				staticArgs := []string{
					"--log-file", logFile.Name(),
					"--output", outFile.Name(),
					"--git-type", string(gitBackend),
				}

				command := cmd.RootCmd()
				command.SetOut(cobraBuf)
				command.SetErr(cobraBuf)
				command.SetArgs(append(staticArgs, test.args...))
				before := time.Now()
				err = command.Execute()
				took := time.Since(before)
				if test.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}

				logData, err := io.ReadAll(logFile)
				assert.NoError(t, err)

				outData, err := io.ReadAll(outFile)
				assert.NoError(t, err)

				test.verify(t, vc, runData{
					logOut: string(logData),
					out:    string(outData),
					cmdOut: cobraBuf.String(),
					took:   took,
				})
			})
		}
	}
}

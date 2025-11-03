package gerrit

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"testing"

	gogerrit "github.com/andygrunwald/go-gerrit"
	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

type goGerritClientMock struct {
	ListProjectsFunc  func(ctx context.Context, opt *gogerrit.ProjectOptions) (*map[string]gogerrit.ProjectInfo, *gogerrit.Response, error)
	QueryChangesFunc  func(ctx context.Context, opt *gogerrit.QueryChangeOptions) (*[]gogerrit.ChangeInfo, *gogerrit.Response, error)
	AbandonChangeFunc func(ctx context.Context, changeID string, input *gogerrit.AbandonInput) (*gogerrit.ChangeInfo, *gogerrit.Response, error)
	SubmitChangeFunc  func(ctx context.Context, changeID string, input *gogerrit.SubmitInput) (*gogerrit.ChangeInfo, *gogerrit.Response, error)
	GetHEADFunc       func(ctx context.Context, projectName string) (string, *gogerrit.Response, error)
}

func (gcm goGerritClientMock) ListProjects(ctx context.Context, opt *gogerrit.ProjectOptions) (*map[string]gogerrit.ProjectInfo, *gogerrit.Response, error) {
	return gcm.ListProjectsFunc(ctx, opt)
}

func (gcm goGerritClientMock) QueryChanges(ctx context.Context, opt *gogerrit.QueryChangeOptions) (*[]gogerrit.ChangeInfo, *gogerrit.Response, error) {
	return gcm.QueryChangesFunc(ctx, opt)
}

func (gcm goGerritClientMock) AbandonChange(ctx context.Context, changeID string, input *gogerrit.AbandonInput) (*gogerrit.ChangeInfo, *gogerrit.Response, error) {
	return gcm.AbandonChangeFunc(ctx, changeID, input)
}

func (gcm goGerritClientMock) SubmitChange(ctx context.Context, changeID string, input *gogerrit.SubmitInput) (*gogerrit.ChangeInfo, *gogerrit.Response, error) {
	return gcm.SubmitChangeFunc(ctx, changeID, input)
}

func (gcm goGerritClientMock) GetHEAD(ctx context.Context, projectName string) (string, *gogerrit.Response, error) {
	if gcm.GetHEADFunc != nil {
		return gcm.GetHEADFunc(ctx, projectName)
	}
	// Default implementation for tests: return "main" for "main-branch-repo", "master" for others
	if projectName == "main-branch-repo" {
		return "refs/heads/main", nil, nil
	}
	return "refs/heads/master", nil, nil
}

var projects = &map[string]gogerrit.ProjectInfo{
	"another-repo-active": {State: "ACTIVE"},
	"main-branch-repo":    {State: "ACTIVE"},
	"repo-active":         {State: "ACTIVE"},
	"repo-read-only":      {State: "READ_ONLY"},
}

func getChangesForQuery(query string) (*[]gogerrit.ChangeInfo, *gogerrit.Response, error) {
	var data = map[string][]gogerrit.ChangeInfo{
		"footer:MultiGitter-Branch=feature+project:repo-active+is:open": {
			{Project: "repo-active", ChangeID: "I123", Branch: "feature", Number: 1000, Status: "NEW"},
		},
		"footer:MultiGitter-Branch=feature": {
			{Project: "repo-active", ChangeID: "I123", Branch: "feature", Number: 1000, Status: "NEW"},
			{Project: "repo-active", ChangeID: "I000", Branch: "feature", Number: 1001, Status: "ABANDONED"},
			{Project: "repo-active", ChangeID: "I644", Branch: "feature", Number: 1002, Submittable: true},
			{Project: "another-repo-active", ChangeID: "I666", Branch: "feature", Number: 1003, Status: "MERGED"},
			{Project: "name-that-do-not-match", ChangeID: "I000", Branch: "feature", Number: 1004, Submittable: true},
		},
		"footer:MultiGitter-Branch=feature+project:another-repo-active+is:open": {},
		"footer:MultiGitter-Branch=feature+project:read-only+is:open": {
			{Project: "read-only", ChangeID: "I456", Branch: "feature", Number: 1004, Status: "NEW"},
			{Project: "read-only", ChangeID: "I789", Branch: "feature", Number: 1004, Status: "NEW"},
		},
	}

	if strings.Contains(query, "throw-error") {
		return nil, nil, assert.AnError
	}

	if changes, ok := data[query]; ok {
		return &changes, nil, nil
	}

	return &[]gogerrit.ChangeInfo{}, nil, nil
}

func TestGetRepositories(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			ListProjectsFunc: func(_ context.Context, opt *gogerrit.ProjectOptions) (*map[string]gogerrit.ProjectInfo, *gogerrit.Response, error) {
				require.Equal(t, "repo", opt.Regex)
				return projects, nil, nil
			},
		},
		baseURL:  "https://gerrit.com",
		username: "admin",
		token:    "token123",
		repoListing: RepositoryListing{
			RepoSearch: "repo",
		},
	}

	repos, err := g.GetRepositories(context.Background())
	require.NoError(t, err)
	require.Len(t, repos, 3)

	expectedRepos := []struct {
		name          string
		defaultBranch string
		cloneURL      string
	}{
		{"another-repo-active", "master", "https://admin:token123@gerrit.com/a/another-repo-active"},
		{"main-branch-repo", "main", "https://admin:token123@gerrit.com/a/main-branch-repo"},
		{"repo-active", "master", "https://admin:token123@gerrit.com/a/repo-active"},
	}
	for idx, expectedRepo := range expectedRepos {
		repo := repos[idx]
		assert.Equal(t, expectedRepo.name, repo.FullName())
		assert.Equal(t, expectedRepo.defaultBranch, repo.DefaultBranch())
		assert.Equal(t, expectedRepo.cloneURL, repo.CloneURL())
	}
}

func TestGetPullRequests(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			ListProjectsFunc: func(_ context.Context, opt *gogerrit.ProjectOptions) (*map[string]gogerrit.ProjectInfo, *gogerrit.Response, error) {
				require.Equal(t, "repo", opt.Regex)
				return projects, nil, nil
			},
			QueryChangesFunc: func(_ context.Context, opt *gogerrit.QueryChangeOptions) (*[]gogerrit.ChangeInfo, *gogerrit.Response, error) {
				return getChangesForQuery(opt.Query[0])
			},
		},
		baseURL:  "https://gerrit.com",
		username: "admin",
		token:    "token123",
		repoListing: RepositoryListing{
			RepoSearch: "repo",
		},
	}
	prs, err := g.GetPullRequests(context.Background(), "feature")
	require.NoError(t, err)
	require.Len(t, prs, 4)

	expectedPRs := []struct {
		project  string
		changeID string
		number   int
		status   scm.PullRequestStatus
	}{
		{"another-repo-active", "I666", 1003, scm.PullRequestStatusMerged},
		{"repo-active", "I123", 1000, scm.PullRequestStatusPending},
		{"repo-active", "I000", 1001, scm.PullRequestStatusClosed},
		{"repo-active", "I644", 1002, scm.PullRequestStatusSuccess},
	}
	for idx, expectedPR := range expectedPRs {
		pr := prs[idx]
		change := pr.(change)
		assert.Equal(t, expectedPR.status, pr.Status())
		assert.Equal(t, strconv.Itoa(change.number)+": "+change.project, pr.String())
		assert.Equal(t, expectedPR.project, change.project)
		assert.Equal(t, expectedPR.changeID, change.changeID)
		assert.Equal(t, "https://gerrit.com/c/"+change.project+"/+/"+strconv.Itoa(change.number), change.URL())
	}
}

func TestGetOpenPullRequest(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			QueryChangesFunc: func(_ context.Context, opt *gogerrit.QueryChangeOptions) (*[]gogerrit.ChangeInfo, *gogerrit.Response, error) {
				return getChangesForQuery(opt.Query[0])
			},
		},
	}
	tests := []struct {
		repository       string
		expectedErr      bool
		expectedChangeID string
	}{
		{"repo-active", false, "I123"},
		{"another-repo-active", false, ""},
		{"read-only", true, ""},
	}
	for _, test := range tests {
		t.Run(test.repository, func(t *testing.T) {
			repo := repository{name: test.repository}
			pr, err := g.GetOpenPullRequest(context.Background(), repo, "feature")

			if test.expectedErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "More than one open change for branch feature in project "+test.repository)
			} else {
				require.NoError(t, err)
				if test.expectedChangeID == "" {
					require.Nil(t, pr)
				} else {
					require.NotNil(t, pr)
					assert.Equal(t, test.expectedChangeID, pr.(change).changeID)
				}
			}
		})
	}
}

func TestCreatePullRequest(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			QueryChangesFunc: func(_ context.Context, opt *gogerrit.QueryChangeOptions) (*[]gogerrit.ChangeInfo, *gogerrit.Response, error) {
				return getChangesForQuery(opt.Query[0])
			},
		},
	}
	repo := repository{name: "repo-active"}
	pr, err := g.CreatePullRequest(context.Background(), repo, repo, scm.NewPullRequest{Head: "feature"})
	require.NoError(t, err)
	require.NotNil(t, pr)
	assert.Equal(t, "I123", pr.(change).changeID)

	_, err = g.CreatePullRequest(context.Background(), repo, repo, scm.NewPullRequest{Head: "unknown-feature"})
	require.Error(t, err)
}

func TestUpdatePullRequest(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			QueryChangesFunc: func(_ context.Context, opt *gogerrit.QueryChangeOptions) (*[]gogerrit.ChangeInfo, *gogerrit.Response, error) {
				return getChangesForQuery(opt.Query[0])
			},
		},
	}
	repo := repository{name: "repo-active"}
	pr, err := g.UpdatePullRequest(context.Background(), repo, change{}, scm.NewPullRequest{Head: "feature"})
	require.NoError(t, err)
	require.NotNil(t, pr)
	assert.Equal(t, "I123", pr.(change).changeID)

	_, err = g.UpdatePullRequest(context.Background(), repo, change{}, scm.NewPullRequest{Head: "unknown-feature"})
	require.Error(t, err)
}

func TestRemoteReference(t *testing.T) {
	g := &Gerrit{}

	tests := []struct {
		baseBranch      string
		featureBranch   string
		skipPullRequest bool
		pushOnly        bool
		expectedRef     string
	}{
		{"master", "feature", false, false, "refs/for/master"},
		{"master", "feature", true, false, "refs/heads/feature"},
		{"master", "feature", false, true, "refs/heads/feature"},
		{"master", "feature", true, true, "refs/heads/feature"},
	}

	for _, test := range tests {
		t.Run(test.baseBranch+"_"+test.featureBranch, func(t *testing.T) {
			ref := g.RemoteReference(test.baseBranch, test.featureBranch, test.skipPullRequest, test.pushOnly)
			assert.Equal(t, test.expectedRef, ref)
		})
	}
}

func TestFeatureBranchExist(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			QueryChangesFunc: func(_ context.Context, opt *gogerrit.QueryChangeOptions) (*[]gogerrit.ChangeInfo, *gogerrit.Response, error) {
				return getChangesForQuery(opt.Query[0])
			},
		},
	}

	tests := []struct {
		branchName string
		expected   bool
	}{
		{"feature", true},
		{"new-feature", false},
	}

	for _, test := range tests {
		t.Run(test.branchName, func(t *testing.T) {
			repo := repository{name: "repo-active"}
			exist, err := g.FeatureBranchExist(context.Background(), repo, test.branchName)
			require.NoError(t, err)
			assert.Equal(t, test.expected, exist)
		})
	}
}

func TestEnhanceCommit(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			QueryChangesFunc: func(_ context.Context, opt *gogerrit.QueryChangeOptions) (*[]gogerrit.ChangeInfo, *gogerrit.Response, error) {
				return getChangesForQuery(opt.Query[0])
			},
		},
	}

	tests := []struct {
		branchName    string
		expectedErr   bool
		changeIDRegex string
	}{
		{"feature", false, "I123"},
		{"new-feature", false, "I[0-9a-f]{40}"},
		{"feature-that-throw-error", true, ""},
	}

	for _, test := range tests {
		t.Run(test.branchName, func(t *testing.T) {
			repo := repository{name: "repo-active"}
			msg, err := g.EnhanceCommit(context.Background(), repo, test.branchName, "dummy commit message")
			if test.expectedErr {
				require.Error(t, err)
				assert.Equal(t, "dummy commit message", msg)
			} else {
				require.NoError(t, err)
				assert.Regexp(t, regexp.MustCompile(
					"dummy commit message\n\nMultiGitter-Branch: "+test.branchName+"\nChange-Id: "+test.changeIDRegex), msg)
			}
		})
	}
}

func TestMergePullRequest(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			SubmitChangeFunc: func(_ context.Context, changeID string, _ *gogerrit.SubmitInput) (*gogerrit.ChangeInfo, *gogerrit.Response, error) {
				// Ensure correct id is used when a change is submitted
				require.Equal(t, "repo-active~master~Icc717a31a47beb9b5d9aeb8a1d374883afe89030", changeID)
				return &gogerrit.ChangeInfo{}, nil, nil
			},
		},
	}
	pr := change{
		id:       "repo-active~master~Icc717a31a47beb9b5d9aeb8a1d374883afe89030",
		project:  "repo-active",
		branch:   "master",
		changeID: "Icc717a31a47beb9b5d9aeb8a1d374883afe89030",
	}
	err := g.MergePullRequest(context.Background(), pr)
	require.NoError(t, err)
}

func TestClosePullRequest(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			AbandonChangeFunc: func(ctx context.Context, changeID string, input *gogerrit.AbandonInput) (*gogerrit.ChangeInfo, *gogerrit.Response, error) {
				// Ensure correct id is used when a change is abandoned
				require.Equal(t, "repo-active~master~Icc717a31a47beb9b5d9aeb8a1d374883afe89030", changeID)
				return &gogerrit.ChangeInfo{}, nil, nil
			},
		},
	}
	pr := change{
		id:       "repo-active~master~Icc717a31a47beb9b5d9aeb8a1d374883afe89030",
		project:  "repo-active",
		branch:   "master",
		changeID: "Icc717a31a47beb9b5d9aeb8a1d374883afe89030",
	}
	err := g.ClosePullRequest(context.Background(), pr)
	require.NoError(t, err)
}

func TestGetDefaultBranch(t *testing.T) {
	tests := []struct {
		name           string
		projectName    string
		mockResponse   string
		mockError      error
		expectedBranch string
		expectError    bool
	}{
		{
			name:           "successful retrieval with refs/heads/ prefix",
			projectName:    "test-project",
			mockResponse:   "refs/heads/main",
			mockError:      nil,
			expectedBranch: "main",
			expectError:    false,
		},
		{
			name:           "successful retrieval with refs/heads/master",
			projectName:    "legacy-project",
			mockResponse:   "refs/heads/master",
			mockError:      nil,
			expectedBranch: "master",
			expectError:    false,
		},
		{
			name:           "response without refs/heads/ prefix",
			projectName:    "bare-project",
			mockResponse:   "develop",
			mockError:      nil,
			expectedBranch: "develop",
			expectError:    false,
		},
		{
			name:           "API error",
			projectName:    "error-project",
			mockResponse:   "",
			mockError:      errors.New("API error"),
			expectedBranch: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Gerrit{
				client: goGerritClientMock{
					GetHEADFunc: func(ctx context.Context, projectName string) (string, *gogerrit.Response, error) {
						require.Equal(t, tt.projectName, projectName)
						return tt.mockResponse, nil, tt.mockError
					},
				},
			}

			branch, err := g.getDefaultBranch(context.Background(), tt.projectName)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get HEAD branch")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBranch, branch)
			}
		})
	}
}

func TestGetRepositoriesWithSpecificList(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			ListProjectsFunc: func(_ context.Context, opt *gogerrit.ProjectOptions) (*map[string]gogerrit.ProjectInfo, *gogerrit.Response, error) {
				// Verify we don't inject repoSearch when using specific repo list
				require.Empty(t, opt.Regex)
				return projects, nil, nil
			},
		},
		baseURL:  "https://gerrit.com",
		username: "admin",
		token:    "token123",
		repoListing: RepositoryListing{
			Repositories: []string{"repo-active"},
		},
	}

	repos, err := g.GetRepositories(context.Background())
	require.NoError(t, err)
	require.Len(t, repos, 1)

	repo := repos[0]
	assert.Equal(t, "repo-active", repo.FullName())
	assert.Equal(t, "master", repo.DefaultBranch())
	assert.Equal(t, "https://admin:token123@gerrit.com/a/repo-active", repo.CloneURL())

	// Test with non-existent repository in the list
	g.repoListing.Repositories = []string{"repo-read-only"}
	repos, err = g.GetRepositories(context.Background())
	require.NoError(t, err)
	require.Len(t, repos, 0)

	// Test with empty repository list
	g.repoListing.Repositories = []string{}
	repos, err = g.GetRepositories(context.Background())
	require.NoError(t, err)
	require.Len(t, repos, 3) // Should return all active repos

	// Test with error from client
	g.client = goGerritClientMock{
		ListProjectsFunc: func(_ context.Context, _ *gogerrit.ProjectOptions) (*map[string]gogerrit.ProjectInfo, *gogerrit.Response, error) {
			return nil, nil, assert.AnError
		},
	}
	repos, err = g.GetRepositories(context.Background())
	require.Error(t, err)
	require.Nil(t, repos)
}

func TestNewWithConfigEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: Config{
				Username: "user",
				Token:    "token",
				BaseURL:  "https://gerrit.example.com",
				RepoListing: RepositoryListing{
					RepoSearch: "test",
				},
			},
			expectError: false,
		},
		{
			name: "empty username",
			config: Config{
				Username: "",
				Token:    "token",
				BaseURL:  "https://gerrit.example.com",
			},
			expectError: false, // Should still work, just no auth
		},
		{
			name: "invalid base URL",
			config: Config{
				Username: "user",
				Token:    "token",
				BaseURL:  "not-a-url",
			},
			expectError: false, // go-gerrit is permissive with URL formats
		},
		{
			name: "empty base URL",
			config: Config{
				Username: "user",
				Token:    "token",
				BaseURL:  "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.config)
			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, client)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
				assert.Equal(t, tt.config.Username, client.username)
				assert.Equal(t, tt.config.Token, client.token)
				assert.Equal(t, tt.config.BaseURL, client.baseURL)
				assert.Equal(t, tt.config.RepoListing, client.repoListing)
			}
		})
	}
}

func TestGetRepositoriesWithBothFilters(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			ListProjectsFunc: func(_ context.Context, opt *gogerrit.ProjectOptions) (*map[string]gogerrit.ProjectInfo, *gogerrit.Response, error) {
				// When both Repositories and RepoSearch are provided, RepoSearch should be ignored
				require.Empty(t, opt.Regex, "RepoSearch should be ignored when specific repositories are provided")
				return projects, nil, nil
			},
		},
		baseURL:  "https://gerrit.com",
		username: "admin",
		token:    "token123",
		repoListing: RepositoryListing{
			Repositories: []string{"repo-active"},
			RepoSearch:   "should-be-ignored", // This should be ignored when Repositories is provided
		},
	}

	repos, err := g.GetRepositories(context.Background())
	require.NoError(t, err)
	require.Len(t, repos, 1)
	assert.Equal(t, "repo-active", repos[0].FullName())
}

func TestGetRepositoriesWithEmptyRepositoriesButRepoSearch(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			ListProjectsFunc: func(_ context.Context, opt *gogerrit.ProjectOptions) (*map[string]gogerrit.ProjectInfo, *gogerrit.Response, error) {
				// When Repositories is empty, RepoSearch should be used
				require.Equal(t, "active", opt.Regex)
				return projects, nil, nil
			},
		},
		baseURL:  "https://gerrit.com",
		username: "admin",
		token:    "token123",
		repoListing: RepositoryListing{
			Repositories: []string{}, // Empty, so RepoSearch should be used
			RepoSearch:   "active",
		},
	}

	repos, err := g.GetRepositories(context.Background())
	require.NoError(t, err)
	require.Len(t, repos, 3) // Should return all active repos matching "active"
}

func TestGetRepositoriesWithCaseSensitiveMatching(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			ListProjectsFunc: func(_ context.Context, _ *gogerrit.ProjectOptions) (*map[string]gogerrit.ProjectInfo, *gogerrit.Response, error) {
				return projects, nil, nil
			},
		},
		baseURL:  "https://gerrit.com",
		username: "admin",
		token:    "token123",
		repoListing: RepositoryListing{
			Repositories: []string{"REPO-ACTIVE"}, // Different case
		},
	}

	repos, err := g.GetRepositories(context.Background())
	require.NoError(t, err)
	require.Len(t, repos, 0) // Should not match due to case sensitivity
}

func TestGetRepositoriesWithInactiveRepoInList(t *testing.T) {
	g := &Gerrit{
		client: goGerritClientMock{
			ListProjectsFunc: func(_ context.Context, _ *gogerrit.ProjectOptions) (*map[string]gogerrit.ProjectInfo, *gogerrit.Response, error) {
				return projects, nil, nil
			},
		},
		baseURL:  "https://gerrit.com",
		username: "admin",
		token:    "token123",
		repoListing: RepositoryListing{
			Repositories: []string{"repo-read-only"}, // This repo is READ_ONLY, not ACTIVE
		},
	}

	repos, err := g.GetRepositories(context.Background())
	require.NoError(t, err)
	require.Len(t, repos, 0) // Should not return inactive repos even if in the list
}

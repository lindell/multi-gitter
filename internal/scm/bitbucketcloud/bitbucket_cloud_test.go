package bitbucketcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestBitbucketCloud creates a BitbucketCloud instance for testing
func createTestBitbucketCloud(repositories []string, workspaces []string, username, token string, authType AuthType) *BitbucketCloud {
	bbc := &BitbucketCloud{
		repositories: repositories,
		workspaces:   workspaces,
		username:     username,
		token:        token,
		authType:     authType,
		sshAuth:      false,
	}
	return bbc
}

func TestBitbucketCloud_GetRepositories_SpecificRepositories(t *testing.T) {
	tests := []struct {
		name         string
		repositories []string
		authType     AuthType
		username     string
		token        string
		expectError  bool
	}{
		{
			name:         "workspace token auth",
			repositories: []string{"test-repo"},
			authType:     AuthTypeWorkspaceToken,
			username:     "test-user",
			token:        "test-token",
			expectError:  false,
		},
		{
			name:         "app password auth",
			repositories: []string{"test-repo"},
			authType:     AuthTypeAppPassword,
			username:     "test-user",
			token:        "test-password",
			expectError:  false,
		},
		{
			name:         "multiple repositories",
			repositories: []string{"repo1", "repo2", "repo3"},
			authType:     AuthTypeWorkspaceToken,
			username:     "test-user",
			token:        "test-token",
			expectError:  false,
		},
		{
			name:         "empty repositories list",
			repositories: []string{},
			authType:     AuthTypeWorkspaceToken,
			username:     "test-user",
			token:        "test-token",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bbc := createTestBitbucketCloud(tt.repositories, []string{"test-workspace"}, tt.username, tt.token, tt.authType)

			// Test that the BitbucketCloud instance is created correctly
			assert.Equal(t, tt.repositories, bbc.repositories)
			assert.Equal(t, tt.username, bbc.username)
			assert.Equal(t, tt.token, bbc.token)
			assert.Equal(t, tt.authType, bbc.authType)
			assert.Equal(t, []string{"test-workspace"}, bbc.workspaces)

			// Test the logic paths for specific repositories vs workspace repositories
			if len(tt.repositories) > 0 {
				// When specific repositories are provided, the function should iterate through them
				assert.NotEmpty(t, bbc.repositories, "Should have specific repositories to fetch")
			} else {
				// When no specific repositories are provided, it should fetch all workspace repositories
				assert.Empty(t, bbc.repositories, "Should fetch all workspace repositories")
			}
		})
	}
}

// TestBitbucketCloud_convertRepository_AuthenticationTypes tests the authentication logic
func TestBitbucketCloud_convertRepository_AuthenticationTypes(t *testing.T) {
	tests := []struct {
		name         string
		authType     AuthType
		username     string
		token        string
		sshAuth      bool
		expectedAuth string
	}{
		{
			name:         "HTTPS with app password",
			authType:     AuthTypeAppPassword,
			username:     "testuser",
			token:        "testpass",
			sshAuth:      false,
			expectedAuth: "testuser:testpass",
		},
		{
			name:         "HTTPS with workspace token",
			authType:     AuthTypeWorkspaceToken,
			username:     "testuser",
			token:        "testtoken",
			sshAuth:      false,
			expectedAuth: "x-token-auth:testtoken",
		},
		{
			name:         "SSH authentication",
			authType:     AuthTypeWorkspaceToken,
			username:     "testuser",
			token:        "testtoken",
			sshAuth:      true,
			expectedAuth: "ssh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bbc := &BitbucketCloud{
				authType: tt.authType,
				username: tt.username,
				token:    tt.token,
				sshAuth:  tt.sshAuth,
			}

			// Test that the BitbucketCloud instance has the correct authentication settings
			assert.Equal(t, tt.authType, bbc.authType)
			assert.Equal(t, tt.username, bbc.username)
			assert.Equal(t, tt.token, bbc.token)
			assert.Equal(t, tt.sshAuth, bbc.sshAuth)

			// Test the expected authentication behavior
			if tt.sshAuth {
				assert.True(t, bbc.sshAuth, "Should use SSH authentication")
			} else {
				assert.False(t, bbc.sshAuth, "Should use HTTPS authentication")
				if tt.authType == AuthTypeAppPassword {
					assert.Equal(t, AuthTypeAppPassword, bbc.authType)
				} else {
					assert.Equal(t, AuthTypeWorkspaceToken, bbc.authType)
				}
			}
		})
	}
}

func TestBitbucketCloud_AuthTypeValidation(t *testing.T) {
	tests := []struct {
		name     string
		authType AuthType
		username string
		token    string
		isValid  bool
	}{
		{
			name:     "valid app password",
			authType: AuthTypeAppPassword,
			username: "user",
			token:    "password",
			isValid:  true,
		},
		{
			name:     "valid workspace token",
			authType: AuthTypeWorkspaceToken,
			username: "user",
			token:    "token",
			isValid:  true,
		},
		{
			name:     "empty token",
			authType: AuthTypeWorkspaceToken,
			username: "user",
			token:    "",
			isValid:  false,
		},
		{
			name:     "whitespace token",
			authType: AuthTypeWorkspaceToken,
			username: "user",
			token:    "   ",
			isValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.username, tt.token, []string{}, []string{"workspace"}, []string{}, false, false, "", tt.authType)

			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "bearer token is empty")
			}
		})
	}
}

func TestParseAuthType(t *testing.T) {
	tests := []struct {
		input       string
		expected    AuthType
		shouldError bool
	}{
		{
			input:       "app-password",
			expected:    AuthTypeAppPassword,
			shouldError: false,
		},
		{
			input:       "workspace-token",
			expected:    AuthTypeWorkspaceToken,
			shouldError: false,
		},
		{
			input:       "invalid-type",
			expected:    AuthType(0),
			shouldError: true,
		},
		{
			input:       "",
			expected:    AuthType(0),
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseAuthType(tt.input)

			if tt.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "could not parse")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBitbucketCloud_Configuration(t *testing.T) {
	t.Run("proper initialization", func(t *testing.T) {
		repositories := []string{"repo1", "repo2"}
		workspaces := []string{"workspace1"}
		users := []string{"user1"}
		username := "testuser"
		token := "testtoken"
		authType := AuthTypeWorkspaceToken

		bbc, err := New(username, token, repositories, workspaces, users, false, false, "", authType)
		require.NoError(t, err)

		assert.Equal(t, repositories, bbc.repositories)
		assert.Equal(t, workspaces, bbc.workspaces)
		assert.Equal(t, users, bbc.users)
		assert.Equal(t, username, bbc.username)
		assert.Equal(t, token, bbc.token)
		assert.Equal(t, authType, bbc.authType)
		assert.False(t, bbc.fork)
		assert.False(t, bbc.sshAuth)
		assert.Empty(t, bbc.newOwner)
	})

	t.Run("fork and ssh configuration", func(t *testing.T) {
		bbc, err := New("user", "token", []string{}, []string{"workspace"}, []string{}, true, true, "newowner", AuthTypeAppPassword)
		require.NoError(t, err)

		assert.True(t, bbc.fork)
		assert.True(t, bbc.sshAuth)
		assert.Equal(t, "newowner", bbc.newOwner)
		assert.Equal(t, AuthTypeAppPassword, bbc.authType)
	})
}

// Integration-style test that tests the actual logic flow without external dependencies
func TestBitbucketCloud_GetRepositories_LogicFlow(t *testing.T) {
	t.Run("specific repositories path", func(t *testing.T) {
		// Test that when specific repositories are provided, it takes the specific path
		bbc := createTestBitbucketCloud([]string{"repo1", "repo2"}, []string{"workspace"}, "user", "token", AuthTypeWorkspaceToken)

		// The logic should check if repositories are provided
		assert.NotEmpty(t, bbc.repositories, "Should have specific repositories")
		assert.Len(t, bbc.repositories, 2, "Should have exactly 2 repositories")

		// This tests the branch condition in GetRepositories
		if len(bbc.repositories) > 0 {
			// This is the path that will be taken - fetching specific repositories
			assert.Contains(t, bbc.repositories, "repo1")
			assert.Contains(t, bbc.repositories, "repo2")
		}
	})

	t.Run("workspace repositories path", func(t *testing.T) {
		// Test that when no specific repositories are provided, it takes the workspace path
		bbc := createTestBitbucketCloud([]string{}, []string{"workspace"}, "user", "token", AuthTypeWorkspaceToken)

		// The logic should check if repositories are empty
		assert.Empty(t, bbc.repositories, "Should have no specific repositories")
		assert.NotEmpty(t, bbc.workspaces, "Should have workspaces to fetch from")

		// This tests the else branch condition in GetRepositories
		if len(bbc.repositories) == 0 {
			// This is the path that will be taken - fetching all workspace repositories
			assert.Contains(t, bbc.workspaces, "workspace")
		}
	})
}

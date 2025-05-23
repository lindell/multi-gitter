package cmd

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGerritClientValidation(t *testing.T) {
	tests := []struct {
		name        string
		setupFlags  func() *pflag.FlagSet
		verifyFlags bool
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config with repo search",
			setupFlags: func() *pflag.FlagSet {
				fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
				fs.String("username", "testuser", "")
				fs.String("base-url", "https://gerrit.example.com", "")
				fs.StringSlice("repo", []string{}, "")
				fs.String("repo-search", "myrepo", "")
				fs.String("token", "testtoken", "")
				return fs
			},
			verifyFlags: true,
			expectError: false,
		},
		{
			name: "valid config with specific repositories",
			setupFlags: func() *pflag.FlagSet {
				fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
				fs.String("username", "testuser", "")
				fs.String("base-url", "https://gerrit.example.com", "")
				fs.StringSlice("repo", []string{"repo1", "repo2"}, "")
				fs.String("repo-search", "", "")
				fs.String("token", "testtoken", "")
				return fs
			},
			verifyFlags: true,
			expectError: false,
		},
		{
			name: "both repo and repo-search defined - should error",
			setupFlags: func() *pflag.FlagSet {
				fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
				fs.String("username", "testuser", "")
				fs.String("base-url", "https://gerrit.example.com", "")
				fs.StringSlice("repo", []string{"repo1"}, "")
				fs.String("repo-search", "myrepo", "")
				fs.String("token", "testtoken", "")
				return fs
			},
			verifyFlags: true,
			expectError: true,
			errorMsg:    "repo and repoSearch can't be defined both",
		},
		{
			name: "both repo and repo-search defined but verifyFlags false - should not error",
			setupFlags: func() *pflag.FlagSet {
				fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
				fs.String("username", "testuser", "")
				fs.String("base-url", "https://gerrit.example.com", "")
				fs.StringSlice("repo", []string{"repo1"}, "")
				fs.String("repo-search", "myrepo", "")
				fs.String("token", "testtoken", "")
				return fs
			},
			verifyFlags: false,
			expectError: false,
		},
		{
			name: "missing username",
			setupFlags: func() *pflag.FlagSet {
				fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
				fs.String("username", "", "")
				fs.String("base-url", "https://gerrit.example.com", "")
				fs.StringSlice("repo", []string{}, "")
				fs.String("repo-search", "", "")
				fs.String("token", "testtoken", "")
				return fs
			},
			verifyFlags: true,
			expectError: true,
			errorMsg:    "no username set",
		},
		{
			name: "missing base-url",
			setupFlags: func() *pflag.FlagSet {
				fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
				fs.String("username", "testuser", "")
				fs.String("base-url", "", "")
				fs.StringSlice("repo", []string{}, "")
				fs.String("repo-search", "", "")
				fs.String("token", "testtoken", "")
				return fs
			},
			verifyFlags: true,
			expectError: true,
			errorMsg:    "no base-url set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.setupFlags()
			vc, err := createGerritClient(fs, tt.verifyFlags)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, vc)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, vc)
			}
		})
	}
}
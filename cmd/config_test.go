package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeStaticConfig(t *testing.T) {
	tests := []struct {
		name          string
		xdgConfigHome string            // relative to the temp dir, empty means unset
		files         map[string]string // config file path (relative to the temp dir) -> token
		expected      string
	}{
		{
			name:          "XDG_CONFIG_HOME takes precedence over legacy",
			xdgConfigHome: "xdg",
			files: map[string]string{
				"xdg/multi-gitter/config.yaml":   "xdg-token",
				"home/.multi-gitter/config.yaml": "legacy-token",
			},
			expected: "xdg-token",
		},
		{
			name: "default XDG location takes precedence over legacy",
			files: map[string]string{
				"home/.config/multi-gitter/config.yaml": "default-xdg-token",
				"home/.multi-gitter/config.yaml":        "legacy-token",
			},
			expected: "default-xdg-token",
		},
		{
			name:          "fallback to legacy location",
			xdgConfigHome: "xdg",
			files: map[string]string{
				"home/.multi-gitter/config.yaml": "legacy-token",
			},
			expected: "legacy-token",
		},
		{
			name:     "no config file",
			expected: "default",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()
			for path, token := range test.files {
				path = filepath.Join(tempDir, path)
				require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
				require.NoError(t, os.WriteFile(path, []byte("token: "+token+"\n"), 0644))
			}

			xdgConfigHome := ""
			if test.xdgConfigHome != "" {
				xdgConfigHome = filepath.Join(tempDir, test.xdgConfigHome)
			}
			t.Setenv("XDG_CONFIG_HOME", xdgConfigHome)
			t.Setenv("HOME", filepath.Join(tempDir, "home"))

			cmd := &cobra.Command{}
			cmd.Flags().String("token", "default", "")
			require.NoError(t, initializeStaticConfig(cmd))

			val, err := cmd.Flags().GetString("token")
			require.NoError(t, err)
			assert.Equal(t, test.expected, val)
		})
	}
}

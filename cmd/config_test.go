package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeStaticConfig_XDGConfigHome(t *testing.T) {
	tempDir := t.TempDir()
	xdgDir := filepath.Join(tempDir, "xdg")
	homeDir := filepath.Join(tempDir, "home")

	xdgConfigDir := filepath.Join(xdgDir, "multi-gitter")
	legacyConfigDir := filepath.Join(homeDir, ".multi-gitter")

	require.NoError(t, os.MkdirAll(xdgConfigDir, 0755))
	require.NoError(t, os.MkdirAll(legacyConfigDir, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(xdgConfigDir, "config.yaml"), []byte("token: xdg-token\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(legacyConfigDir, "config.yaml"), []byte("token: legacy-token\n"), 0644))

	t.Setenv("XDG_CONFIG_HOME", xdgDir)
	t.Setenv("HOME", homeDir)

	cmd := &cobra.Command{}
	cmd.Flags().String("token", "", "")

	err := initializeStaticConfig(cmd)
	require.NoError(t, err)

	val, err := cmd.Flags().GetString("token")
	require.NoError(t, err)
	assert.Equal(t, "xdg-token", val)
}

func TestInitializeStaticConfig_DefaultXDG(t *testing.T) {
	tempDir := t.TempDir()
	homeDir := filepath.Join(tempDir, "home")

	defaultXdgConfigDir := filepath.Join(homeDir, ".config", "multi-gitter")
	legacyConfigDir := filepath.Join(homeDir, ".multi-gitter")

	require.NoError(t, os.MkdirAll(defaultXdgConfigDir, 0755))
	require.NoError(t, os.MkdirAll(legacyConfigDir, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(defaultXdgConfigDir, "config.yaml"), []byte("token: default-xdg-token\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(legacyConfigDir, "config.yaml"), []byte("token: legacy-token\n"), 0644))

	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", homeDir)

	cmd := &cobra.Command{}
	cmd.Flags().String("token", "", "")

	err := initializeStaticConfig(cmd)
	require.NoError(t, err)

	val, err := cmd.Flags().GetString("token")
	require.NoError(t, err)
	assert.Equal(t, "default-xdg-token", val)
}

func TestInitializeStaticConfig_LegacyFallback(t *testing.T) {
	t.Run("with XDG_CONFIG_HOME unset", func(t *testing.T) {
		tempDir := t.TempDir()
		homeDir := filepath.Join(tempDir, "home")
		legacyConfigDir := filepath.Join(homeDir, ".multi-gitter")

		require.NoError(t, os.MkdirAll(legacyConfigDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(legacyConfigDir, "config.yaml"), []byte("token: legacy-token\n"), 0644))

		t.Setenv("XDG_CONFIG_HOME", "")
		t.Setenv("HOME", homeDir)

		cmd := &cobra.Command{}
		cmd.Flags().String("token", "", "")

		err := initializeStaticConfig(cmd)
		require.NoError(t, err)

		val, err := cmd.Flags().GetString("token")
		require.NoError(t, err)
		assert.Equal(t, "legacy-token", val)
	})

	t.Run("with XDG_CONFIG_HOME set to empty directory", func(t *testing.T) {
		tempDir := t.TempDir()
		xdgDir := filepath.Join(tempDir, "empty-xdg")
		homeDir := filepath.Join(tempDir, "home")
		legacyConfigDir := filepath.Join(homeDir, ".multi-gitter")

		require.NoError(t, os.MkdirAll(xdgDir, 0755))
		require.NoError(t, os.MkdirAll(legacyConfigDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(legacyConfigDir, "config.yaml"), []byte("token: legacy-token\n"), 0644))

		t.Setenv("XDG_CONFIG_HOME", xdgDir)
		t.Setenv("HOME", homeDir)

		cmd := &cobra.Command{}
		cmd.Flags().String("token", "", "")

		err := initializeStaticConfig(cmd)
		require.NoError(t, err)

		val, err := cmd.Flags().GetString("token")
		require.NoError(t, err)
		assert.Equal(t, "legacy-token", val)
	})
}

func TestInitializeStaticConfig_NoConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, "xdg"))
	t.Setenv("HOME", filepath.Join(tempDir, "home"))

	cmd := &cobra.Command{}
	cmd.Flags().String("token", "default", "")

	err := initializeStaticConfig(cmd)
	require.NoError(t, err)

	val, err := cmd.Flags().GetString("token")
	require.NoError(t, err)
	assert.Equal(t, "default", val)
}

func TestInitializeStaticConfig_FlagPrecedence(t *testing.T) {
	tempDir := t.TempDir()
	homeDir := filepath.Join(tempDir, "home")
	configDir := filepath.Join(homeDir, ".config", "multi-gitter")

	require.NoError(t, os.MkdirAll(configDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte("token: config-token\n"), 0644))

	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", homeDir)

	cmd := &cobra.Command{}
	cmd.Flags().String("token", "", "")
	require.NoError(t, cmd.Flags().Set("token", "cli-token"))

	err := initializeStaticConfig(cmd)
	require.NoError(t, err)

	val, err := cmd.Flags().GetString("token")
	require.NoError(t, err)
	assert.Equal(t, "cli-token", val)
}

func TestInitializeConfig_DynamicOverStatic(t *testing.T) {
	tempDir := t.TempDir()
	homeDir := filepath.Join(tempDir, "home")
	staticConfigDir := filepath.Join(homeDir, ".config", "multi-gitter")

	require.NoError(t, os.MkdirAll(staticConfigDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(staticConfigDir, "config.yaml"), []byte("token: static-token\nlog-level: debug\n"), 0644))

	dynamicConfigFile := filepath.Join(tempDir, "dynamic.yaml")
	require.NoError(t, os.WriteFile(dynamicConfigFile, []byte("token: dynamic-token\n"), 0644))

	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", homeDir)

	cmd := &cobra.Command{}
	configureConfig(cmd)
	cmd.Flags().String("token", "", "")
	cmd.Flags().String("log-level", "", "")
	require.NoError(t, cmd.Flags().Set("config", dynamicConfigFile))

	err := initializeConfig(cmd)
	require.NoError(t, err)

	tokenVal, err := cmd.Flags().GetString("token")
	require.NoError(t, err)
	assert.Equal(t, "dynamic-token", tokenVal)

	logLevelVal, err := cmd.Flags().GetString("log-level")
	require.NoError(t, err)
	assert.Equal(t, "debug", logLevelVal)
}

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_MissingFileCreatesDefault(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dl", "config.toml")
	dataDir := filepath.Join(tmpDir, "data", "dl")

	cfg, err := loadConfig(configPath, dataDir)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, dataDir, cfg.DataDir)

	// Verify file was written
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "data_dir")

	// Verify data dir was created
	info, err := os.Stat(dataDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestLoadConfig_ExistingFile(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")
	customDir := filepath.Join(tmpDir, "custom", "path")
	fallbackDir := filepath.Join(tmpDir, "fallback", "dl")

	content := "data_dir = \"" + customDir + "\"\n"
	require.NoError(t, os.WriteFile(configPath, []byte(content), 0o600))

	cfg, err := loadConfig(configPath, fallbackDir)
	require.NoError(t, err)
	assert.Equal(t, customDir, cfg.DataDir)
}

func TestLoadConfig_EmptyDataDirFallsBack(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")
	fallbackDir := filepath.Join(tmpDir, "fallback", "dl")

	content := "data_dir = \"\"\n"
	require.NoError(t, os.WriteFile(configPath, []byte(content), 0o600))

	cfg, err := loadConfig(configPath, fallbackDir)
	require.NoError(t, err)
	assert.Equal(t, fallbackDir, cfg.DataDir)
}

func TestLoadConfig_InvalidTOML(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")
	dataDir := filepath.Join(tmpDir, "data", "dl")

	require.NoError(t, os.WriteFile(configPath, []byte("not valid toml [[["), 0o600))

	_, err := loadConfig(configPath, dataDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode config")
}

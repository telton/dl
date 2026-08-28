package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_VersionFlag(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	stdin := os.Stdin

	err := run(&stdout, stdin, []string{"dl", "--version"})
	require.NoError(t, err)
	assert.Contains(t, stdout.String(), "dl")
}

func TestRun_PipeMode(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")
	dataDir := filepath.Join(tmpDir, "data")

	require.NoError(t, os.WriteFile(configPath, []byte("data_dir = \""+dataDir+"\"\n"), 0o600))

	// Create piped stdin
	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	go func() {
		writer.WriteString("pipe test entry")
		writer.Close()
	}()

	var stdout bytes.Buffer
	err = run(&stdout, reader, []string{"dl", "--config", configPath})
	require.NoError(t, err)

	// Verify
	content, err := readTodayEntries(dataDir)
	require.NoError(t, err)
	assert.Contains(t, content, "pipe test entry")
}

func TestRun_EmptyPipe(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")
	dataDir := filepath.Join(tmpDir, "data")

	require.NoError(t, os.WriteFile(configPath, []byte("data_dir = \""+dataDir+"\"\n"), 0o600))

	// Create piped stdin with empty content
	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	go func() {
		writer.Close()
	}()

	var stdout bytes.Buffer
	err = run(&stdout, reader, []string{"dl", "--config", configPath})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty entry")
}

func TestRun_ProjectSubdir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")
	dataDir := filepath.Join(tmpDir, "data")
	projectDir := filepath.Join(dataDir, "myproject")

	require.NoError(t, os.WriteFile(configPath, []byte("data_dir = \""+dataDir+"\"\n"), 0o600))

	// Create piped stdin
	reader, writer, err := os.Pipe()
	require.NoError(t, err)

	go func() {
		writer.WriteString("project specific entry")
		writer.Close()
	}()

	var stdout bytes.Buffer
	err = run(&stdout, reader, []string{"dl", "--config", configPath, "--project", "myproject"})
	require.NoError(t, err)

	// Verify it wrote to the project subdir
	content, err := readTodayEntries(projectDir)
	require.NoError(t, err)
	assert.Contains(t, content, "project specific entry")

	// Verify it did NOT write to the root data dir
	rootContent, err := readTodayEntries(dataDir)
	require.NoError(t, err)
	assert.Empty(t, rootContent)
}

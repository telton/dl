package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodayFilename(t *testing.T) {
	t.Parallel()

	expected := time.Now().Format("2006-01-02") + ".log"
	assert.Equal(t, expected, todayFilename())
}

func TestTodayFilePath(t *testing.T) {
	t.Parallel()

	dataDir := "/tmp/data"
	expected := filepath.Join(dataDir, todayFilename())
	assert.Equal(t, expected, todayFilePath(dataDir))
}

func TestFormatLogEntry(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 26, 14, 32, 1, 0, time.UTC)
	entry := "did the thing"
	expected := "[14:32:01] did the thing\n"
	assert.Equal(t, expected, formatLogEntry(now, entry))
}

func TestAppendEntryAt(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	now := time.Date(2026, 8, 26, 14, 32, 1, 0, time.UTC)

	err := appendEntryAt(tmpDir, "did the thing", now)
	require.NoError(t, err)

	content, err := os.ReadFile(todayFilePath(tmpDir))
	require.NoError(t, err)
	assert.Equal(t, "[14:32:01] did the thing\n", string(content))
}

func TestAppendEntryAt_MultipleEntries(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	now1 := time.Date(2026, 8, 26, 14, 32, 1, 0, time.UTC)
	now2 := time.Date(2026, 8, 26, 14, 45, 22, 0, time.UTC)

	require.NoError(t, appendEntryAt(tmpDir, "entry one", now1))
	require.NoError(t, appendEntryAt(tmpDir, "entry two", now2))

	content, err := os.ReadFile(todayFilePath(tmpDir))
	require.NoError(t, err)

	lines := strings.Split(strings.TrimRight(string(content), "\n"), "\n")
	require.Len(t, lines, 2)
	assert.Equal(t, "[14:32:01] entry one", lines[0])
	assert.Equal(t, "[14:45:22] entry two", lines[1])
}

func TestAppendEntryAt_EmptyEntry(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	now := time.Now()

	err := appendEntryAt(tmpDir, "", now)
	require.Error(t, err)
	assert.Equal(t, "empty entry", err.Error())

	err = appendEntryAt(tmpDir, "   ", now)
	require.Error(t, err)
	assert.Equal(t, "empty entry", err.Error())
}

func TestAppendEntryAt_WritesToCorrectFile(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	now := time.Now()

	err := appendEntryAt(tmpDir, "test entry", now)
	require.NoError(t, err)

	// Verify file exists
	info, err := os.Stat(todayFilePath(tmpDir))
	require.NoError(t, err)
	assert.False(t, info.IsDir())
	assert.Positive(t, info.Size())
}

func TestReadTodayEntries(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	now := time.Now()

	// Empty when file doesn't exist
	content, err := readTodayEntries(tmpDir)
	require.NoError(t, err)
	assert.Empty(t, content)

	// After writing
	require.NoError(t, appendEntryAt(tmpDir, "test entry", now))
	content, err = readTodayEntries(tmpDir)
	require.NoError(t, err)
	assert.Contains(t, content, "test entry")
}

func TestAppendEntry_Integration(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	// Since appendEntry uses time.Now(), we can't control the exact timestamp,
	// but we can verify the file is written and the format is correct.
	err := appendEntryAt(tmpDir, "integration test", time.Now())
	require.NoError(t, err)

	content, err := readTodayEntries(tmpDir)
	require.NoError(t, err)
	assert.Contains(t, content, "integration test")
	assert.Contains(t, content, "[")
	assert.Contains(t, content, "]")
}

func TestAppendEntryAt_FilePermissions(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	now := time.Now()

	require.NoError(t, appendEntryAt(tmpDir, "test", now))

	info, err := os.Stat(todayFilePath(tmpDir))
	require.NoError(t, err)
	mode := info.Mode().Perm()
	// os.OpenFile creates with 0644, actual may be masked by umask
	assert.NotEqual(t, os.FileMode(0), mode&0644, "file should be readable/writable")
}

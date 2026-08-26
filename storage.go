package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func todayFilename() string {
	return time.Now().Format("2006-01-02") + ".log"
}

func todayFilePath(dataDir string) string {
	return filepath.Join(dataDir, todayFilename())
}

func formatLogEntry(now time.Time, entry string) string {
	return fmt.Sprintf("[%s] %s\n", now.Format("15:04:05"), entry)
}

func appendEntryAt(dataDir string, entry string, now time.Time) error {
	path := todayFilePath(dataDir)
	entry = strings.TrimSpace(entry)
	if entry == "" {
		return errors.New("empty entry")
	}

	line := formatLogEntry(now, entry)

	//nolint:gosec // data path is user-configured
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer logFile.Close()

	_, err = logFile.WriteString(line)
	if err != nil {
		return fmt.Errorf("write entry: %w", err)
	}

	return nil
}

func appendEntry(dataDir, entry string) error {
	return appendEntryAt(dataDir, entry, time.Now())
}

func readTodayEntries(dataDir string) (string, error) {
	path := todayFilePath(dataDir)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}

		return "", fmt.Errorf("read entries: %w", err)
	}

	return string(data), nil
}

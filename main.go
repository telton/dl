package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func run(stdout io.Writer, stdin *os.File, args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	showVersion := fs.Bool("version", false, "print version and exit")
	configPath := fs.String("config", "", "path to config file")
	project := fs.String("project", "", "project name (subdirectory under data_dir)")

	if err := fs.Parse(args[1:]); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	if *showVersion {
		printVersion(stdout)

		return nil
	}

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	dataDir := cfg.DataDir
	if *project != "" {
		dataDir = filepath.Join(dataDir, *project)
		if err := os.MkdirAll(dataDir, 0750); err != nil {
			return fmt.Errorf("create project dir: %w", err)
		}
	}

	// Detect if stdin is a pipe or redirected
	stat, err := stdin.Stat()
	if err != nil {
		return fmt.Errorf("stat stdin: %w", err)
	}

	if stat.Mode()&os.ModeCharDevice == 0 {
		data, err := io.ReadAll(stdin)
		if err != nil {
			return fmt.Errorf("read stdin: %w", err)
		}

		entry := strings.TrimSpace(string(data))
		if entry == "" {
			return errors.New("empty entry")
		}

		if err := appendEntry(dataDir, entry); err != nil {
			return fmt.Errorf("append entry: %w", err)
		}

		return nil
	}

	// Print existing entries
	todayPath := todayFilePath(dataDir)
	data, err := os.ReadFile(todayPath)
	if err == nil && len(data) > 0 {
		fmt.Fprintln(stdout, string(data))
	}

	fmt.Fprintln(stdout, "dl — ctrl+c to quit")

	scanner := bufio.NewScanner(stdin)
	for {
		fmt.Fprint(stdout, "> ")

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if err := appendEntry(dataDir, line); err != nil {
			fmt.Fprintf(stdout, "Error: %v\n", err)
		}
	}

	return nil
}

func main() {
	if err := run(os.Stdout, os.Stdin, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Package dl provides a minimal devlog CLI.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

type Config struct {
	DataDir string `toml:"data_dir"`
}

func defaultConfigPath() string {
	return filepath.Join(xdg.ConfigHome, "dl", "config.toml")
}

func loadConfig(configPath, dataDir string) (*Config, error) {
	cfg := &Config{
		DataDir: dataDir,
	}

	//nolint:gosec // user-configurable config path
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(configPath), 0750); err != nil {
				return nil, fmt.Errorf("create config dir: %w", err)
			}

			if err := os.MkdirAll(cfg.DataDir, 0750); err != nil {
				return nil, fmt.Errorf("create data dir: %w", err)
			}

			//nolint:gosec // user-configurable config path
			outFile, err := os.Create(configPath)
			if err != nil {
				return nil, fmt.Errorf("create config file: %w", err)
			}
			defer outFile.Close()

			enc := toml.NewEncoder(outFile)
			if err := enc.Encode(cfg); err != nil {
				return nil, fmt.Errorf("encode default config: %w", err)
			}

			return cfg, nil
		}

		return nil, fmt.Errorf("read config: %w", err)
	}

	if _, err := toml.Decode(string(data), cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	if cfg.DataDir == "" {
		cfg.DataDir = dataDir
	}

	if err := os.MkdirAll(cfg.DataDir, 0750); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	return cfg, nil
}

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = defaultConfigPath()
	}

	dataDir := filepath.Join(xdg.DataHome, "dl")

	return loadConfig(configPath, dataDir)
}

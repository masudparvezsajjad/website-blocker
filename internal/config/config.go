package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	Enabled        bool     `json:"enabled"`
	BlockPagePort  int      `json:"block_page_port"`
	BlockedDomains []string `json:"blocked_domains"`
	RedirectIPv4   string   `json:"redirect_ipv4"`
	RedirectIPv6   string   `json:"redirect_ipv6"`
}

const (
	AppDirName = "AdultBlocker"
	ConfigFile = "config.json"
)

func AppSupportDir() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		// fallback to home directory if UserConfigDir is unavailable
		home, err := os.UserHomeDir()
		if err != nil {
			return filepath.Join("/", AppDirName)
		}
		return filepath.Join(home, ".config", AppDirName)
	}
	return filepath.Join(dir, AppDirName)
}

func ConfigPath() string {
	return filepath.Join(AppSupportDir(), ConfigFile)
}

func DefaultConfig() *Config {
	return &Config{
		Enabled:       false,
		BlockPagePort: 8088,
		BlockedDomains: []string{
			"pornhub.com",
			"www.pornhub.com",
			"xvideos.com",
			"www.xvideos.com",
			"xnxx.com",
			"www.xnxx.com",
		},
		RedirectIPv4: "127.0.0.1",
		RedirectIPv6: "::1",
	}
}

func EnsureDir() error {
	return os.MkdirAll(AppSupportDir(), 0755)
}

func Load() (*Config, error) {
	path := ConfigPath()
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfg := DefaultConfig()
			if err := Save(cfg); err != nil {
				return nil, err
			}
			return cfg, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	if cfg.BlockPagePort == 0 {
		cfg.BlockPagePort = 8088
	}
	if cfg.RedirectIPv4 == "" {
		cfg.RedirectIPv4 = "127.0.0.1"
	}
	if cfg.RedirectIPv6 == "" {
		cfg.RedirectIPv6 = "::1"
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	if err := EnsureDir(); err != nil {
		return err
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigPath(), b, 0644)
}

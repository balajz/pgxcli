// Package config provides utilities
// for loading and managing configuration settings for the application.
//
// Configuration is loaded from an embedded default configuration file.
// merged with a user configuration file located in the user's config directory.
package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/balajz/pgxcli/internal/perrors"
	"github.com/spf13/viper"
)

const (
	// Default indicates that a path-based option should use its default location.
	Default  = "default"
	filename = "config.toml"
	appName  = "pgxcli"
)

// Config represents the top-level application configuration.
type Config struct {
	Main  MainConfig  `mapstructure:"main" toml:"main"`
	Table TableConfig `mapstructure:"table" toml:"table"`
}

// MainConfig contains general CLI and session settings.
type MainConfig struct {
	Prompt      string               `mapstructure:"prompt" toml:"prompt"`
	Style       SyntaxHighlightStyle `mapstructure:"style" toml:"style"`
	HistoryFile string               `mapstructure:"history_file" toml:"history_file"`
	LogFile     string               `mapstructure:"log_file" toml:"log_file"`
	Pager       string               `mapstructure:"pager" toml:"pager"`
	OnError     OnErrorAction        `mapstructure:"on_error" toml:"on_error"`
}

// TableConfig contains output table rendering settings.
type TableConfig struct {
	Style TableStyle       `mapstructure:"style" toml:"style"`
	Color TableColorConfig `mapstructure:"color" toml:"color"`
}

// TableColorConfig contains color settings for table elements.
type TableColorConfig struct {
	Header  TableColor `mapstructure:"header" toml:"header"`
	Column  TableColor `mapstructure:"column" toml:"column"`
	Caption TableColor `mapstructure:"caption" toml:"caption"`
}

// Load reads the embedded default configuration and merges with user configuration.
func Load() (*Config, error) {
	userPath, err := UserConfigPath()
	if err != nil {
		return nil, err
	}
	if err := ensureUserConfig(userPath); err != nil {
		return nil, err
	}

	defaultV := viper.New()
	defaultV.SetConfigType("toml")
	if err := defaultV.ReadConfig(bytes.NewReader(defaultConfigFile)); err != nil {
		return nil, fmt.Errorf("read default config: %w", err)
	}

	userV := viper.New()
	userV.SetConfigFile(userPath)
	if err := defaultV.ReadConfig(bytes.NewReader(defaultConfigFile)); err != nil {
		return nil, perrors.Wrap(
			err,
			perrors.WithMessage("failed to read config"),
			perrors.WithDetails(
				"type", "user",
			),
		)
	}
	// user settings land on top of default settings
	if err := defaultV.MergeConfigMap(userV.AllSettings()); err != nil {
		return nil, perrors.Wrap(
			err,
			perrors.WithMessage("failed to merge cofig"),
		)
	}

	var cfg Config
	if err := defaultV.Unmarshal(&cfg); err != nil {
		return nil, perrors.Wrap(
			err,
			perrors.WithMessage("failed to unmarshal config"),
			perrors.WithDetails(
				"type", "merged",
			),
		)
	}

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}
	return &cfg, nil
}

// GetDefaultConfig returns the default configuration embedded in the binary.
func GetDefaultConfig() (*Config, error) {
	defaultV := viper.New()
	defaultV.SetConfigType("toml")
	if err := defaultV.ReadConfig(bytes.NewReader(defaultConfigFile)); err != nil {
		return nil, perrors.Wrap(
			err,
			perrors.WithMessage("failed to read config"),
			perrors.WithDetails(
				"type", "default",
			),
		)
	}

	var cfg Config
	if err := defaultV.Unmarshal(&cfg); err != nil {
		return nil, perrors.Wrap(
			err,
			perrors.WithMessage("failed to unmarshal config"),
			perrors.WithDetails(
				"type", "default",
			),
		)
	}
	return &cfg, nil
}

// UserConfigPath returns the user config path (for example: ~/.config/pgxcli/config.toml).
func UserConfigPath() (string, error) {
	if userdir := os.Getenv("XDG_CONFIG_HOME"); userdir != "" {
		return filepath.Join(userdir, appName, filename), nil
	}

	userdir, err := os.UserConfigDir()
	if err != nil {
		return "", perrors.Wrap(
			err,
			perrors.WithMessage("failed to get the user config directory"),
		)
	}
	return filepath.Join(userdir, appName, filename), nil
}

// ensureUserConfig write embed on firt run
func ensureUserConfig(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return perrors.Wrap(
			err,
			perrors.WithMessage("failed to create config directory"),
			perrors.WithDetails(
				"path", path,
			),
		)
	}

	if err := os.WriteFile(path, defaultConfigFile, 0o644); err != nil {
		return perrors.Wrap(
			err,
			perrors.WithMessage("failed to write default config"),
			perrors.WithDetails(
				"path", path,
			),
		)
	}
	return nil
}

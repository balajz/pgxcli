package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate_Success(t *testing.T) {
	cfg := Config{
		Main: main{
			Prompt:      "test> ",
			Style:       "monokai",
			HistoryFile: "default",
			LogFile:     "default",
		},
	}

	err := validate(cfg)
	assert.NoError(t, err)
}

func TestValidate_MultipleErrors(t *testing.T) {
	cfg := Config{
		Main: main{
			Prompt:      "",
			Style:       "",
			HistoryFile: "",
			LogFile:     "",
		},
	}

	err := validate(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "prompt must not be empty")
	assert.Contains(t, err.Error(), "style must not be empty")
	assert.Contains(t, err.Error(), "history file path must not be empty")
	assert.Contains(t, err.Error(), "log file path must not be empty")
}

func TestLoad_ValidationFailsOnEmptyPrompt(t *testing.T) {
	tempDir := t.TempDir()
	originalConfigDir := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if originalConfigDir != "" {
			os.Setenv("XDG_CONFIG_HOME", originalConfigDir)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
	}()

	os.Setenv("XDG_CONFIG_HOME", tempDir)

	userConfigPath := filepath.Join(tempDir, appName, filename)
	require.NoError(t, os.MkdirAll(filepath.Dir(userConfigPath), 0o700))

	userConfig := `[main]
prompt = ""
style = "monokai"
history_file = "default"
log_file = "default"
`
	require.NoError(t, os.WriteFile(userConfigPath, []byte(userConfig), 0o644))

	_, err := Load()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validate config")
	assert.Contains(t, err.Error(), "prompt must not be empty")
}

func TestLoad_ValidationFailsOnEmptyStyle(t *testing.T) {
	tempDir := t.TempDir()
	originalConfigDir := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		if originalConfigDir != "" {
			os.Setenv("XDG_CONFIG_HOME", originalConfigDir)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
	}()

	os.Setenv("XDG_CONFIG_HOME", tempDir)

	userConfigPath := filepath.Join(tempDir, appName, filename)
	require.NoError(t, os.MkdirAll(filepath.Dir(userConfigPath), 0o700))

	userConfig := `[main]
prompt = "test> "
style = ""
history_file = "default"
log_file = "default"
`
	require.NoError(t, os.WriteFile(userConfigPath, []byte(userConfig), 0o644))

	_, err := Load()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validate config")
	assert.Contains(t, err.Error(), "style must not be empty")
}

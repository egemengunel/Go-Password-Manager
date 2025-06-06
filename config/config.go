package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// Config holds application configuration
type Config struct {
	VaultPath    string
	ConfigDir    string
	DefaultVault string
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()

	var configDir string
	switch runtime.GOOS {
	case "windows":
		configDir = filepath.Join(os.Getenv("APPDATA"), "gopassman")
	case "darwin":
		configDir = filepath.Join(homeDir, ".config", "gopassman")
	default: // linux and others
		configDir = filepath.Join(homeDir, ".config", "gopassman")
	}

	return &Config{
		VaultPath:    filepath.Join(configDir, "vault.gpv"),
		ConfigDir:    configDir,
		DefaultVault: "default",
	}
}

// EnsureConfigDir creates the configuration directory if it doesn't exist
func (c *Config) EnsureConfigDir() error {
	return os.MkdirAll(c.ConfigDir, 0700)
}

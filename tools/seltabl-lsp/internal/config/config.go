// Package config provides a set of functions for managing the
// configuration folder for the lanaguage server.
package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
)

// Config is the configuration for the server
type Config struct {
	// ConfigPath is the path to the configuration folder
	ConfigPath string
}

// CreateConfigDir creates a new config directory and returns a config
func CreateConfigDir() (*Config, error) {
	path, err := homedir.Expand("~/.config/seltabl-lsp/")
	if err != nil {
		return nil, fmt.Errorf("failed to expand home directory: %w", err)
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		if os.IsExist(err) {
			return &Config{ConfigPath: path}, nil
		}
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	return &Config{
		ConfigPath: path,
	}, nil
}

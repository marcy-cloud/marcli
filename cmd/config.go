package cmd

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration - keeping track of our version and builds! 💕
type Config struct {
	Version         string `yaml:"version"`         // Our cute version number! ✨
	Build           int    `yaml:"build"`           // Build counter - we're so organized! 🎀
	ExitAfterCommand bool  `yaml:"exitAfterCommand"` // Whether to exit TUI after running a command or return to menu
}

const configFile = "config.yml" // Where we keep our config, obviously! 💖

// LoadConfig loads the configuration from config.yml - so reliable! ✨
func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to config.yml - keeping everything organized! 💅
func SaveConfig(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// IncrementBuild increments the build number and saves the config - we're so organized! 🎀
func IncrementBuild() error {
	config, err := LoadConfig()
	if err != nil {
		// If config doesn't exist, create a cute default one! ✨
		config = &Config{
			Version: "0.1.0",
			Build:   0,
		}
	}

	config.Build++
	return SaveConfig(config)
}

// GetVersion returns the version string - formatted so nicely! 💖
func GetVersion() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s (build %d)", config.Version, config.Build), nil
}


package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	CommandSets map[string]CommandSet `yaml:"commands"`
	Global      GlobalConfig          `yaml:"global"`
}

// CommandSet represents a group of related commands
type CommandSet struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Commands    []string `yaml:"commands"`
	Dir         string   `yaml:"dir"`
	AutoRestart bool     `yaml:"auto_restart"`
	Env         []string `yaml:"env"`
}

// GlobalConfig represents global settings
type GlobalConfig struct {
	LogFile     string `yaml:"log_file"`
	MaxOutput   int    `yaml:"max_output_lines"`
	RefreshRate int    `yaml:"refresh_rate_ms"`
}

// Load loads configuration from a file
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if config.Global.RefreshRate == 0 {
		config.Global.RefreshRate = 100
	}
	if config.Global.MaxOutput == 0 {
		config.Global.MaxOutput = 1000
	}

	return &config, nil
}

// Save saves configuration to a file
func (c *Config) Save(filename string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		CommandSets: map[string]CommandSet{
			"example": {
				Name:        "Example",
				Description: "Example command set",
				Commands:    []string{"echo 'Hello World'", "sleep 5"},
				Dir:         ".",
				AutoRestart: false,
			},
		},
		Global: GlobalConfig{
			LogFile:     "cmdpool.log",
			MaxOutput:   1000,
			RefreshRate: 100,
		},
	}
} 
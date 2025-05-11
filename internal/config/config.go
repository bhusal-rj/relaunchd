package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name    string        `yaml:"name"`
	Command CommandConfig `yaml:"command"`
	Watch   WatchConfig   `yaml:"watch"`
	Process ProcessConfig `yaml:"process"`
}

type CommandConfig struct {
	Exec string            `yaml:"exec"`
	Cwd  string            `yaml:"cwd"`
	Env  map[string]string `yaml:"env"`
}

type WatchConfig struct {
	Paths   []string `yaml:"paths"`
	Exclude []string `yaml:"exclude"`
}

type ProcessConfig struct {
	Background   bool `yaml:"background"`
	MaxRestarts  int  `yaml:"max_restarts"`
	RestartDelay int  `yaml:"restart_delay"`
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = "config.yaml"
	}

	// Get the abssoulute path of the config file
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.New("failed to get absolute path of config file")
	}

	//Check if the config file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, errors.New("config file does not exist")
	}

	// Read the config file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %v", err)
	}

	// Apply defaults and validate the config

	if err := applyDefaults(config); err != nil {
		return nil, fmt.Errorf("failed to apply defaults: %v", err)
	}

	if err := validate(config); err != nil {
		return nil, fmt.Errorf("failed to validate config: %v", err)
	}

	return config, nil
}

// applyDefaults sets default values for the config fields
func applyDefaults(config *Config) error {
	if config.Name == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory:", err)
			return fmt.Errorf("failed to get current working directory: %v", err)
		}
		config.Name = filepath.Base(cwd)
	}

	if config.Command.Cwd == "" {
		config.Command.Cwd = "."
	}

	return nil
}

// validate checks if the configuration is vlaid
func validate(config *Config) error {
	if config.Command.Exec == "" {
		return fmt.Errorf("command exec is required")
	}
	return nil
}

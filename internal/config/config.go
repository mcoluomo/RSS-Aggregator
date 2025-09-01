package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	var config Config

	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("UseHomeDir error: %w", err)
	}
	configPath := filepath.Join(homeDir, ".gatorconfig.json")

	return configPath, nil
}

func (config *Config) SetUser(userName string) ([]byte, error) {
	config.Current_user_name = userName

	data, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to encode JSON: %w", err)
	}
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		return nil, fmt.Errorf("failed to write to file %w", err)
	}

	return nil, nil
}

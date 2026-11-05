package config

import (
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DatabaseURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	fullPath := filepath.Join(home, configFileName)
	return fullPath, nil
}

func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	
	contents, err := os.Open(fullPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}
	defer contents.Close()

	decoder := json.NewDecoder(contents)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}


func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}


func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}


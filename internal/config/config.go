package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/robinskaba/roge/internal/roblox"
)

type Config struct {
	ApiKey     string             `json:"api_key"`
	AuthorId   string             `json:"author_id"`
	AuthorType roblox.CreatorType `json:"author_type"`
}

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(configDir, "roge")

	err = os.MkdirAll(appDir, 0755) // read+write+run permissions
	if err != nil {
		return "", err
	}

	return filepath.Join(appDir, "config.json"), nil
}

func SaveConfig(cfg Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600) // sets permission to only current user
}

func LoadConfig() (Config, error) {
	var cfg Config

	path, err := getConfigPath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// no config set up yet, use default
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

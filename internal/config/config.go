package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DataFilePath string
}

func GetDefaultConfig() (Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, err
	}

	appDir := filepath.Join(configDir, "mindgo")
	fullPath := filepath.Join(appDir, "moods.json")

	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		return Config{}, err
	}

	return Config{DataFilePath: fullPath}, nil
}

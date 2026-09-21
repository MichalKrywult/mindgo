package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultConfigFileName = "config.json"
	defaultDataFileName   = "moods.json"
	defaultExportFileName = "moods.csv"
	appDirectoryName      = "mindgo"

	defaultConfig = `{"configPath":"moods.json"}`
)

func EnsureConfigFile() error {
	configPath, err := buildPath("", defaultConfigFileName)
	if err != nil {
		return fmt.Errorf("failed to get user config directory: %w", err)
	}

	_, err = os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(configPath), 0755)
			if err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}

			err = os.WriteFile(configPath, []byte(defaultConfig), 0644)
			if err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}
			return nil
		}

		return fmt.Errorf("failed to read file %s: %w", configPath, err)
	}

	return nil
}


func buildPath(file, defaultFileName string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	var fullPath string
	if file != "" {
		fullPath = filepath.Join(configDir, appDirectoryName, file)
	} else {
		fullPath = filepath.Join(configDir, appDirectoryName, defaultFileName)
	}

	return fullPath, nil

}

func BuildDataPath(file string) (string, error) {
	return buildPath(file, defaultDataFileName)
}

func BuildExportPath() (string, error) {
	return buildPath("", defaultExportFileName)
}

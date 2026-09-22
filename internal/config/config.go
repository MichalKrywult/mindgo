package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultConfigFileName = "config.json"
	defaultDataFileName   = "moods.json"
	defaultExportFileName = "moods.csv"
	appDirectoryName      = "mindgo"

	defaultConfig = `{"dataPath":"moods.json"}`
)

type Config struct {
	DataPath string `json:"dataPath"`
}

func EnsureConfigFile(configPath string) error {
	_, err := os.Stat(configPath)
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

func LoadConfigFromFile(configPath string) (Config, error) {

	data, err := os.ReadFile(configPath) //the path is aready ensured
	if err != nil {
		return Config{}, fmt.Errorf("failed to read file %s: %w", configPath, err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config from %s: %w", configPath, err)
	}

	return config, nil
}

func SaveConfigToFile(configPath string, config Config) error {
	// data, err := json.Marshal() would write everything in one line
	// while technically correct, it's unreadable
	data, err := json.MarshalIndent(config, "", "  ") // writes data readable to humans
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", configPath, err)
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

func BuildConfigPath() (string, error) {
	return buildPath("", defaultConfigFileName)
}

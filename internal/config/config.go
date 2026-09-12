package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultDataFileName   = "moods.json"
	defaultExportFileName = "moods.csv"
	appDirectoryName      = "mindgo"
)

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

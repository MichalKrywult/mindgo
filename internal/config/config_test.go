package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildDataPath(t *testing.T) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("failed to get user config directory: %v", err)
	}

	tests := []struct {
		name         string
		file         string
		expectedPath string
	}{
		{
			name:         "default path",
			file:         "",
			expectedPath: filepath.Join(configDir, appDirectoryName, defaultDataFileName),
		},
		{
			name:         "custom file",
			file:         "custom.json",
			expectedPath: filepath.Join(configDir, appDirectoryName, "custom.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath, err := BuildDataPath(tt.file)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if fullPath != tt.expectedPath {
				t.Errorf("expected path: %s, received path: %s", tt.expectedPath, fullPath)
			}

		})
	}
}

func TestBuildExportPath(t *testing.T) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("failed to get user config directory: %v", err)
	}

	tests := []struct {
		name         string
		expectedPath string
	}{
		{
			name:         "default file",
			expectedPath: filepath.Join(configDir, appDirectoryName, defaultExportFileName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath, err := BuildExportPath()
			if err != nil {
				t.Errorf("unexpected error: %v", err)

			}

			if fullPath != tt.expectedPath {
				t.Errorf("expected path: %s, received path: %s", tt.expectedPath, fullPath)
			}
		})
	}
}

func TestBuildConfigPath(t *testing.T) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("failed to get user config directory: %v", err)
	}

	tests := []struct {
		name         string
		expectedPath string
	}{
		{
			name:         "default file",
			expectedPath: filepath.Join(configDir, appDirectoryName, defaultConfigFileName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath, err := BuildConfigPath()
			if err != nil {
				t.Errorf("unexpected error: %v", err)

			}

			if fullPath != tt.expectedPath {
				t.Errorf("expected path: %s, received path: %s", tt.expectedPath, fullPath)
			}
		})
	}
}

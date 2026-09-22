package config

import (
	"encoding/json"
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

func TestEnsureConfigFile(t *testing.T) {
	t.Run("creates config file when it does not exist", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "config", "config.json")

		err := EnsureConfigFile(configPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("failed to read created config: %v", err)
		}

		if string(data) != defaultConfig {
			t.Errorf("config content = %q, expected %q", string(data), defaultConfig)
		}
	})

	t.Run("does nothing when config already exists", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "config.json")

		original := `{"dataPath":"moods.json"}`

		err := os.WriteFile(configPath, []byte(original), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = EnsureConfigFile(configPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatal(err)
		}

		if string(data) != original {
			t.Errorf("file was modified: got %q, expected %q", string(data), original)
		}
	})
}

func TestLoadConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")

	content := `{"dataPath":"moods.json"}`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	config, err := LoadConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.DataPath != "moods.json" {
		t.Errorf("expected DataPath to be: %v, got: %v", "moods.json", config.DataPath)
	}
}

func TestLoadConfigFromInvalidFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")

	content := `{"invalid"___json1`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfigFromFile(configPath)
	if err == nil {
		t.Fatalf("expected error for invalid json")
	}
}

func TestSaveConfigToFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")

	content := Config{DataPath: "moods.json"}

	err := SaveConfigToFile(configPath, content)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config from %s: %v", configPath, err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		t.Errorf("failed to unmarshal config from %s: %v", configPath, err)
	}

	if config != content {
		t.Errorf("got %v, expected %v", config, content)
	}
}

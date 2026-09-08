package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildPath(t *testing.T) {
	dir, err := os.UserConfigDir()

	if err != nil {
		t.Fatalf("Something went wrong with getting os.UserConfigDir(): %v", err)
	}

	tests := []struct {
		name     string
		flags    flags
		expected string
	}{
		{
			name: "custom file",
			flags: flags{
				file: "config.json",
			},
			expected: filepath.Join(dir, "mindgo", "config.json"),
		},
		{
			name: "custom path",
			flags: flags{
				path: "mypath/test.json",
			},
			expected: "mypath/test.json",
		},
		{
			name:     "no flags",
			flags:    flags{},
			expected: filepath.Join(dir, "mindgo", "moods.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := BuildPath(tt.flags)

			if err != nil {
				t.Fatalf("Something went wrong: %v: ", err)
			}

			if path != tt.expected {
				t.Errorf("Something went wrong, expected: %v, got: %v", tt.expected, path)
			}
		})
	}
}

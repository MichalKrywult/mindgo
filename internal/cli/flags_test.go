package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MichalKrywult/mindgo/internal/config"
)

func TestParseFlags(t *testing.T) {

	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}() // anonymous function restoring original args

	tests := []struct {
		name     string
		args     []string
		expected flags
		err      bool
	}{
		{
			name:     "file flag",
			args:     []string{"program", "-file", "config.json"},
			expected: flags{file: "config.json"},
			err:      false,
		},
		{
			name:     "path flag",
			args:     []string{"program", "-path", "/tmp/config.json"},
			expected: flags{path: "/tmp/config.json"},
			err:      false,
		},
		{
			name:     "both flags",
			args:     []string{"program", "-file", "config.json", "-path", "/tmp/config.json"},
			expected: flags{},
			err:      true,
		},
		{
			name:     "no flags",
			args:     []string{"program"},
			expected: flags{file: ""},
			err:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			parsedFlags, err := ParseFlags()
			if tt.err {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
			if parsedFlags != tt.expected {
				t.Fatalf("expected: %v, got: %v", tt.expected, parsedFlags)
			}
		})
	}
}

func TestSelectDataPath(t *testing.T) {
	dir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("os.UserConfigDir() error = %v", err)
	}

	dataPath := filepath.Join(dir, "mindgo", "moods.json")

	tests := []struct {
		name     string
		flags    flags
		expected string
	}{
		{
			name: "file flag",
			flags: flags{
				file: "moods.json",
			},
			expected: dataPath,
		},
		{
			name:     "path flag",
			flags:    flags{path: "mypath/test.json"},
			expected: "mypath/test.json",
		},
		{
			name:     "no flags",
			flags:    flags{},
			expected: "my/test/moods.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := SelectDataPath(tt.flags, config.Config{DataPath: "my/test/moods.json"})
			if err != nil {
				t.Fatalf("SelectDataPath() error = %v", err)
			}

			if path != tt.expected {
				t.Errorf("SelectDataPath() = %q, want %q", path, tt.expected)
			}
		})
	}
}

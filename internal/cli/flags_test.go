package cli

import (
	"os"
	"path/filepath"
	"testing"
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
			expected: flags{file: defaultDataFilePath},
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
				t.Fatalf("Something went wrong: %v ", err)
			}

			if path != tt.expected {
				t.Errorf("Something went wrong, expected: %v, got: %v", tt.expected, path)
			}
		})
	}
}

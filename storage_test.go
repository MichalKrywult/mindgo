package main

import (
	"path/filepath"
	"testing"
)

func TestFileStorage_SaveAndLoad(t *testing.T) {
	// os.MkdirTemp creates temporary dictionary
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_moods.json")

	storage := FileStorage{filename: filePath}

	entries := []MoodEntry{{ID: 1, Mood: 5, Note: "Test"}}
	err := storage.Save(entries)
	if err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	loaded, err := storage.Load()
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}

	if len(loaded) != 1 || loaded[0].Mood != 5 {
		t.Errorf("loaded data mismatch")
	}
}

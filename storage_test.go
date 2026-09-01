package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileStorage_SaveAndLoad(t *testing.T) {
	// os.MkdirTemp creates temporary dictionary
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_moods.json")

	storage := &FileStorage{filename: filePath}

	entries := []MoodEntry{{Mood: 5, Note: "Test"}}
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

func TestFileStorage_LoadNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "non_existent.json")

	storage := &FileStorage{filename: filePath}

	entries, err := storage.Load()
	if err != nil {
		t.Fatalf("expected no error for non-existent file, got: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestFileStorage_LoadCorruptedFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "invalid.json")

	err := os.WriteFile(filePath, []byte("not json"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	storage := &FileStorage{filename: filePath}

	_, err = storage.Load()

	if err == nil {
		t.Fatal("expected error for invalid json format")
	}

}

func TestFileStorage_LoadEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "empty.json")

	err := os.WriteFile(filePath, []byte(""), 0644)
	if err != nil {
		t.Fatal(err)
	}

	storage := &FileStorage{filename: filePath}

	entries, err := storage.Load()

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if len(entries) != 0 {
		t.Fatalf("expected 0, got %v", len(entries))
	}

}

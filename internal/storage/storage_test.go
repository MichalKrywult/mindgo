package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

func TestFileStorage_SaveAndLoad(t *testing.T) {
	// os.MkdirTemp creates temporary dictionary
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_moods.json")

	storage := &FileStorage{Filename: filePath}

	entries := []domain.MoodEntry{{Mood: 5, Note: "Test"}}
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

	storage := &FileStorage{Filename: filePath}

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

	storage := &FileStorage{Filename: filePath}

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

	storage := &FileStorage{Filename: filePath}

	entries, err := storage.Load()

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if len(entries) != 0 {
		t.Fatalf("expected 0, got %v", len(entries))
	}

}

func TestMockStorage_SaveCopiesEntries(t *testing.T) {
	storage := &MockStorage{}
	entries := []domain.MoodEntry{{Mood: 5, Note: "Old"}}

	err := storage.Save(entries)
	if err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	entries[0].Mood = 4
	entries[0].Note = "New"

	loaded, err := storage.Load()
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}

	if loaded[0].Mood != 5 {
		t.Errorf("storage was modified through original slice")
	}

	if loaded[0].Note != "Old" {
		t.Errorf("storage was modified through original slice")
	}
}
func TestMockStorage_LoadCopiesEntries(t *testing.T) {
	storage := &MockStorage{}
	entries := []domain.MoodEntry{{Mood: 5, Note: "Old"}}

	err := storage.Save(entries)
	if err != nil {
		t.Fatalf("unexpected save error: %v", err)
	}

	loaded, err := storage.Load()
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}

	loaded[0].Mood = 4
	loaded[0].Note = "New"

	loadedAgain, err := storage.Load()
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}

	if loadedAgain[0].Mood != 5 {
		t.Errorf("expected mood: %v, got %d", entries[0].Mood, loadedAgain[0].Mood)
	}

	if loadedAgain[0].Note != "Old" {
		t.Errorf("expected note %q, got %q", entries[0].Note, loadedAgain[0].Note)
	}
}

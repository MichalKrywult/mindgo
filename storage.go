package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// MockStorage pretends to be real storage for tests
// it remembers last state of storage
type MockStorage struct {
	entries []MoodEntry
}

func (m *MockStorage) Save(entries []MoodEntry) error {
	m.entries = entries
	return nil
}

func (m *MockStorage) Load() ([]MoodEntry, error) {
	entries := make([]MoodEntry, len(m.entries))
	copy(entries, m.entries)

	return entries, nil
}

type Storage interface {
	Save([]MoodEntry) error
	Load() ([]MoodEntry, error)
}

type FileStorage struct {
	filename string
}

func (fs *FileStorage) Save(entries []MoodEntry) error {
	// data, err := json.Marshal(entries) would write everything in one line
	// while technically correct, it's unreadable
	data, err := json.MarshalIndent(entries, "", "  ") // writes data readable to humans
	if err != nil {
		return fmt.Errorf("failed to marshal entries: %w", err)
	}

	err = os.WriteFile(fs.filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", fs.filename, err)
	}

	return nil
}

func (fs *FileStorage) Load() ([]MoodEntry, error) {
	data, err := os.ReadFile(fs.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []MoodEntry{}, nil
		}
		return nil, fmt.Errorf("failed to read file %s: %w", fs.filename, err)
	}

	if len(data) == 0 {
		return []MoodEntry{}, nil
	}

	var entries []MoodEntry

	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal entries from %s: %w", fs.filename, err)
	}

	return entries, nil
}

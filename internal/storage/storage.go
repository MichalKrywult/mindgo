package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

// MockStorage pretends to be real storage for tests
// it remembers last state of storage
type MockStorage struct {
	entries []domain.MoodEntry
}

func (m *MockStorage) Save(entries []domain.MoodEntry) error {
	m.entries = entries
	return nil
}

func (m *MockStorage) Load() ([]domain.MoodEntry, error) {
	entries := make([]domain.MoodEntry, len(m.entries))
	copy(entries, m.entries)

	return entries, nil
}

type Storage interface {
	Save([]domain.MoodEntry) error
	Load() ([]domain.MoodEntry, error)
}

type FileStorage struct {
	Filename string
}

func (fs *FileStorage) Save(entries []domain.MoodEntry) error {
	// data, err := json.Marshal(entries) would write everything in one line
	// while technically correct, it's unreadable
	data, err := json.MarshalIndent(entries, "", "  ") // writes data readable to humans
	if err != nil {
		return fmt.Errorf("failed to marshal entries: %w", err)
	}

	err = os.WriteFile(fs.Filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", fs.Filename, err)
	}

	return nil
}

func (fs *FileStorage) Load() ([]domain.MoodEntry, error) {
	data, err := os.ReadFile(fs.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.MoodEntry{}, nil
		}
		return nil, fmt.Errorf("failed to read file %s: %w", fs.Filename, err)
	}

	if len(data) == 0 {
		return []domain.MoodEntry{}, nil
	}

	var entries []domain.MoodEntry

	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal entries from %s: %w", fs.Filename, err)
	}

	return entries, nil
}

package main

import (
	"encoding/json"
	"os"
)

type Storage interface {
	Save([]MoodEntry) error
	Load() ([]MoodEntry, error)
}

type FileStorage struct {
	filename string
}

func (fs FileStorage) Save(entries []MoodEntry) error {
	// data, err := json.Marshal(entries) would write everything in one line
	// while technically correct, it's unreadable
	data, err := json.MarshalIndent(entries, "", "  ") // writes data readable to humans
	if err != nil {
		return err
	}

	err = os.WriteFile(fs.filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (fs FileStorage) Load() ([]MoodEntry, error) {
	data, err := os.ReadFile(fs.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []MoodEntry{}, nil
		}
		return nil, err
	}

	var entries []MoodEntry

	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

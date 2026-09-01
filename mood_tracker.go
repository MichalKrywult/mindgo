package main

import (
	"fmt"
)

type MoodTracker struct {
	entries []MoodEntry
	entryID int
	storage Storage
}

func NewMoodTracker(storage Storage) (*MoodTracker, error) {
	// read data from the file
	entries, err := storage.Load()
	if err != nil {
		return nil, err
	}

	// search for max id in data that is being loaded
	maxID := 0
	for _, e := range entries {
		if e.ID > maxID {
			maxID = e.ID
		}
	}

	// return pointer to the created struct
	return &MoodTracker{
		entries: entries,
		entryID: maxID,
		storage: storage,
	}, nil
}

func (tracker *MoodTracker) AddEntry(entry MoodEntry) error {
	tracker.entryID++
	entry.ID = tracker.entryID
	tracker.entries = append(tracker.entries, entry)
	err := tracker.storage.Save(tracker.entries)
	if err != nil {
		return fmt.Errorf("unexpected error occurred when saving to file: %w", err)
	}
	return nil
}

func (tracker *MoodTracker) GetEntries() []MoodEntry {
	entries := make([]MoodEntry, len(tracker.entries))
	//we have to create space for that copy with make()
	// copy() doesn't make a slice bigger
	copy(entries, tracker.entries)
	//now we can copy tracker.entries into entries

	return entries
}

/* right now this function is useless
func (tracker *MoodTracker) findEntryByID(id int) (MoodEntry, error) {
	for _, entry := range tracker.entries {
		if entry.ID == id {
			return entry, nil
		}
	}
	return MoodEntry{}, fmt.Errorf("ID %d doesn't exist", id)
}*/

func (tracker *MoodTracker) findIndexByID(id int) (int, error) {
	//again, without pointer receiver (*) because this method doesn't modify the tracker
	for index, entry := range tracker.entries {
		if entry.ID == id {
			return index, nil
		}
	}
	return -1, fmt.Errorf("ID %d doesn't exist", id)
}

func (tracker *MoodTracker) RemoveEntryByID(id int) error {
	index, err := tracker.findIndexByID(id)
	if err != nil {
		return err
	}

	tracker.entries = append(tracker.entries[:index], tracker.entries[index+1:]...)
	// append(slice, element)
	// ... separates slice into separate arguments
	return nil
}

func (tracker *MoodTracker) EditEntryByID(id int, entry MoodEntry) error {
	index, err := tracker.findIndexByID(id)
	if err != nil {
		return err
	}

	entry.ID = id
	tracker.entries[index] = entry

	return nil
}

package main

import "fmt"

type MoodTracker struct {
	entries []MoodEntry
	entryID int
}

func (tracker *MoodTracker) AddEntry(entry MoodEntry) {
	// pointer receiver (*) because this method modifies the tracker
	tracker.entryID++
	entry.ID = tracker.entryID
	tracker.entries = append(tracker.entries, entry)
}

func (tracker MoodTracker) GetEntries() []MoodEntry {
	//without pointer receiver (*) because this method doesn't modify the tracker
	entries := make([]MoodEntry, len(tracker.entries))
	//we have to create space for that copy with make()
	// copy() doesn't make a slice bigger
	copy(entries, tracker.entries)
	//now we can copy tracker.entries into entries

	return entries
}

func (tracker MoodTracker) findEntryByID(id int) (MoodEntry, error) {
	//without pointer receiver (*) because this method doesn't modify the tracker
	for _, entry := range tracker.entries {
		if entry.ID == id {
			return entry, nil
		}
	}
	return MoodEntry{}, fmt.Errorf("ID %d doesn't exist", id)
}

func (tracker MoodTracker) findIndexByID(id int) (int, error) {
	//again, without pointer receiver (*) because this method doesn't modify the tracker
	for index, entry := range tracker.entries {
		if entry.ID == id {
			return index, nil
		}
	}
	return -1, fmt.Errorf("ID %d doesn't exist", id)
}

func (tracker *MoodTracker) RemoveEntry(id int) error {
	// pointer receiver (*) because this method modifies the tracker
	index, err := tracker.findIndexByID(id)
	if err != nil {
		return err
	}

	tracker.entries = append(tracker.entries[:index], tracker.entries[index+1:]...)
	// append(slice, element)
	// ... separates slice into separate arguments
	return nil
}

package main

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
	// copy() doesn't make a slice bigger,
	// we have to create space for that copy with make()
	copy(entries, tracker.entries)
	//now we can copy tracker.entries into entries
	return entries
}

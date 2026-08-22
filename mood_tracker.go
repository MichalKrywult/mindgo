package main

type MoodTracker struct {
	entries []MoodEntry
}

func (tracker *MoodTracker) AddEntry(entry MoodEntry) {
	tracker.entries = append(tracker.entries, entry)
}

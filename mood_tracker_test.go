package main

import (
	"testing"
)

func TestAddEntry(t *testing.T) {
	tracker := MoodTracker{}

	entry := MoodEntry{Mood: 2, Date: "2022/12/13", Note: "Test"}
	tracker.AddEntry(entry)

	if len(tracker.entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(tracker.entries))
	}

	if tracker.entries[0].Mood != 2 {
		t.Errorf("Expected mood to be 2, got %d", tracker.entries[0].Mood)
	}

	if tracker.entries[0].Date != "2022/12/13" {
		t.Errorf("Expected date to be %q, got %s", "2022/12/13", tracker.entries[0].Date)
	}

	if tracker.entries[0].Note != "Test" {
		t.Errorf("Expected note to be %q, got %s", "Test", tracker.entries[0].Note)
	}
}

func TestGetEntries(t *testing.T) {
	tracker := MoodTracker{}

	entry := MoodEntry{Mood: 2, Date: "2022/12/13", Note: "Test"}

	tracker.AddEntry(entry)
	tracker.AddEntry(entry)

	data := tracker.GetEntries()
	if len(data) != len(tracker.entries) {
		t.Errorf("Expected %d entries, got %d", len(tracker.entries), len(data))
	}

	data[0].Mood = 99
	if tracker.entries[0].Mood == 99 {
		t.Error("Modifying returned entries changed tracker state")
	}
}

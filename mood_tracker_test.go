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

	if data[0].ID != 1 || data[1].ID != 2 {
		t.Error("Entries ID are invalid")
	}
}

func TestFindIndexByID(t *testing.T) {
	tracker := MoodTracker{}

	entry := MoodEntry{Mood: 2, Date: "2022/12/13", Note: "Test"}

	tracker.AddEntry(entry)
	tracker.AddEntry(entry)

	index, err := tracker.findIndexByID(1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if index != 0 {
		t.Errorf("Expected %d, got %d", 0, index)
	}

	index, err = tracker.findIndexByID(2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if index != 1 {
		t.Errorf("Expected %d, got %d", 1, index)
	}

	index, err = tracker.findIndexByID(3)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	if index != -1 {
		t.Errorf("Expected %d, got %d", -1, index)
	}
}

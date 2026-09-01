package main

import (
	"testing"
	"time"
)

func TestAddEntry(t *testing.T) {
	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}
	time := time.Now()

	entry := MoodEntry{Mood: 2, Date: time, Note: "Test"}
	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	if len(tracker.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(tracker.entries))
	}

	if tracker.entries[0].Mood != 2 {
		t.Errorf("expected mood to be 2, got %d", tracker.entries[0].Mood)
	}

	if tracker.entries[0].Date != time {
		t.Errorf("expected date to be %q, got %s", time, tracker.entries[0].Date)
	}

	if tracker.entries[0].Note != "Test" {
		t.Errorf("expected note to be %q, got %s", "Test", tracker.entries[0].Note)
	}
}

func TestGetEntries(t *testing.T) {
	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}
	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	data := tracker.GetEntries()
	if len(data) != len(tracker.entries) {
		t.Errorf("expected %d entries, got %d", len(tracker.entries), len(data))
	}

	data[0].Mood = 99
	if tracker.entries[0].Mood == 99 {
		t.Error("modifying returned entries changed tracker state")
	}

	if data[0].ID != 1 || data[1].ID != 2 {
		t.Error("entries ID are invalid")
	}
}

func TestFindIndexByID(t *testing.T) {
	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}
	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	index, err := tracker.findIndexByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if index != 0 {
		t.Errorf("expected %d, got %d", 0, index)
	}

	index, err = tracker.findIndexByID(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if index != 1 {
		t.Errorf("expected %d, got %d", 1, index)
	}

	index, err = tracker.findIndexByID(3)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if index != -1 {
		t.Errorf("expected %d, got %d", -1, index)
	}
}

func TestRemoveEntry(t *testing.T) {
	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}
	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}
	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}
	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	err = tracker.RemoveEntryByID(2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	data := tracker.GetEntries()
	if len(data) != 2 {
		t.Error("entry wasn't removed from the tracker")
	}

	if data[0].ID != 1 || data[1].ID != 3 {
		t.Error("wrong entry was removed from the tracker")
	}
}

func TestRemoveEntryWithInvalidID(t *testing.T) {
	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}
	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	err = tracker.RemoveEntryByID(4)
	if err == nil {
		t.Error("expected an error when removing an entry with an invalid ID")
	}

	if len(tracker.entries) != 1 {
		t.Error("tracker was modified after attempting to remove an invalid ID")
	}
}

func TestEditEntryByID(t *testing.T) {
	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	err = tracker.AddEntry(MoodEntry{Mood: 2, Date: time.Now(), Note: "Old"})
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	err = tracker.EditEntryByID(1, MoodEntry{Mood: 5, Date: time.Now(), Note: "New"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data := tracker.GetEntries()

	if len(data) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(data))
	}

	if data[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", data[0].ID)
	}

	if data[0].Mood != 5 {
		t.Errorf("expected mood 5, got %d", data[0].Mood)
	}

	if data[0].Note != "New" {
		t.Errorf("expected note %q, got %q", "New", data[0].Note)
	}
}

func TestEditEntryByIDWithInvalidID(t *testing.T) {
	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}
	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	err = tracker.EditEntryByID(2, entry)
	if err == nil || len(tracker.GetEntries()) != 1 {
		t.Error("expected an error for invalid ID")
	}

	data := tracker.GetEntries()
	if data[0].ID != 1 {
		t.Error("transaction was edited despite wrong ID")
	}
}

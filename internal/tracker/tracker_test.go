package tracker

import (
	"testing"
	"time"

	"github.com/MichalKrywult/mindgo/internal/domain"
	"github.com/MichalKrywult/mindgo/internal/storage"
)

func TestAddEntry(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}
	time := time.Now()

	entry := domain.MoodEntry{Mood: 2, Date: time, Note: "Test"}
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

func TestInvalidAddEntry(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.ErrorStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}
	time := time.Now()

	entry := domain.MoodEntry{Mood: 2, Date: time, Note: "Test"}
	err = tracker.AddEntry(entry)
	if err == nil {
		t.Fatalf("expected an error when adding an entry")
	}

	if len(tracker.entries) != 0 {
		t.Errorf("expected tracker to be empty, got: %v instead", len(tracker.entries))
	}

	if tracker.entryID != 0 {
		t.Errorf("expected entryID to be 0, got: %v instead", tracker.entryID)
	}
}

func TestGetEntries(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := domain.MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}

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

func TestEditEntryByIndex(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := domain.MoodEntry{Mood: 2, Date: time.Now(), Note: "Old"}
	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("unexpected error occured when adding entry, %v", err)
	}

	newEntry := domain.MoodEntry{Mood: 8, Date: time.Now(), Note: "New"}

	err = tracker.EditEntryByIndex(0, newEntry)
	if err != nil {
		t.Fatalf("unexpected error occured when editing entry, %v", err)
	}

	if tracker.entries[0].Note != "New" {
		t.Errorf("expected 'Note' to be 'New', got 'Note' = %s", tracker.entries[0].Note)
	}
}

func TestEditEntryByIndexForInvalidIndex(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}
	entry := domain.MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}

	err = tracker.EditEntryByIndex(0, entry)
	if err == nil {
		t.Fatal("expected error for invalid index")
	}
}

func TestRemoveEntryByIndex(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := domain.MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}
	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("unexpected error occured when adding entry, %v", err)
	}

	err = tracker.RemoveEntryByIndex(0)
	if err != nil {
		t.Fatalf("unexpected error occured when removing entry, %v", err)
	}

	if len(tracker.entries) != 0 {
		t.Fatal("entry wasn't removed from the tracker")
	}
}

func TestRemoveEntryByIndexForInvalidIndex(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := domain.MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}
	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("unexpected error occured when adding entry, %v", err)
	}

	err = tracker.RemoveEntryByIndex(2)
	if err == nil {
		t.Fatalf("expected error for invalid index")
	}

	if len(tracker.entries) != 1 {
		t.Fatal("entry was removed from the tracker, despite invalid index")
	}
}

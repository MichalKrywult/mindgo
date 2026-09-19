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

	now := time.Now()
	entry := domain.MoodEntry{Mood: 2, Date: now, Note: "Test"}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("unexpected error occurred when adding entry, %v", err)
	}

	if len(tracker.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(tracker.entries))
	}

	if tracker.entries[0].Mood != 2 {
		t.Errorf("expected mood to be 2, got %d", tracker.entries[0].Mood)
	}

	if tracker.entries[0].Date != now {
		t.Errorf("expected date to be %q, got %s", now, tracker.entries[0].Date)
	}

	if tracker.entries[0].Note != "Test" {
		t.Errorf("expected note to be %q, got %s", "Test", tracker.entries[0].Note)
	}
}

func TestAddEntryStorageError(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.ErrorStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	now := time.Now()
	entry := domain.MoodEntry{Mood: 2, Date: now, Note: "Test"}

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
		t.Fatalf("unexpected error occurred when adding entry, %v", err)
	}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("unexpected error occurred when adding entry, %v", err)
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

	entry := domain.MoodEntry{
		Mood: 2,
		Date: time.Now(),
		Note: "Old",
	}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("unexpected error occurred when adding entry, %v", err)
	}

	newEntry := domain.MoodEntry{
		Mood: 8,
		Date: time.Now(),
		Note: "New",
	}

	err = tracker.EditEntryByIndex(0, newEntry)
	if err != nil {
		t.Fatalf("unexpected error occurred when editing entry, %v", err)
	}

	got := tracker.entries[0]

	if got.Mood != newEntry.Mood {
		t.Errorf("expected mood to be %d, got %d", newEntry.Mood, got.Mood)
	}

	if got.Date != newEntry.Date {
		t.Errorf("expected date to be %q, got %s", newEntry.Date, got.Date)
	}

	if got.Note != newEntry.Note {
		t.Errorf("expected note to be %q, got %q", newEntry.Note, got.Note)
	}

	// ID should not change when editing an entry.
	if got.ID != 1 {
		t.Errorf("expected ID to remain 1, got %d", got.ID)
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

func TestEditEntryByIndexForNegativeIndex(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := domain.MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}

	err = tracker.EditEntryByIndex(-1, entry)
	if err == nil {
		t.Fatal("expected error for negative index")
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
		t.Fatalf("unexpected error occurred when adding entry, %v", err)
	}

	err = tracker.RemoveEntryByIndex(0)
	if err != nil {
		t.Fatalf("unexpected error occurred when removing entry, %v", err)
	}

	if len(tracker.entries) != 0 {
		t.Fatal("entry wasn't removed from the tracker")
	}
}

func TestRemoveEntryByIndexFromMiddle(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entries := []domain.MoodEntry{
		{Mood: 1, Date: time.Now(), Note: "First"},
		{Mood: 5, Date: time.Now(), Note: "Second"},
		{Mood: 10, Date: time.Now(), Note: "Third"},
	}

	for _, entry := range entries {
		err = tracker.AddEntry(entry)
		if err != nil {
			t.Fatalf("unexpected error occurred when adding entry, %v", err)
		}
	}

	err = tracker.RemoveEntryByIndex(1)
	if err != nil {
		t.Fatalf("unexpected error occurred when removing entry, %v", err)
	}

	data := tracker.GetEntries()

	if len(data) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(data))
	}

	if data[0].Note != "First" {
		t.Errorf("expected first entry to be 'First', got %q", data[0].Note)
	}

	if data[1].Note != "Third" {
		t.Errorf("expected second entry to be 'Third', got %q", data[1].Note)
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
		t.Fatalf("unexpected error occurred when adding entry, %v", err)
	}

	err = tracker.RemoveEntryByIndex(2)
	if err == nil {
		t.Fatalf("expected error for invalid index")
	}

	if len(tracker.entries) != 1 {
		t.Fatal("entry was removed from the tracker, despite invalid index")
	}
}

func TestRemoveEntryByIndexForNegativeIndex(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := domain.MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("unexpected error occurred when adding entry, %v", err)
	}

	err = tracker.RemoveEntryByIndex(-1)
	if err == nil {
		t.Fatalf("expected error for negative index")
	}

	if len(tracker.entries) != 1 {
		t.Fatal("entry was removed from the tracker, despite invalid index")
	}
}

func TestCreatingTrackerBackup(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := domain.MoodEntry{Mood: 8, Note: "Test"}
	entry2 := domain.MoodEntry{Mood: 2, Note: "Test2"}

	tracker.AddEntry(entry)

	backup := tracker.createTrackerBackup()
	tracker.EditEntryByIndex(1, entry2)
	tracker.AddEntry(entry2)

	if backup.entries[0].Mood != 8 {
		t.Errorf("expected mood %d, got %d", 8, backup.entries[0].Mood)
	}

	if backup.entries[0].Note != "Test" {
		t.Errorf("expected note %q, got %q", "Test", backup.entries[0].Note)
	}

	if backup.entryID != 1 {
		t.Errorf("expected entryID %d, got %d", 1, backup.entryID)
	}
}

func TestRestoringTrackerBackup(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := domain.MoodEntry{Mood: 8, Note: "Test"}
	entry2 := domain.MoodEntry{Mood: 2, Note: "Test2"}

	tracker.AddEntry(entry)

	backup := tracker.createTrackerBackup()

	tracker.EditEntryByIndex(1, entry2)
	tracker.AddEntry(entry2)

	tracker.restoreTrackerFromBackup(backup)

	if len(tracker.entries) != 1 {
		t.Errorf("expected entries length %d, got %d", 1, len(tracker.entries))
	}

	if tracker.entries[0].Mood != 8 {
		t.Errorf("expected mood %d, got %d", 8, tracker.entries[0].Mood)
	}

	if tracker.entries[0].Note != "Test" {
		t.Errorf("expected note %q, got %q", "Test", tracker.entries[0].Note)
	}

	if tracker.entryID != 1 {
		t.Errorf("expected entryID %d, got %d", 1, tracker.entryID)
	}
}

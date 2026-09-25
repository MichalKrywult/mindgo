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
		t.Fatalf("unexpected error when adding entry: %v", err)
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
	tracker, err := NewMoodTracker(&storage.SaveErrorStorage{})
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
		t.Errorf("expected tracker to be empty, got %d entries instead", len(tracker.entries))
	}

	if tracker.entryID != 0 {
		t.Errorf("expected entryID to be 0, got %d instead", tracker.entryID)
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
		t.Fatalf("unexpected error when adding entry: %v", err)
	}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("unexpected error when adding entry: %v", err)
	}

	data := tracker.GetEntries()

	if len(data) != len(tracker.entries) {
		t.Errorf("expected %d entries, got %d", len(tracker.entries), len(data))
	}

	data[0].Mood = 99
	if tracker.entries[0].Mood == 99 {
		t.Error("modifying returned entries changed the tracker state")
	}

	if data[0].ID != 1 || data[1].ID != 2 {
		t.Error("entry IDs are invalid")
	}
}

func TestEditEntryByIndex(t *testing.T) {
	tests := []struct {
		name           string
		initialEntries []domain.MoodEntry
		index          int
		entry          domain.MoodEntry
		expectedEntry  *domain.MoodEntry // pointers can be nil, values would be initialized with default
		expectedError  bool
	}{
		{
			name:           "negative index",
			initialEntries: []domain.MoodEntry{},
			index:          -1,
			entry:          domain.MoodEntry{Mood: 2, Note: "Test"},
			expectedError:  true,
		},
		{
			name:           "empty tracker",
			initialEntries: []domain.MoodEntry{},
			index:          0,
			entry:          domain.MoodEntry{Mood: 2, Note: "Test"},
			expectedError:  true,
		},
		{
			name:           "valid index",
			initialEntries: []domain.MoodEntry{{ID: 1, Mood: 3, Note: "Original"}},
			index:          0,
			entry:          domain.MoodEntry{Mood: 5, Note: "Updated"},
			expectedEntry:  &domain.MoodEntry{ID: 1, Mood: 5, Note: "Updated"},
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &storage.MockStorage{}

			err := mockStorage.Save(tt.initialEntries)
			if err != nil {
				t.Fatalf("failed to save initial entries: %v", err)
			}

			tracker, err := NewMoodTracker(mockStorage)
			if err != nil {
				t.Fatalf("failed to create tracker: %v", err)
			}

			err = tracker.EditEntryByIndex(tt.index, tt.entry)

			if tt.expectedError && err == nil {
				t.Error("expected an error, got nil")
			}

			if !tt.expectedError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tt.expectedEntry != nil {
				got := tracker.entries[tt.index]

				if got.ID != tt.expectedEntry.ID {
					t.Errorf("expected ID to be %d, got %d", tt.expectedEntry.ID, got.ID)
				}

				if got.Mood != tt.expectedEntry.Mood {
					t.Errorf("expected mood to be %d, got %d", tt.expectedEntry.Mood, got.Mood)
				}

				if got.Note != tt.expectedEntry.Note {
					t.Errorf("expected note to be %q, got %q", tt.expectedEntry.Note, got.Note)
				}
			}
		})
	}
}

func TestEditEntryByIndexStorageError(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.SaveErrorStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := domain.MoodEntry{Mood: 2, Note: "Original"}

	tracker.entries = []domain.MoodEntry{entry}
	tracker.entryID = 0

	newEntry := domain.MoodEntry{Mood: 5, Note: "Edited"}

	err = tracker.EditEntryByIndex(0, newEntry)
	if err == nil {
		t.Fatalf("expected an error when editing an entry")
	}

	if len(tracker.entries) != 1 {
		t.Errorf("expected tracker to contain 1 entry, got %d", len(tracker.entries))
	}

	if tracker.entries[0] != entry {
		t.Errorf("expected original entry to be restored, got %v instead", tracker.entries[0])
	}

	if tracker.entryID != entry.ID {
		t.Errorf("expected entryID to be %d, got %d instead", entry.ID, tracker.entryID)
	}
}

func TestRemoveEntryByIndex(t *testing.T) {

	tests := []struct {
		name            string
		index           int
		initialEntries  []domain.MoodEntry
		expectedEntries []domain.MoodEntry
		expectedError   bool
	}{
		{
			name:            "empty tracker",
			index:           0,
			initialEntries:  []domain.MoodEntry{},
			expectedEntries: []domain.MoodEntry{},
			expectedError:   true,
		},
		{
			name:            "valid index",
			index:           0,
			initialEntries:  []domain.MoodEntry{{Mood: 2, Note: "Test"}},
			expectedEntries: []domain.MoodEntry{},
			expectedError:   false,
		},
		{
			name:            "index out of range",
			index:           1,
			initialEntries:  []domain.MoodEntry{{Mood: 2, Note: "Test"}},
			expectedEntries: []domain.MoodEntry{{Mood: 2, Note: "Test"}},
			expectedError:   true,
		},
		{
			name:            "negative index",
			index:           -1,
			initialEntries:  []domain.MoodEntry{{Mood: 2, Note: "Test"}},
			expectedEntries: []domain.MoodEntry{{Mood: 2, Note: "Test"}},
			expectedError:   true,
		},
		{
			name:            "valid index, multiple entries",
			index:           1,
			initialEntries:  []domain.MoodEntry{{Mood: 1, Note: "Test"}, {Mood: 2, Note: "Test"}, {Mood: 3, Note: "Test"}},
			expectedEntries: []domain.MoodEntry{{Mood: 1, Note: "Test"}, {Mood: 3, Note: "Test"}},
			expectedError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &storage.MockStorage{}

			err := mockStorage.Save(tt.initialEntries)
			if err != nil {
				t.Fatalf("failed to save initial entries: %v", err)
			}

			tracker, err := NewMoodTracker(mockStorage)
			if err != nil {
				t.Fatalf("failed to create tracker: %v", err)
			}

			err = tracker.RemoveEntryByIndex(tt.index)
			if tt.expectedError && err == nil {
				t.Error("expected an error, got nil")
			}

			if !tt.expectedError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(tracker.entries) != len(tt.expectedEntries) {
				t.Fatal("length of tracker and expected length is not equal")
			}

			for i, expectedEntry := range tt.expectedEntries {
				got := tracker.entries[i]

				if expectedEntry.Mood != got.Mood {
					t.Errorf("expected mood: %v, got: %v instead", expectedEntry.Mood, got.Mood)
				}

				if expectedEntry.Note != got.Note {
					t.Errorf("expected note: %v, got: %v instead", expectedEntry.Note, got.Note)
				}
			}
		})
	}
}

func TestCreatingTrackerBackup(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := domain.MoodEntry{Mood: 8, Note: "Test"}
	entry2 := domain.MoodEntry{Mood: 2, Note: "Test2"}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("unexpected error when adding entry: %v", err)
	}

	backup := tracker.createTrackerBackup()

	err = tracker.EditEntryByIndex(0, entry2)
	if err != nil {
		t.Fatalf("unexpected error when editing entry: %v", err)
	}

	err = tracker.AddEntry(entry2)
	if err != nil {
		t.Fatalf("unexpected error when adding entry: %v", err)
	}

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

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("unexpected error when adding entry: %v", err)
	}

	backup := tracker.createTrackerBackup()

	err = tracker.EditEntryByIndex(0, entry2)
	if err != nil {
		t.Fatalf("unexpected error when editing entry: %v", err)
	}

	err = tracker.AddEntry(entry2)
	if err != nil {
		t.Fatalf("unexpected error when adding entry: %v", err)
	}

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

func TestRemoveEntryByIndexStorageError(t *testing.T) {
	tracker, err := NewMoodTracker(&storage.SaveErrorStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry := domain.MoodEntry{Mood: 2, Note: "Test"}

	tracker.entries = []domain.MoodEntry{entry}
	tracker.entryID = 0

	err = tracker.RemoveEntryByIndex(0)
	if err == nil {
		t.Fatalf("expected an error when removing an entry")
	}

	if len(tracker.entries) != 1 {
		t.Errorf("expected tracker to contain 1 entry, got %d", len(tracker.entries))
	}

	if tracker.entries[0] != entry {
		t.Errorf("expected original entry to be restored, got %v instead", tracker.entries[0])
	}

	if tracker.entryID != entry.ID {
		t.Errorf("expected entryID to be %d, got %d instead", entry.ID, tracker.entryID)
	}
}

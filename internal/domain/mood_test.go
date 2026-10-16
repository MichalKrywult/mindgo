package domain

import (
	"testing"
	"time"
)

func TestNewMoodEntryDoesNotShareTags(t *testing.T) {
	tags := []string{"test1", "test2"}

	entry, err := NewMoodEntry(2, time.Now(), "Test", tags)
	if err != nil {
		t.Fatalf("error when creating entry: %v", err)
	}

	tags[0] = "test3"

	if entry.Tags[0] != "test1" {
		t.Errorf("entry tags changed: got %q, want %q", entry.Tags[0], "test1")
	}
}

func TestNewMoodEntryValidMood(t *testing.T) {
	entry, err := NewMoodEntry(5, time.Now(), "Test note", []string{"test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if entry.Mood != 5 {
		t.Errorf("got mood %d, want 5", entry.Mood)
	}

	if entry.Note != "Test note" {
		t.Errorf("got note %q, want %q", entry.Note, "Test note")
	}

	if len(entry.Tags) != 1 || entry.Tags[0] != "test" {
		t.Errorf("got tags %v, want [test]", entry.Tags)
	}
}

func TestNewMoodEntryInvalidMood(t *testing.T) {
	tests := []struct {
		name string
		mood int
	}{
		{name: "below minimum", mood: 0},
		{name: "above maximum", mood: 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMoodEntry(tt.mood, time.Now(), "Test", nil)
			if err == nil {
				t.Fatal("expected error for invalid mood")
			}
		})
	}
}

func TestNewMoodEntryDuplicateTags(t *testing.T) {
	tags := []string{"test", "test"}

	_, err := NewMoodEntry(5, time.Now(), "Test", tags)
	if err == nil {
		t.Fatal("expected error for duplicate tags")
	}
}
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

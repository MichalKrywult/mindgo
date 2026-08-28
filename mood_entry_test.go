package main

import (
	"testing"
	"time"
)

func TestNewMoodEntryWithCorrectData(t *testing.T) {
	time := time.Now()
	entry, err := NewMoodEntry(2, time, "Test")

	if err != nil {
		t.Fatalf("NewMoodEntry() returned an error: %v", err)
	}

	if entry.Mood != 2 {
		t.Errorf("Mood = %d, expected 2", entry.Mood)
	}

	if entry.Date != time {
		t.Errorf("Date = %q, expected %q", entry.Date, time)
	}

	if entry.Note != "Test" {
		t.Errorf("Note = %q, expected %q", entry.Note, "Test")
	}
}

func TestNewMoodEntryWithMoodValueAboveRange(t *testing.T) {
	_, err := NewMoodEntry(12, time.Now(), "Test")

	if err == nil {
		t.Error("expected an error for an invalid mood value")
	}
}

func TestNewMoodEntryWithMoodValueBelowRange(t *testing.T) {
	_, err := NewMoodEntry(0, time.Now(), "Test")

	if err == nil {
		t.Error("expected an error for an invalid mood value")
	}
}

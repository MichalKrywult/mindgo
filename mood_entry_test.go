package main

import (
	"testing"
)

func TestNewMoodEntryWithCorrectData(t *testing.T) {
	entry, err := NewMoodEntry(2, "2022/12/13", "Test")

	if err != nil {
		t.Fatalf("NewMoodEntry() returned an error: %v", err)
	}

	if entry.Mood != 2 {
		t.Errorf("Mood = %d, expected 2", entry.Mood)
	}

	if entry.Date != "2022/12/13" {
		t.Errorf("Date = %q, expected %q", entry.Date, "2022/12/13")
	}

	if entry.Note != "Test" {
		t.Errorf("Note = %q, expected %q", entry.Note, "Test")
	}
}

func TestNewMoodEntryWithIncorrectDate(t *testing.T) {
	_, err := NewMoodEntry(2, "", "Test")

	if err == nil {
		t.Error("Expected an error for an empty date")
	}
}

func TestNewMoodEntryWithMoodValueAboveRange(t *testing.T) {
	_, err := NewMoodEntry(12, "2022/12/13", "Test")

	if err == nil {
		t.Error("Expected an error for an invalid mood value")
	}
}

func TestNewMoodEntryWithMoodValueBelowRange(t *testing.T) {
	_, err := NewMoodEntry(0, "2022/12/13", "Test")

	if err == nil {
		t.Error("Expected an error for an invalid mood value")
	}
}

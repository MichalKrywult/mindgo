package main

import (
	"testing"
	"time"
)

func TestCalculateAverageMood(t *testing.T) {
	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry1 := MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}
	time.Sleep(100 * time.Millisecond)
	entry2 := MoodEntry{Mood: 4, Date: time.Now(), Note: "Test"}

	err = tracker.addEntry(entry1)
	if err != nil {
		t.Fatalf("failed to add entry: %v", err)
	}

	err = tracker.addEntry(entry2)
	if err != nil {
		t.Fatalf("failed to add entry: %v", err)
	}

	average, err := tracker.CalculateAverageMood()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if average != 3 {
		t.Errorf("calculated average is wrong, expected: %v, got: %v", 3, average)
	}
}

func TestCalculateAverageMoodWhenEmpty(t *testing.T) {
	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	_, err = tracker.CalculateAverageMood()

	if err == nil {
		t.Error("expected an error for trying to calculate average on empty tracker ")
	}
}

package main

import (
	"testing"
	"time"
)

func TestCalculateAverageMood(t *testing.T) {
	tracker := MoodTracker{}

	entry1 := MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}
	time.Sleep(100 * time.Millisecond)
	entry2 := MoodEntry{Mood: 4, Date: time.Now(), Note: "Test"}

	tracker.AddEntry(entry1)
	tracker.AddEntry(entry2)

	average, err := tracker.CalculateAverageMood()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if average != 3 {
		t.Errorf("Calculated average is wrong, expected: %v, got: %v", 3, average)
	}
}

func TestCalculateAverageMoodWhenEmpty(t *testing.T) {
	tracker := MoodTracker{}

	_, err := tracker.CalculateAverageMood()

	if err == nil {
		t.Error("Expected an error for trying to calculate average on empty tracker ")
	}
}

package main

import (
	"strings"
	"testing"
	"time"
)

func TestCLIAddsMoodEntry(t *testing.T) {
	input := strings.NewReader("1\n2\nGreat day\n0\n")

	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	cli := NewCLI(tracker, input)

	cli.show()

	entries := tracker.GetEntries()

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].Mood != 2 {
		t.Errorf("expected mood 2, got %d", entries[0].Mood)
	}

	if entries[0].Note != "Great day" {
		t.Errorf("expected note %q, got %q", "Great day", entries[0].Note)
	}
}

func TestCLIEditMoodEntry(t *testing.T) {
	input := strings.NewReader("2\n1\n8\nGood day\n0\n")
	tracker, err := NewMoodTracker(&MockStorage{})

	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry, err := NewMoodEntry(2, time.Now(), "Bad day")
	if err != nil {
		t.Fatal(err)
	}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	cli := NewCLI(tracker, input)
	cli.show()

	entries := tracker.GetEntries()

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].Mood != 8 {
		t.Errorf("expected mood 8, got %d", entries[0].Mood)
	}

	if entries[0].Note != "Good day" {
		t.Errorf("expected note %q, got %q", "Good day", entries[0].Note)
	}
}

func TestCLIRemoveMoodEntry(t *testing.T) {
	input := strings.NewReader("3\n1\n0\n")

	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry, err := NewMoodEntry(2, time.Now(), "Test")
	if err != nil {
		t.Fatal(err)
	}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	cli := NewCLI(tracker, input)
	cli.show()

	entries := tracker.GetEntries()

	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestCLIRemoveMoodEntryWithInvalidIndex(t *testing.T) {
	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry, err := NewMoodEntry(2, time.Now(), "Test")
	if err != nil {
		t.Fatal(err)
	}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	input := strings.NewReader("3\n99\n0\n")

	cli := NewCLI(tracker, input)
	cli.show()

	entries := tracker.GetEntries()

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestCLIRemoveMoodEntryWithInvalidInput(t *testing.T) {
	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry, err := NewMoodEntry(2, time.Now(), "Test")
	if err != nil {
		t.Fatal(err)
	}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	input := strings.NewReader("3\nabc\n0\n")

	cli := NewCLI(tracker, input)
	cli.show()

	entries := tracker.GetEntries()

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestCLIEditWithEmptyTracker(t *testing.T) {
	tracker, err := NewMoodTracker(&MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	input := strings.NewReader("2\n0\n")

	cli := NewCLI(tracker, input)
	cli.show()

	entries := tracker.GetEntries()

	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

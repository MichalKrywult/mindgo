package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/MichalKrywult/mindgo/internal/domain"
	"github.com/MichalKrywult/mindgo/internal/storage"
	"github.com/MichalKrywult/mindgo/internal/tracker"
)

func TestCLIAddsMoodEntry(t *testing.T) {
	input := strings.NewReader("1\n2\nGreat day\n0\n")

	tracker, err := tracker.NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	cli := NewCLI(tracker, input)

	cli.Show()

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
	tracker, err := tracker.NewMoodTracker(&storage.MockStorage{})

	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry, err := domain.NewMoodEntry(2, time.Now(), "Bad day")
	if err != nil {
		t.Fatal(err)
	}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	cli := NewCLI(tracker, input)
	cli.Show()

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

	tracker, err := tracker.NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry, err := domain.NewMoodEntry(2, time.Now(), "Test")
	if err != nil {
		t.Fatal(err)
	}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	cli := NewCLI(tracker, input)
	cli.Show()

	entries := tracker.GetEntries()

	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestCLIRemoveMoodEntryWithInvalidIndex(t *testing.T) {
	tracker, err := tracker.NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry, err := domain.NewMoodEntry(2, time.Now(), "Test")
	if err != nil {
		t.Fatal(err)
	}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	input := strings.NewReader("3\n99\n0\n")

	cli := NewCLI(tracker, input)
	cli.Show()

	entries := tracker.GetEntries()

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestCLIRemoveMoodEntryWithInvalidInput(t *testing.T) {
	tracker, err := tracker.NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry, err := domain.NewMoodEntry(2, time.Now(), "Test")
	if err != nil {
		t.Fatal(err)
	}

	err = tracker.AddEntry(entry)
	if err != nil {
		t.Fatalf("Unexpected error occured when adding entry, %v", err)
	}

	input := strings.NewReader("3\nabc\n0\n")

	cli := NewCLI(tracker, input)
	cli.Show()

	entries := tracker.GetEntries()

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestCLIEditWithEmptyTracker(t *testing.T) {
	tracker, err := tracker.NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	input := strings.NewReader("2\n0\n")

	cli := NewCLI(tracker, input)
	cli.Show()

	entries := tracker.GetEntries()

	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestIsIndexValid(t *testing.T) {
	tracker, err := tracker.NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	input := strings.NewReader("0\n")
	cli := NewCLI(tracker, input)

	if cli.IsIndexValid(0) {
		t.Error("expected IsIndexValid(0) to be false for empty tracker")
	}

	entry := domain.MoodEntry{Mood: 5, Date: time.Now(), Note: "Test"}
	_ = tracker.AddEntry(entry)

	if !cli.IsIndexValid(0) {
		t.Error("expected IsIndexValid(0) to be true")
	}
	if cli.IsIndexValid(1) {
		t.Error("expected IsIndexValid(1) to be false")
	}
	if cli.IsIndexValid(-1) {
		t.Error("expected IsIndexValid(-1) to be false")
	}
}

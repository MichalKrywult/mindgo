package stats

import (
	"testing"
	"time"

	"github.com/MichalKrywult/mindgo/internal/domain"
	"github.com/MichalKrywult/mindgo/internal/storage"
	"github.com/MichalKrywult/mindgo/internal/tracker"
)

func TestCalculateAverageMood(t *testing.T) {
	tracker, err := tracker.NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	entry1 := domain.MoodEntry{Mood: 2, Date: time.Now(), Note: "Test"}
	time.Sleep(100 * time.Millisecond)
	entry2 := domain.MoodEntry{Mood: 4, Date: time.Now(), Note: "Test"}

	err = tracker.AddEntry(entry1)
	if err != nil {
		t.Fatalf("failed to add entry: %v", err)
	}

	err = tracker.AddEntry(entry2)
	if err != nil {
		t.Fatalf("failed to add entry: %v", err)
	}

	average, err := CalculateAverageMood(tracker.GetEntries())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if average != 3 {
		t.Errorf("calculated average is wrong, expected: %v, got: %v", 3, average)
	}
}

func TestCalculateAverageMoodWhenEmpty(t *testing.T) {
	tracker, err := tracker.NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	_, err = CalculateAverageMood(tracker.GetEntries())

	if err == nil {
		t.Error("expected an error for trying to calculate average on empty tracker ")
	}
}

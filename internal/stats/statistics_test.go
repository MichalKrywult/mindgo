package stats

import (
	"testing"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

func TestCalculateAverageMood(t *testing.T) {
	entries := []domain.MoodEntry{
		{Mood: 2, Note: "Test"},
		{Mood: 4, Note: "Test"},
	}

	average, err := CalculateAverageMood(entries)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if average != 3 {
		t.Errorf("expected average 3, got %v", average)
	}
}

func TestCalculateAverageMoodWhenEmpty(t *testing.T) {
	entries := []domain.MoodEntry{}

	_, err := CalculateAverageMood(entries)

	if err == nil {
		t.Fatal("expected error for empty entries")
	}
}

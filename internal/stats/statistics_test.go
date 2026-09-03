package stats

import (
	"testing"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

func TestCalculateStatsSummary(t *testing.T) {
	entries := []domain.MoodEntry{
		{Mood: 2, Note: "Test"},
		{Mood: 4, Note: "Test"},
	}

	summary, err := CalculateStatsSummary(entries)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.TotalCount != 2 {
		t.Errorf("expected total count 2, got %v", summary.TotalCount)
	}

	if summary.MinMood != 2 {
		t.Errorf("expected min mood 2, got %v", summary.MinMood)
	}

	if summary.MaxMood != 4 {
		t.Errorf("expected max mood 4, got %v", summary.MaxMood)
	}

	if summary.Average != 3 {
		t.Errorf("expected average 3, got %v", summary.Average)
	}
}

func TestCalculateStatsSummaryWhenEmpty(t *testing.T) {
	entries := []domain.MoodEntry{}

	_, err := CalculateStatsSummary(entries)

	if err == nil {
		t.Fatal("expected error for empty entries")
	}
}

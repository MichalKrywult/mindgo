package stats_test

import (
	"testing"

	"github.com/MichalKrywult/mindgo/internal/domain"
	"github.com/MichalKrywult/mindgo/internal/stats"
)

func TestCalculateStatsSummary(t *testing.T) {
	entries := []domain.MoodEntry{
		{Mood: 2, Note: "Test"},
		{Mood: 4, Note: "Test"},
	}

	result, err := stats.CalculateStats(entries)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	summary := result.Summary

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

	if result.Dist[2] != 1 {
		t.Errorf("expected mood 2 count 1, got %v", result.Dist[2])
	}

}

func TestCalculateStatsSummaryWhenEmpty(t *testing.T) {
	entries := []domain.MoodEntry{}

	_, err := stats.CalculateStats(entries)

	if err == nil {
		t.Fatal("expected error for empty entries")
	}
}

func TestRenderHistogram(t *testing.T) {
	dist := map[int]int{
		1:  3,
		5:  1,
		8:  5,
		10: 1,
	}

	histogram, err := stats.RenderHistogram(
		dist,
		domain.MinMoodValue,
		domain.MaxMoodValue,
	)

	if err != nil {
		t.Fatalf("something went wrong with formatting histogram: %v", err)
	}

	want := ` 1 | ███ (3)
 2 |  (0)
 3 |  (0)
 4 |  (0)
 5 | █ (1)
 6 |  (0)
 7 |  (0)
 8 | █████ (5)
 9 |  (0)
10 | █ (1)
`

	if histogram != want {
		t.Errorf("expected %q, got %q", want, histogram)
	}
}

func TestRenderHistogramWhenEmpty(t *testing.T) {
	dist := map[int]int{}

	_, err := stats.RenderHistogram(
		dist,
		domain.MinMoodValue,
		domain.MaxMoodValue,
	)

	if err == nil {
		t.Fatal("expected an error for empty map")
	}
}

package stats_test

import (
	"reflect"
	"testing"
	"time"

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

func TestCalculateStatsByWeekday(t *testing.T) {
	tests := []struct {
		name     string
		entries  []domain.MoodEntry
		expected map[time.Weekday]float64
	}{
		{
			name: "multiple entries on the same weekday",
			entries: []domain.MoodEntry{
				{Date: time.Date(2026, 9, 7, 0, 0, 0, 0, time.UTC), Mood: 4},
				{Date: time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC), Mood: 3},
				{Date: time.Date(2026, 9, 8, 0, 0, 0, 0, time.UTC), Mood: 5},
			},
			expected: map[time.Weekday]float64{
				time.Monday:  3.5,
				time.Tuesday: 5,
			},
		},
		{
			name:     "empty list",
			entries:  []domain.MoodEntry{},
			expected: map[time.Weekday]float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statsByWeekday := stats.CalculateStatsByWeekday(tt.entries)

			if !reflect.DeepEqual(statsByWeekday, tt.expected) { //allows to compare maps
				t.Errorf("expected %v, got %v", tt.expected, statsByWeekday)
			}

		})
	}
}

func TestCalculateStatsByTag(t *testing.T) {
	tests := []struct {
		name     string
		entries  []domain.MoodEntry
		expected map[string]stats.Summary
	}{
		{
			name: "multiple entries, without duplicate tags in single entry",
			entries: []domain.MoodEntry{
				{
					Mood: 5,
					Tags: []string{"work", "exercise"},
				},
				{
					Mood: 3,
					Tags: []string{"work"},
				},
				{
					Mood: 4,
					Tags: []string{"exercise"},
				}},
			expected: map[string]stats.Summary{
				"work": {
					TotalCount: 2,
					Average:    4,
					MinMood:    3,
					MaxMood:    5,
				},
				"exercise": {
					TotalCount: 2,
					Average:    4.5,
					MinMood:    4,
					MaxMood:    5,
				},
			},
		},
		{
			name: "multiple entries, with duplicate tags in single entry",
			entries: []domain.MoodEntry{
				{
					Mood: 5,
					Tags: []string{"work", "exercise", "work", "work"},
				},
				{
					Mood: 3,
					Tags: []string{"work"},
				},
				{
					Mood: 4,
					Tags: []string{"exercise"},
				}},
			expected: map[string]stats.Summary{
				"work": {
					TotalCount: 2,
					Average:    4,
					MinMood:    3,
					MaxMood:    5,
				},
				"exercise": {
					TotalCount: 2,
					Average:    4.5,
					MinMood:    4,
					MaxMood:    5,
				},
			},
		},
		{
			name: "no tags",
			entries: []domain.MoodEntry{
				{
					Mood: 5,
					Tags: []string{},
				},
			},
			expected: map[string]stats.Summary{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stats.CalculateStatsByTag(tt.entries)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

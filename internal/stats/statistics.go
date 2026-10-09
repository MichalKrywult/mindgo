package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

type Summary struct {
	TotalCount int
	Average    float64
	MinMood    int
	MaxMood    int
}

type Stats struct {
	Summary Summary
	Dist    map[int]int
}

func CalculateStats(entries []domain.MoodEntry) (Stats, error) {
	if len(entries) == 0 {
		return Stats{}, fmt.Errorf("calculating stats is not possible when entries list is empty")
	}

	average := 0.0
	minMood := entries[0].Mood
	maxMood := entries[0].Mood
	dist := make(map[int]int) // make is necessary in order to add entries

	for _, entry := range entries {
		average += float64(entry.Mood)
		dist[entry.Mood]++

		if entry.Mood < minMood {
			minMood = entry.Mood
		}

		if entry.Mood > maxMood {
			maxMood = entry.Mood
		}
	}

	return Stats{
		Summary: Summary{
			TotalCount: len(entries),
			MinMood:    minMood,
			MaxMood:    maxMood,
			Average:    average / float64(len(entries)),
		},
		Dist: dist,
	}, nil
}

func RenderHistogram(dist map[int]int, minScale, maxScale int) (string, error) {

	if len(dist) == 0 {
		return "", fmt.Errorf("length of dist cannot be 0")
	}
	var builder strings.Builder // builder from strings standard library

	for mood := minScale; mood <= maxScale; mood++ {
		count := dist[mood]
		bars := strings.Repeat("█", count)

		fmt.Fprintf(&builder, "%2d | %s (%d)\n", mood, bars, count)
		//%2d means that the length of d is 2 digits long - ' 3, '10' etc
	}

	return builder.String(), nil
}

func CalculateStatsByWeekday(entries []domain.MoodEntry) map[time.Weekday]float64 {
	weekdays := make(map[time.Weekday]struct { // anonymoous struct
		Sum   float64
		Count int
	})

	for _, entry := range entries {
		day := entry.Date.Weekday() //we can modify only 'local' copy

		stats := weekdays[day]
		stats.Sum += float64(entry.Mood)
		stats.Count++

		weekdays[day] = stats
	}

	averages := make(map[time.Weekday]float64)

	for day, stats := range weekdays {
		averages[day] = stats.Sum / float64(stats.Count)
	}

	return averages
}

func CalculateStatsByTag(entries []domain.MoodEntry) map[string]Summary {

	tagStats := make(map[string]struct {
		Sum   float64
		Count int
		Min   int
		Max   int
	})

	for _, entry := range entries {
		seen := make(map[string]bool)

		for _, tag := range entry.Tags {
			if seen[tag] == true {
				continue
			}

			seen[tag] = true

			stats, exists := tagStats[tag]
			mood := entry.Mood
			stats.Sum += float64(mood)
			stats.Count++

			if !exists {
				stats.Max = mood
				stats.Min = mood
				tagStats[tag] = stats
				continue
			}

			if mood > stats.Max {
				stats.Max = mood
			}

			if mood < stats.Min {
				stats.Min = mood
			}

			tagStats[tag] = stats
		}
	}

	tags := make(map[string]Summary)
	for tag, stats := range tagStats {
		tags[tag] = Summary{
			TotalCount: stats.Count,
			MaxMood:    stats.Max,
			MinMood:    stats.Min,
			Average:    stats.Sum / float64(stats.Count),
		}
	}

	return tags
}

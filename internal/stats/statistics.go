package stats

import (
	"fmt"
	"strings"

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

		builder.WriteString(fmt.Sprintf("%2d | %s (%d)\n", mood, bars, count))
		//%2d means that the length of d is 2 digits long - ' 3, '10' etc
	}

	return builder.String(), nil
}

package stats

import (
	"fmt"

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

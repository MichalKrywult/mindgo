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

func CalculateStatsSummary(entries []domain.MoodEntry) (Summary, error) {

	if len(entries) == 0 {
		return Summary{}, fmt.Errorf("calculating summary is not possible when entries list is empty")
	}

	average := 0.0
	minMood := 11
	maxMood := -1

	for _, entry := range entries {
		average += float64(entry.Mood)

		if entry.Mood < minMood {
			minMood = entry.Mood
		}
		if entry.Mood > maxMood {
			maxMood = entry.Mood
		}
	}

	average /= float64(len(entries))
	totalCount := len(entries)

	return Summary{TotalCount: totalCount,
		MinMood: minMood,
		MaxMood: maxMood,
		Average: average}, nil
}

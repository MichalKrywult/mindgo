package stats

import (
	"fmt"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

func CalculateAverageMood(entries []domain.MoodEntry) (float64, error) {
	if len(entries) == 0 {
		return 0, fmt.Errorf("calculating average is not possible when entries list is empty")
	}

	sum := 0.0
	for _, entry := range entries {
		sum += float64(entry.Mood)
	}

	return sum / float64(len(entries)), nil
}
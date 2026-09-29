package domain

import (
	"errors"
	"time"
)

const (
	MinMoodValue = 1
	MaxMoodValue = 10
)

type MoodEntry struct {
	ID   int       `json:"id"`
	Mood int       `json:"mood"`
	Date time.Time `json:"date"`
	Note string    `json:"note"`
	Tags []string  `json:"tags"`
}

func NewMoodEntry(mood int, date time.Time, note string, tags []string) (MoodEntry, error) {
	if mood < MinMoodValue || mood > MaxMoodValue {
		return MoodEntry{}, errors.New("invalid Mood range")
	}

	moodEntry := MoodEntry{Mood: mood, Date: date, Note: note, Tags: tags}
	return moodEntry, nil
}

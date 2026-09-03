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
}

func NewMoodEntry(mood int, date time.Time, note string) (MoodEntry, error) {
	if mood < MinMoodValue || mood > MaxMoodValue {
		return MoodEntry{}, errors.New("invalid Mood range")
	}

	moodEntry := MoodEntry{Mood: mood, Date: date, Note: note}
	return moodEntry, nil
}

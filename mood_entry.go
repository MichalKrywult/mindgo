package main

import (
	"errors"
	"time"
)

type MoodEntry struct {
	ID   int
	Mood int
	Date time.Time
	Note string
}

func NewMoodEntry(mood int, date time.Time, note string) (MoodEntry, error) {

	if mood < 1 || mood > 10 {
		return MoodEntry{}, errors.New("Invalid Mood range")
	}

	moodEntry := MoodEntry{Mood: mood, Date: date, Note: note}
	return moodEntry, nil
}

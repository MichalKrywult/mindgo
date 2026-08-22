package main

import "errors"

type MoodEntry struct {
	ID   int
	Mood int
	Date string
	Note string
}

func NewMoodEntry(mood int, date, note string) (MoodEntry, error) {

	if mood < 1 || mood > 10 {
		return MoodEntry{}, errors.New("Invalid Mood range")
	}

	if date == "" {
		return MoodEntry{}, errors.New("Date cannot be empty")
	}

	moodEntry := MoodEntry{Mood: mood, Date: date, Note: note}
	return moodEntry, nil
}

package main

import "fmt"

func (tracker MoodTracker) CalculateAverageMood() (float64, error) {

	if len(tracker.entries) == 0 {
		return 0, fmt.Errorf("Calculating average is not possible when tracker is empty")
	}

	average := 0.0
	for _, entry := range tracker.entries {
		average += float64(entry.Mood)
		// entry.Mood is currently an int, so it has to be converted to float64 type
	}
	average = average / float64(len(tracker.entries))

	return average, nil
}

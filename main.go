package main

import "fmt"

func main() {
	entry, err := NewMoodEntry(2, "2022/12/13", "Test")
	if err != nil {
		fmt.Println("Something went wrong:", err)
		return
	}

	tracker := MoodTracker{}
	tracker.AddEntry(entry)
	tracker.AddEntry(entry)

	fmt.Println(tracker.GetEntries())
	fmt.Println("Mood created and added correctly")
	fmt.Println(tracker.findIndexByID(2))
	fmt.Println(tracker.findEntryByID(2))

	err = tracker.RemoveEntry(1)
	if err != nil {
		fmt.Println("Couldn't remove that entry")
	} else {
		fmt.Println("Entry removed correctly")
	}
	showCLI()
}

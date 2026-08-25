package main

import "fmt"

func displayMenuAndReadChoice() int {
	fmt.Println("=====MENU=====")
	fmt.Println("1. New entry")
	fmt.Println("2. Remove entry")
	fmt.Println("3. Show history")
	fmt.Println("4. Show statistics")
	fmt.Println("0. Exit")

	fmt.Println("Your choice: ")
	var choice int
	fmt.Scan(&choice) //& means using adress of a variable, Scan saves input to that variable

	return choice
}

func displayAllEntries(tracker MoodTracker) {

	entries := tracker.GetEntries()

	if len(entries) == 0 {
		fmt.Print("You don't have any entries\n")
		return
	}

	fmt.Print("All of your entries:\n")
	for i := range entries {
		fmt.Printf("%d. %v | %v | %v\n", i+1, entries[i].Mood, entries[i].Date, entries[i].Note)
	}
}

func displayAverageMood(tracker MoodTracker) {
	average, err := tracker.CalculateAverageMood()
	if err != nil {
		fmt.Print("You don't have any entries\n")
		return
	}

	fmt.Printf("Your average mood: %v\n", average)
}

func readMoodEntry() (MoodEntry, error) {
	var entry MoodEntry

	fmt.Println("Your mood score:")
	fmt.Scan(&entry.Mood)

	fmt.Println("Date:")
	fmt.Scan(&entry.Date)

	fmt.Println("Note:")
	fmt.Scan(&entry.Note)

	entry, err := NewMoodEntry(entry.Mood, entry.Date, entry.Note)
	return entry, err
}

func (tracker *MoodTracker) GetEntryID(index int) (int, error) {
	if index < 0 || index >= len(tracker.entries) {
		return 0, fmt.Errorf("invalid entry index")
	}

	return tracker.entries[index].ID, nil
}

func showCLI(tracker *MoodTracker) { // *MoodTracker means tracker is a pointer to a MoodTracker type

	for {
		choice := displayMenuAndReadChoice()
		switch choice {
		case 1:
			entry, err := readMoodEntry()
			if err != nil {
				fmt.Println("Something went wrong:", err)
				return
			}

			tracker.AddEntry(entry)
			fmt.Println("Entry added!")

		case 2:
			fmt.Println("Which entry would you like to delete?")
			displayAllEntries(*tracker)

			fmt.Println("Number of entry to delete:")
			var choice int
			fmt.Scan(&choice)

			entryID, err := tracker.GetEntryID(choice - 1)
			if err != nil {
				fmt.Println("Invalid entry number")
				continue
			}

			err = tracker.RemoveEntryByID(entryID)
			if err != nil {
				fmt.Println("Unexpected error:", err)
				return
			}

			fmt.Println("Entry removed successfully")

		case 3:
			fmt.Println("Choice 3")
			displayAllEntries(*tracker) // passing a copy of a tracker to the function
		case 4:
			fmt.Println("Choice 4")
			displayAverageMood(*tracker)
		case 0:
			fmt.Println("Exit")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

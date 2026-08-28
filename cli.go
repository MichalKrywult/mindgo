package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type CLI struct {
	tracker *MoodTracker
	scanner *bufio.Scanner
}

func NewCLI(tracker *MoodTracker, input io.Reader) *CLI {
	return &CLI{
		tracker: tracker,
		scanner: bufio.NewScanner(input),
	}
}

func (cli *CLI) readLine() (string, error) {
	if !cli.scanner.Scan() {
		return "", cli.scanner.Err()
	}

	return strings.TrimSpace(cli.scanner.Text()), nil
}

func (cli *CLI) readInt() (int, error) {
	text, err := cli.readLine()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(text)
}

func (cli *CLI) displayMenuAndReadChoice() int {
	fmt.Println("=====MENU=====")
	fmt.Println("1. New entry")
	fmt.Println("2. Edit entry")
	fmt.Println("3. Remove entry")
	fmt.Println("4. Show history")
	fmt.Println("5. Show statistics")
	fmt.Println("0. Exit")

	fmt.Print("Your choice: ")
	choice, err := cli.readInt()
	if err != nil {
		return -1
	}

	return choice
}

func displayAllEntries(tracker MoodTracker) error {

	entries := tracker.GetEntries()

	if len(entries) == 0 {
		return fmt.Errorf("You don't have any entries\n")
	}

	fmt.Print("All of your entries:\n")
	for i := range entries {
		fmt.Printf(
			"%d. %d | %s | %s\n",
			i+1,
			entries[i].Mood,
			entries[i].Date.Format("2006-01-02 15:04"),
			entries[i].Note,
		)

	}
	return nil
}

func displayAverageMood(tracker MoodTracker) {
	average, err := tracker.CalculateAverageMood()
	if err != nil {
		fmt.Print("You don't have any entries\n")
		return
	}

	fmt.Printf("Your average mood: %v\n", average)
}

func (cli *CLI) readNewMoodEntry() (MoodEntry, error) {
	fmt.Print("Mood score (1-10): ")

	mood, err := cli.readInt()
	if err != nil {
		return MoodEntry{}, err
	}

	fmt.Print("Note: ")

	note, err := cli.readLine()
	if err != nil {
		return MoodEntry{}, err
	}

	return NewMoodEntry(mood, time.Now(), note)
}

func (tracker MoodTracker) getEntryID(index int) (int, error) {
	if index < 0 || index >= len(tracker.entries) {
		return 0, fmt.Errorf("invalid entry index")
	}

	return tracker.entries[index].ID, nil
}

func (cli *CLI) show() {
	for {
		choice := cli.displayMenuAndReadChoice()

		switch choice {
		case 1:
			entry, err := cli.readNewMoodEntry()
			if err != nil {
				fmt.Println("Something went wrong:", err)
				continue
			}

			cli.tracker.AddEntry(entry)
			fmt.Println("Entry added!")

		case 2:
			fmt.Println("Which entry would you like to edit?")
			err := displayAllEntries(*cli.tracker)
			if err != nil {
				fmt.Print(err)
				continue
			}

			fmt.Println("Number of entry to edit:")
			choice, err := cli.readInt()
			if err != nil {
				fmt.Println("Invalid entry")
				continue
			}

			entryID, err := cli.tracker.getEntryID(choice - 1)
			if err != nil {
				fmt.Println("Invalid entry number")
				continue
			}

			entry, err := cli.readNewMoodEntry()
			if err != nil {
				fmt.Println("Something went wrong:", err)
				continue
			}

			err = cli.tracker.EditEntryByID(entryID, entry)
			if err != nil {
				fmt.Println("Something went wrong:", err)
				continue
			}

		case 3:
			fmt.Println("Which entry would you like to delete?")
			err := displayAllEntries(*cli.tracker)
			if err != nil {
				fmt.Print(err)
				continue
			}

			fmt.Println("Number of entry to delete:")
			choice, err := cli.readInt()
			if err != nil {
				fmt.Println("Invalid entry")
				continue
			}

			entryID, err := cli.tracker.getEntryID(choice - 1)
			if err != nil {
				fmt.Println("Invalid entry number")
				continue
			}

			err = cli.tracker.RemoveEntryByID(entryID)
			if err != nil {
				fmt.Println("Unexpected error:", err)
				continue
			}

			fmt.Println("Entry removed successfully")

		case 4:
			err := displayAllEntries(*cli.tracker) // passing a copy of the MoodTracker to the function

			if err != nil {
				fmt.Print(err)
				continue
			}

		case 5:
			displayAverageMood(*cli.tracker)

		case 0:
			fmt.Println("Exit")
			return

		default:
			fmt.Println("Invalid choice")
		}
	}
}

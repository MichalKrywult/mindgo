package cli

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/MichalKrywult/mindgo/internal/domain"
	"github.com/MichalKrywult/mindgo/internal/tracker"
)

type CLI struct {
	tracker *tracker.MoodTracker
	scanner *bufio.Scanner
}

func NewCLI(tracker *tracker.MoodTracker, input io.Reader) *CLI {
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

func (cli *CLI) hasEntries() bool {
	return len(cli.tracker.GetEntries()) > 0
}

func (cli *CLI) IsIndexValid(index int) bool {
	return index >= 0 && index < len(cli.tracker.GetEntries())
}

func (cli *CLI) displayAllEntries() {
	entries := cli.tracker.GetEntries()

	if len(entries) == 0 {
		fmt.Println("You don't have any entries")
		return
	}

	fmt.Println("All of your entries:")
	for i := range entries {
		fmt.Printf(
			"%d. %d | %s | %s\n",
			i+1,
			entries[i].Mood,
			entries[i].Date.Format("2006-01-02 15:04"),
			entries[i].Note,
		)
	}
}

func (cli *CLI) displayAverageMood() {
	average, err := cli.tracker.CalculateAverageMood()
	if err != nil {
		fmt.Println("You don't have any entries")
		return
	}

	fmt.Printf("Your average mood: %v\n", average)
}

func (cli *CLI) readNewMoodEntry() (domain.MoodEntry, error) {
	fmt.Print("Mood score (1-10): ")

	mood, err := cli.readInt()
	if err != nil {
		return domain.MoodEntry{}, err
	}

	fmt.Print("Note: ")

	note, err := cli.readLine()
	if err != nil {
		return domain.MoodEntry{}, err
	}

	return domain.NewMoodEntry(mood, time.Now(), note)
}

func (cli *CLI) Show() {
	for {
		choice := cli.displayMenuAndReadChoice()

		switch choice {
		case 1:
			entry, err := cli.readNewMoodEntry()
			if err != nil {
				fmt.Println("Something went wrong:", err)
				continue
			}

			err = cli.tracker.AddEntry(entry)
			if err != nil {
				fmt.Println("Unexpected error occured:", err)
				continue
			}
			fmt.Println("Entry added!")

		case 2:
			if !cli.hasEntries() {
				fmt.Println("You don't have any entries")
				continue
			}

			fmt.Println("Which entry would you like to edit?")
			cli.displayAllEntries()

			fmt.Println("Number of entry to edit:")
			choice, err := cli.readInt()
			if err != nil {
				fmt.Println("Invalid entry")
				continue
			}
			index := choice - 1

			if !cli.IsIndexValid(index) {
				fmt.Println("Invalid entry number")
				continue
			}

			entry, err := cli.readNewMoodEntry()
			if err != nil {
				fmt.Printf("Something went wrong with adding new entry: %v", err)
				continue
			}

			err = cli.tracker.EditEntryByIndex(index, entry)
			if err != nil {
				fmt.Printf("Something went wrong with editing entry: %v", err)
				continue
			}

		case 3:
			if !cli.hasEntries() {
				fmt.Println("You don't have any entries")
				continue
			}

			fmt.Println("Which entry would you like to delete?")
			cli.displayAllEntries()

			fmt.Println("Number of entry to delete:")
			choice, err := cli.readInt()
			if err != nil {
				fmt.Println("Invalid entry")
				continue
			}
			index := choice - 1

			if !cli.IsIndexValid(index) {
				fmt.Println("Invalid entry number")
				continue
			}

			err = cli.tracker.RemoveEntryByIndex(index)
			if err != nil {
				fmt.Println("Something went wrong with removing entry:", err)
				continue
			}

		case 4:
			cli.displayAllEntries()

		case 5:
			cli.displayAverageMood()

		case 0:
			fmt.Println("Exit")
			return

		default:
			fmt.Println("Invalid choice")
		}
	}
}

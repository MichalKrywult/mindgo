package cli

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/MichalKrywult/mindgo/internal/domain"
	"github.com/MichalKrywult/mindgo/internal/export"
	"github.com/MichalKrywult/mindgo/internal/stats"
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
		if cli.scanner.Err() == nil { // in that case EOF occurs
			return "", io.EOF
		}
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

func (cli *CLI) displayMenuAndReadChoice() (int, error) {
	fmt.Println("=====MENU=====")
	fmt.Println("1. New entry")
	fmt.Println("2. Edit entry")
	fmt.Println("3. Remove entry")
	fmt.Println("4. Show history")
	fmt.Println("5. Show statistics")
	fmt.Println("6. Export to csv")
	fmt.Println("0. Exit")

	fmt.Print("Your choice: ")
	choice, err := cli.readInt()
	if err != nil {
		return 0, err
	}

	return choice, nil
}

func (cli *CLI) hasEntries() bool {
	return len(cli.tracker.GetEntries()) > 0
}

func (cli *CLI) isIndexValid(index int) bool {
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

func (cli *CLI) handleAddEntry() error {

	entry, err := cli.readNewMoodEntry()
	if err != nil {
		return fmt.Errorf("Something went wrong: %v", err)
	}

	err = cli.tracker.AddEntry(entry)
	if err != nil {
		return fmt.Errorf("Unexpected error occured: %v", err)
	}
	return nil
}

func (cli *CLI) handleEditEntry() error {
	if !cli.hasEntries() {
		return fmt.Errorf("You don't have any entries")
	}

	fmt.Println("Which entry would you like to edit?")
	cli.displayAllEntries()

	fmt.Println("Number of entry to edit:")
	choice, err := cli.readInt()
	if err != nil {
		return fmt.Errorf("Invalid entry")
	}
	index := choice - 1

	if !cli.isIndexValid(index) {
		return fmt.Errorf("Invalid entry number")
	}

	entry, err := cli.readNewMoodEntry()
	if err != nil {
		return fmt.Errorf("Something went wrong with adding new entry: %v", err)
	}

	err = cli.tracker.EditEntryByIndex(index, entry)
	if err != nil {
		return fmt.Errorf("Something went wrong with editing entry: %v\n", err)
	}

	return nil
}

func (cli *CLI) handleRemoveEntry() error {

	if !cli.hasEntries() {
		return fmt.Errorf("You don't have any entries")
	}

	fmt.Println("Which entry would you like to delete?")
	cli.displayAllEntries()

	fmt.Println("Number of entry to delete:")
	choice, err := cli.readInt()
	if err != nil {
		return fmt.Errorf("Invalid entry")

	}
	index := choice - 1

	if !cli.isIndexValid(index) {
		return fmt.Errorf("Invalid entry number")
	}

	err = cli.tracker.RemoveEntryByIndex(index)
	if err != nil {
		return fmt.Errorf("Something went wrong with removing entry: %v", err)

	}

	return nil
}

func (cli *CLI) handleDisplayingStatistics() error {
	entries := cli.tracker.GetEntries()
	statsData, err := stats.CalculateStats(entries)
	if err != nil {
		return fmt.Errorf("Something went wrong: %v", err)

	}

	fmt.Printf("Total count: %d\n", statsData.Summary.TotalCount)
	fmt.Printf("Average: %.2f (Min: %d, Max: %d)\n",
		statsData.Summary.Average,
		statsData.Summary.MinMood,
		statsData.Summary.MaxMood)

	fmt.Println("Histogram:")
	histogram, err := stats.RenderHistogram(
		statsData.Dist,
		domain.MinMoodValue,
		domain.MaxMoodValue)

	if err != nil {
		return fmt.Errorf("Something went wrong: %v", err)
	}

	fmt.Print(histogram)

	fmt.Println("Average mood for each day:")
	statsByDay := stats.CalculateStatsByWeekday(entries)
	days := []time.Weekday{ //for some reason days in Go start from Sunday,
		//so it's needed to create slice of days
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
		time.Sunday,
	}

	for _, day := range days {
		average, exists := statsByDay[day]

		if !exists {
			fmt.Printf("%s: N/A\n", day)
			continue
		}

		fmt.Printf("%s: %.2f\n", day, average)
	}
	return nil
}

func (cli *CLI) handleExport() error {

	entries := cli.tracker.GetEntries()

	err := export.ExportEntriesToCSV(entries)
	if err != nil {
		return fmt.Errorf("Unexpected error occured: %v", err)

	}

	return nil
}

func (cli *CLI) Show() {
	for {
		choice, err := cli.displayMenuAndReadChoice()
		if err != nil {
			if err == io.EOF {
				return
			}

			fmt.Printf("Invalid choice: %v\n", err)
			continue
		}

		switch choice {
		case 1:
			err := cli.handleAddEntry()
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("Entry added!")
		case 2:
			err := cli.handleEditEntry()
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("Entry edited!")
		case 3:
			err := cli.handleRemoveEntry()
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("Entry removed!")
		case 4:
			cli.displayAllEntries()

		case 5:
			err := cli.handleDisplayingStatistics()
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("Statistics displayed!")
		case 6:
			err := cli.handleExport()
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("Entries exported succesfully")
		case 0:
			fmt.Println("Exit")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

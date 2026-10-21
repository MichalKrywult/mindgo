package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
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

type SortMode int

const (
	SortByDateNewest SortMode = iota // untyped int that is used in const, zero indexed
	SortByDateOldest
)

func NewCLI(tracker *tracker.MoodTracker, input io.Reader) *CLI {
	return &CLI{
		tracker: tracker,
		scanner: bufio.NewScanner(input),
	}
}

func (cli *CLI) readLine() (string, error) {
	if !cli.scanner.Scan() {
		if cli.scanner.Err() == nil {
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

func (cli *CLI) confirmAction(prompt string) (bool, error) {
	fmt.Println(prompt)
	fmt.Print("yes/no: ")

	input, err := cli.readLine()
	if err != nil {
		return false, err
	}
	if strings.EqualFold(input, "yes") || strings.EqualFold(input, "y") {
		return true, nil
	}

	return false, nil
}

func (cli *CLI) displayMainMenuAndReadChoice() (int, error) {
	fmt.Println("=====MENU=====")
	fmt.Println("1. New entry")
	fmt.Println("2. Edit entry")
	fmt.Println("3. Remove entry")
	fmt.Println("4. Show history")
	fmt.Println("5. Show statistics")
	fmt.Println("6. Export to csv")
	fmt.Println("7. REMOVE ALL ENTRIES")
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

func (cli *CLI) handleDisplayingAllEntries() error {
	if !cli.hasEntries() {
		return errors.New("you don't have any entries")
	}

	fmt.Println("How to sort them?")
	fmt.Println("1. Newest to oldest")
	fmt.Println("2. Oldest to newest")
	fmt.Println("0. Go back")

	fmt.Println("Your choice:")
	choice, err := cli.readInt()
	if err != nil {
		return errors.New("invalid choice")
	}

	switch choice {
	case 0:
		return nil
	case 1:
		_, err := cli.displayAllEntries(SortByDateNewest)
		return err
	case 2:
		_, err := cli.displayAllEntries(SortByDateOldest)
		return err
	default:
		return errors.New("invalid choice")
	}
}

func (cli *CLI) displayAllEntries(mode SortMode) (map[int]int, error) {

	entries := cli.tracker.GetEntries()

	if len(entries) == 0 {
		return map[int]int{}, errors.New("You don't have any entries")
	}

	switch mode {
	case SortByDateNewest:
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Date.After(entries[j].Date)
			//anonymous func that compares i,j (indexes in slice)
			//whether they need to be switched
		})
	case SortByDateOldest:
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Date.Before(entries[j].Date)
		})
	default:
		return map[int]int{}, errors.New("Invalid sort choice")
	}

	sortedIndexToID := make(map[int]int, len(entries))

	fmt.Println("All of your entries:")
	for i := range entries {
		fmt.Printf(
			"%d. %d | %s | %s | %s\n",
			i+1,
			entries[i].Mood,
			entries[i].Date.Format("2006-01-02 15:04"),
			entries[i].Note,
			strings.Join(entries[i].Tags, ", "),
		)

		sortedIndexToID[i+1] = entries[i].ID
	}

	return sortedIndexToID, nil
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

	fmt.Print("Tags (comma-separated): ")

	tagsInput, err := cli.readLine()
	if err != nil {
		return domain.MoodEntry{}, err
	}

	var tags []string
	for _, tag := range strings.Split(tagsInput, ",") {
		tag = strings.TrimSpace(tag)

		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return domain.NewMoodEntry(mood, time.Now(), note, tags)
}

func (cli *CLI) handleAddEntry() error {
	entry, err := cli.readNewMoodEntry()
	if err != nil {
		return fmt.Errorf("failed to read new entry: %w", err)
	}

	err = cli.tracker.AddEntry(entry)
	if err != nil {
		return fmt.Errorf("failed to add entry: %w", err)
	}

	return nil
}

func (cli *CLI) handleEditEntry() error {
	if !cli.hasEntries() {
		return errors.New("you don't have any entries")
	}

	fmt.Println("Which entry would you like to edit?")
	indexToID, err := cli.displayAllEntries(SortByDateNewest)
	if err != nil {
		return err
	}

	fmt.Println("Number of entry to edit:")
	choice, err := cli.readInt()
	if err != nil {
		return err
	}

	ID := indexToID[choice]

	ID, exists := indexToID[choice]
	if !exists {
		return errors.New("invalid entry number")
	}

	entry, err := cli.readNewMoodEntry()
	if err != nil {
		return fmt.Errorf("failed to read new entry: %w", err)
	}

	err = cli.tracker.EditEntryByID(ID, entry)
	if err != nil {
		return fmt.Errorf("failed to edit entry: %w", err)
	}

	return nil
}

func (cli *CLI) handleRemoveEntry() error {
	if !cli.hasEntries() {
		return errors.New("you don't have any entries")
	}

	fmt.Println("Which entry would you like to delete?")
	indexToID, err := cli.displayAllEntries(SortByDateNewest)
	if err != nil {
		return err
	}

	fmt.Println("Number of entry to delete:")
	choice, err := cli.readInt()
	if err != nil {
		return errors.New("invalid entry")
	}

	ID, exists := indexToID[choice]
	if !exists {
		return errors.New("invalid entry number")
	}
	
	err = cli.tracker.RemoveEntryByID(ID)
	if err != nil {
		return fmt.Errorf("failed to remove entry: %w", err)
	}

	return nil
}

func (cli *CLI) handleRemoveAllEntries() error {
	if !cli.hasEntries() {
		return errors.New("you don't have any entries")
	}

	prompt := `Are you sure that you want to remove all entries?
This action cannot be undone and your entries will be lost forever.`

	confirmed, err := cli.confirmAction(prompt)
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Println("Operation rejected, entries are safe!")
		return nil
	}

	err = cli.tracker.RemoveAllEntries()
	if err != nil {
		return fmt.Errorf("failed to remove all entries: %w", err)
	}

	fmt.Println("Entries removed!")
	return nil
}

func (cli *CLI) handleStatistics() error {
	entries := cli.tracker.GetEntries()

	statsData, err := stats.CalculateStats(entries)
	if err != nil {
		return fmt.Errorf("failed to calculate statistics: %w", err)
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
		domain.MaxMoodValue,
	)

	if err != nil {
		return fmt.Errorf("failed to render histogram: %w", err)
	}

	fmt.Print(histogram)

	fmt.Println("Average mood for each day:")

	statsByDay := stats.CalculateStatsByWeekday(entries)

	days := []time.Weekday{
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
		return fmt.Errorf("failed to export entries: %w", err)
	}

	return nil
}

func (cli *CLI) Show() {
	for {
		choice, err := cli.displayMainMenuAndReadChoice()
		if err != nil {
			if errors.Is(err, io.EOF) {
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
			err := cli.handleDisplayingAllEntries()
			if err != nil {
				fmt.Println(err)
			}

		case 5:
			err := cli.handleStatistics()
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

			fmt.Println("Entries exported successfully")

		case 7:
			err := cli.handleRemoveAllEntries()
			if err != nil {
				fmt.Println(err)
				continue
			}

		case 0:
			fmt.Println("Exit")
			return

		default:
			fmt.Println("Invalid choice")
		}
	}
}

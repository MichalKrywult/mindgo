package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

const (
	defaultExportCSVFilePath = "moods.csv"
)

func exportCSV(entries []domain.MoodEntry, w io.Writer) error {

	csWriter := csv.NewWriter(w)
	for _, entry := range entries {
		err := csWriter.Write([]string{
			strconv.Itoa(entry.ID),                 // converts int to string
			strconv.Itoa(entry.Mood),               // converts int to string
			entry.Date.Format("2006-01-02T15:04Z"), // converts date to string (fixed date format)
			entry.Note})
		if err != nil {
			return err
		}
	}
	csWriter.Flush() // function that empties the buffor
	err := csWriter.Error()
	if err != nil {
		return err
	}
	return nil
}

func createExportFile(fileName string) (*os.File, error) {

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("unexpected error occured: %w", err)
	}

	csvPath := filepath.Join(configDir, "mindgo", fileName)
	file, err := os.Create(csvPath)
	if err != nil {
		return nil, fmt.Errorf("unexpected error occured: %w", err)
	}

	return file, nil
}

func ExportEntriesToCSV(entries []domain.MoodEntry) error {

	file, err := createExportFile(defaultExportCSVFilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	err = exportCSV(entries, file)
	if err != nil {
		return err
	}

	return nil
}

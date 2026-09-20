package export

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/MichalKrywult/mindgo/internal/config"
	"github.com/MichalKrywult/mindgo/internal/domain"
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

func createExportFile() (*os.File, error) {

	csvPath, err := config.BuildExportPath()
	if err != nil {
		return nil, fmt.Errorf("unexpected error occured: %w", err)
	}

	file, err := os.Create(csvPath)
	if err != nil {
		return nil, fmt.Errorf("unexpected error occured: %w", err)
	}

	return file, nil
}


func ExportEntriesToCSV(entries []domain.MoodEntry) error {
	file, err := createExportFile()
	if err != nil {
		return err
	}

	exportErr := exportCSV(entries, file)
	closeErr := file.Close()

	if exportErr != nil {
		return exportErr
	}

	if closeErr != nil {
		return closeErr
	}

	return nil
}

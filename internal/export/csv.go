package export

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

func ExportCSV(entries []domain.MoodEntry, w io.Writer) error {

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

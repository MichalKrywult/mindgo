package export

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

type errorWriter struct{}

func (w errorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("unexpected error")

}

func TestExportCSV(t *testing.T) {

	tests := []struct {
		name     string
		entries  []domain.MoodEntry
		expected string
	}{
		{
			name: "correct line (single line)",
			entries: []domain.MoodEntry{
				{Date: time.Date(2026, 9, 7, 0, 0, 0, 0, time.UTC), Mood: 4, Note: "Test", ID: 1},
			},
			expected: "1,4,2026-09-07T00:00Z,Test\n",
		},
		{
			name: "correct input (multiple lines)",
			entries: []domain.MoodEntry{
				{Date: time.Date(2026, 9, 7, 0, 0, 0, 0, time.UTC), Mood: 4, Note: "Test", ID: 1},
				{Date: time.Date(2025, 1, 7, 0, 0, 0, 0, time.UTC), Mood: 8, Note: "Test2 ,", ID: 2},
				{Date: time.Date(124, 1, 7, 0, 0, 0, 0, time.UTC), Mood: 2, Note: "Test3 space", ID: 3},
			},
			expected: "1,4,2026-09-07T00:00Z,Test\n" +
				"2,8,2025-01-07T00:00Z,\"Test2 ,\"\n" +
				"3,2,0124-01-07T00:00Z,Test3 space\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer // place in the RAM memory, that pretends to be a file
			err := ExportCSV(tt.entries, &buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			data := buf.String()
			if data != tt.expected {
				t.Errorf("expected=%v, got=%v", tt.expected, data)
			}

		})
	}
}

func TestExportCSVWriterError(t *testing.T) {
	var buf errorWriter
	entries := []domain.MoodEntry{{Mood: 4, Note: "Test"}}
	err := ExportCSV(entries, &buf)
	if err == nil {
		t.Fatal("expected error for invalid writer")
	}
}

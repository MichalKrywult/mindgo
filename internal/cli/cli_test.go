package cli

import (
	"strings"
	"testing"
	"time"

	"github.com/MichalKrywult/mindgo/internal/domain"
	"github.com/MichalKrywult/mindgo/internal/storage"
	"github.com/MichalKrywult/mindgo/internal/tracker"
)

func TestIsIndexValid(t *testing.T) {
	tracker, err := tracker.NewMoodTracker(&storage.MockStorage{})
	if err != nil {
		t.Fatalf("failed to create tracker: %v", err)
	}

	input := strings.NewReader("0\n")
	cli := NewCLI(tracker, input)

	if cli.isIndexValid(0) {
		t.Error("expected isIndexValid(0) to be false for empty tracker")
	}

	entry := domain.MoodEntry{Mood: 5, Date: time.Now(), Note: "Test"}
	_ = tracker.AddEntry(entry)

	if !cli.isIndexValid(0) {
		t.Error("expected isIndexValid(0) to be true")
	}
	if cli.isIndexValid(1) {
		t.Error("expected isIndexValid(1) to be false")
	}
	if cli.isIndexValid(-1) {
		t.Error("expected isIndexValid(-1) to be false")
	}
}

func TestCLI(t *testing.T) {
	tests := []struct {
		name          string
		initialEntry  *domain.MoodEntry //pointers can be nil, values cannot
		input         string
		expectedCount int
		expectedMood  int
		expectedNote  string
	}{
		{
			name:          "test invalid input",
			initialEntry:  nil, // return the pointer
			input:         "abc\nabc\n",
			expectedCount: 0,
			expectedMood:  0,
			expectedNote:  "",
		},
		{
			name:          "test add entry",
			initialEntry:  nil,
			input:         "1\n8\nGreat day\n0\n",
			expectedCount: 1,
			expectedMood:  8,
			expectedNote:  "Great day",
		},
		{
			name:          "test edit entry",
			initialEntry:  &domain.MoodEntry{Mood: 2, Note: "Bad day"}, // return the pointer
			input:         "2\n1\n8\nGreat day\n0\n",
			expectedCount: 1,
			expectedMood:  8,
			expectedNote:  "Great day",
		},
		{
			name:          "test editing empty tracker",
			initialEntry:  nil, // return the pointer
			input:         "3\n0\n",
			expectedCount: 0,
			expectedMood:  0,
			expectedNote:  "",
		},
		{
			name:          "test removing entry",
			initialEntry:  &domain.MoodEntry{Mood: 2, Note: "Bad day"}, // return the pointer
			input:         "3\n1\n0\n",
			expectedCount: 0,
			expectedMood:  0,
			expectedNote:  "",
		},
		{
			name:          "test removing with invalid index",
			initialEntry:  nil, // return the pointer
			input:         "3\n0\n",
			expectedCount: 0,
			expectedMood:  0,
			expectedNote:  "",
		},
		{
			name:          "test removing with invalid input",
			initialEntry:  nil, // return the pointer
			input:         "3\nabc\n0\n",
			expectedCount: 0,
			expectedMood:  0,
			expectedNote:  "",
		},
		{
			name:          "test removing with invalid input",
			initialEntry:  nil, // return the pointer
			input:         "abc\nabc\n0\n",
			expectedCount: 0,
			expectedMood:  0,
			expectedNote:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker, err := tracker.NewMoodTracker(&storage.MockStorage{})
			if err != nil {
				t.Fatalf("failed to create tracker: %v", err)
			}

			if tt.initialEntry != nil {
				err := tracker.AddEntry(*tt.initialEntry)
				if err != nil {
					t.Fatalf("failed to add inital entry %v", err)
				}
			}

			cli := NewCLI(tracker, strings.NewReader(tt.input))
			cli.Show()

			entries := tracker.GetEntries()

			if len(entries) != tt.expectedCount {
				t.Fatalf("expected %d, got %d", tt.expectedCount, len(entries))
			}

			if len(entries) == 0 {
				return
			}

			if entries[0].Mood != tt.expectedMood {
				t.Errorf("expected mood %d, got %d", tt.expectedMood, entries[0].Mood)
			}

			if entries[0].Note != tt.expectedNote {
				t.Errorf("expected note %q, got %q", tt.expectedNote, entries[0].Note)
			}
		})
	}
}

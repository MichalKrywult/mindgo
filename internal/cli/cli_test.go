package cli

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MichalKrywult/mindgo/internal/config"
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
		name            string
		initialEntries  []domain.MoodEntry
		input           string
		expectedEntries []domain.MoodEntry
	}{
		{
			name:            "test invalid input",
			initialEntries:  nil,
			input:           "abc\nabc\n",
			expectedEntries: []domain.MoodEntry{},
		},
		{
			name:           "test add entry",
			initialEntries: nil,
			input:          "1\n8\nGreat day\n0\n",
			expectedEntries: []domain.MoodEntry{
				{Mood: 8, Note: "Great day"},
			},
		},
		{
			name: "test edit entry",
			initialEntries: []domain.MoodEntry{
				{Mood: 2, Note: "Bad day"},
			},
			input: "2\n1\n8\nGreat day\n0\n",
			expectedEntries: []domain.MoodEntry{
				{Mood: 8, Note: "Great day"},
			},
		},
		{
			name:            "test editing empty tracker",
			initialEntries:  nil,
			input:           "3\n0\n",
			expectedEntries: []domain.MoodEntry{},
		},
		{
			name: "test removing entry",
			initialEntries: []domain.MoodEntry{
				{Mood: 2, Note: "Bad day"},
			},
			input:           "3\n1\n0\n",
			expectedEntries: []domain.MoodEntry{},
		},
		{
			name:            "test removing with invalid index",
			initialEntries:  nil,
			input:           "3\n0\n",
			expectedEntries: []domain.MoodEntry{},
		},
		{
			name: "test removing ALL entries success",
			initialEntries: []domain.MoodEntry{
				{Mood: 2, Note: "Bad day"},
				{Mood: 2, Note: "Bad day"},
				{Mood: 2, Note: "Bad day"},
			},
			input:           "7\nyes\n0\n",
			expectedEntries: []domain.MoodEntry{},
		},
		{
			name: "test removing ALL entries rejection",
			initialEntries: []domain.MoodEntry{
				{Mood: 2, Note: "Bad day"},
				{Mood: 2, Note: "Bad day"},
				{Mood: 2, Note: "Bad day"},
			},
			input: "7\nno\n0\n",
			expectedEntries: []domain.MoodEntry{
				{Mood: 2, Note: "Bad day"},
				{Mood: 2, Note: "Bad day"},
				{Mood: 2, Note: "Bad day"},
			},
		},
		{
			name: "test removing ALL entries invalid confirmation",
			initialEntries: []domain.MoodEntry{
				{Mood: 2, Note: "Bad day"},
				{Mood: 2, Note: "Bad day"},
				{Mood: 2, Note: "Bad day"},
			},
			input: "7\nabc\n0\n",
			expectedEntries: []domain.MoodEntry{
				{Mood: 2, Note: "Bad day"},
				{Mood: 2, Note: "Bad day"},
				{Mood: 2, Note: "Bad day"},
			},
		},
		{
			name:            "test removing ALL entries empty tracker",
			initialEntries:  []domain.MoodEntry{},
			input:           "7\n0\n",
			expectedEntries: []domain.MoodEntry{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker, err := tracker.NewMoodTracker(&storage.MockStorage{})
			if err != nil {
				t.Fatalf("failed to create tracker: %v", err)
			}

			for _, entry := range tt.initialEntries {
				err := tracker.AddEntry(entry)
				if err != nil {
					t.Fatalf("failed to add initial entry: %v", err)
				}
			}

			cli := NewCLI(tracker, strings.NewReader(tt.input))
			cli.Show()

			entries := tracker.GetEntries()

			if len(entries) != len(tt.expectedEntries) {
				t.Fatalf("expected %d entries, got %d",
					len(tt.expectedEntries),
					len(entries))
			}

			for i, expected := range tt.expectedEntries {
				actual := entries[i]

				if actual.Mood != expected.Mood {
					t.Errorf("entry %d: expected mood %d, got %d",
						i,
						expected.Mood,
						actual.Mood)
				}

				if actual.Note != expected.Note {
					t.Errorf("entry %d: expected note %q, got %q",
						i,
						expected.Note,
						actual.Note)
				}
			}
		})
	}
}

func TestPersistentDataPath(t *testing.T) {

	//first run
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")

	parsedFlags := flags{path: filepath.Join(dir, "path.json")}
	conf := config.Config{}

	correctPath, err := SelectDataPath(parsedFlags, conf)
	if err != nil {
		t.Fatalf("failed to select path: %v", err)
	}

	conf.DataPath = correctPath
	err = config.SaveConfigToFile(configPath, conf)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	//second run
	newConf, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	loadedPath, err := SelectDataPath(flags{}, newConf)
	if err != nil {
		t.Fatalf("failed to select path: %v", err)
	}

	if loadedPath != correctPath {
		t.Fatalf("paths are not equal, got: %v, wanted: %v", loadedPath, correctPath)
	}
}

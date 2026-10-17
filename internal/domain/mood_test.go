package domain

import (
	"testing"
	"time"
)

func TestNewMoodEntryDoesNotShareTags(t *testing.T) {
	tags := []string{"test1", "test2"}

	entry, err := NewMoodEntry(2, time.Now(), "Test", tags)
	if err != nil {
		t.Fatalf("error when creating entry: %v", err)
	}

	tags[0] = "test3"

	if entry.Tags[0] != "test1" {
		t.Errorf("entry tags changed: got %q, want %q", entry.Tags[0], "test1")
	}
}

func TestNewMoodEntry(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		mood     int
		note     string
		tags     []string
		wantErr  bool
		wantMood int
		wantNote string
		wantTags []string
	}{
		{
			name:     "valid mood",
			mood:     5,
			note:     "Test note",
			tags:     []string{"test"},
			wantMood: 5,
			wantNote: "Test note",
			wantTags: []string{"test"},
		},
		{
			name:    "mood below minimum",
			mood:    0,
			note:    "Test",
			wantErr: true,
		},
		{
			name:    "mood above maximum",
			mood:    11,
			note:    "Test",
			wantErr: true,
		},
		{
			name:    "negative mood",
			mood:    -1,
			note:    "Test",
			wantErr: true,
		},
		{
			name:    "duplicate tags",
			mood:    5,
			note:    "Test",
			tags:    []string{"test", "test"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := NewMoodEntry(tt.mood, now, tt.note, tt.tags)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if entry.Mood != tt.wantMood {
				t.Errorf("got mood %d, want %d", entry.Mood, tt.wantMood)
			}

			if entry.Note != tt.wantNote {
				t.Errorf("got note %q, want %q", entry.Note, tt.wantNote)
			}

			if len(entry.Tags) != len(tt.wantTags) {
				t.Fatalf("got tags %v, want %v", entry.Tags, tt.wantTags)
			}

			for i, tag := range tt.wantTags {
				if entry.Tags[i] != tag {
					t.Errorf("got tag %q at index %d, want %q", entry.Tags[i], i, tag)
				}
			}
		})
	}
}

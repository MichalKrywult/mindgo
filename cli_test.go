package main

import (
	"strings"
	"testing"
)

func TestCLIAddsMoodEntry(t *testing.T) {
	input := strings.NewReader("1\n2\nGreat day\n0\n")

	tracker := MoodTracker{}
	cli := NewCLI(&tracker, input)

	cli.show()

	entries := tracker.GetEntries()

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].Mood != 2 {
		t.Errorf("expected mood 2, got %d", entries[0].Mood)
	}

	if entries[0].Note != "Great day" {
		t.Errorf("expected note %q, got %q", "Great day", entries[0].Note)
	}
}

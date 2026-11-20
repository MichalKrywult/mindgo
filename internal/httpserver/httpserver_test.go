package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

type mockProvider struct {
	entries []domain.MoodEntry
}

func (m mockProvider) GetEntries() []domain.MoodEntry {
	return m.entries
}

func TestGreetUser(t *testing.T) {

	server := NewHTTPServer(mockProvider{})

	req, err := http.NewRequest(http.MethodGet, "/", nil) // no body
	if err != nil {
		t.Fatal("creating request failed")
	}

	rec := httptest.NewRecorder() // pretends to be a receiver
	server.greetUser(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("expected %v, got %v", "text/plain; charset=utf-8", rec.Header().Get("Content-Type"))
	}

	if rec.Body.String() != "Hello from Mindgo!\n" {
		t.Errorf("expected %v, got %v", "Hello from Mindgo!\n", rec.Body.String())
	}
}

func TestGetEntries(t *testing.T) {

	tests := []struct {
		name            string
		expectedEntries []domain.MoodEntry
	}{
		{
			name:            "no entries",
			expectedEntries: []domain.MoodEntry{},
		},
		{
			name:            "one entry",
			expectedEntries: []domain.MoodEntry{{ID:1, Mood: 1}},
		},
		{
			name:            "multiple entries",
			expectedEntries: []domain.MoodEntry{{ID:1, Mood: 1},{ID:2, Mood: 2},{ID:3, Mood: 3},{ID:4, Mood: 4}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := NewHTTPServer(mockProvider{entries: tt.expectedEntries})

			req, err := http.NewRequest(http.MethodGet, "/api/entries", nil)
			if err != nil {
				t.Fatal("creating request failed")
			}

			rec := httptest.NewRecorder() // pretends to be a receiver
			server.getEntries(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("expected status %v, got %v", http.StatusOK, rec.Code)
			}

			if rec.Header().Get("Content-Type") != "application/json; charset=utf-8" {
				t.Errorf("expected %v, got %v", "application/json; charset=utf-8", rec.Header().Get("Content-Type"))
			}

			var entries []domain.MoodEntry
			err = json.NewDecoder(rec.Body).Decode(&entries)
			if err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if !reflect.DeepEqual(entries, tt.expectedEntries) {
				t.Errorf("expected %v, got %v", tt.expectedEntries, entries)
			}

		})
	}
}

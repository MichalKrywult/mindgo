package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

type mockProvider struct{}

func (m mockProvider) GetEntries() []domain.MoodEntry {
	return []domain.MoodEntry{}
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
	}

	server := NewHTTPServer(mockProvider{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

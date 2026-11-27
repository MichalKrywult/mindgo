package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

type mockProvider struct {
	entries []domain.MoodEntry
}

func (m mockProvider) GetEntries() []domain.MoodEntry {
	return m.entries
}

func (m *mockProvider) AddEntry(addEntry domain.MoodEntry) error {
	m.entries = append(m.entries, addEntry)
	return nil
}

func TestRouting(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		url          string
		body         string
		expectedCode int
	}{
		{
			name:         "unknown path",
			url:          "/unknown",
			method:       http.MethodGet,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "default path",
			url:          "/",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
		},
		{
			name:         "entries path, correct Get method",
			url:          "/api/entries",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
		},
		{
			name:         "entries path, wrong method",
			url:          "/api/entries",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "create entry path, Post method",
			url:          "/api/entries",
			method:       http.MethodPost,
			expectedCode: http.StatusCreated,
			body:         `{"mood":5}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewHTTPServer(&mockProvider{})

			req, err := http.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal("creating request failed")
			}

			rec := httptest.NewRecorder()       // pretends to be a receiver
			server.routes().ServeHTTP(rec, req) // let the mux choose correct path

			if rec.Code != tt.expectedCode {
				t.Errorf("expected status %v, got %v", tt.expectedCode, rec.Code)
			}
		})
	}
}

func TestGreetUser(t *testing.T) {

	server := NewHTTPServer(&mockProvider{})

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
			expectedEntries: []domain.MoodEntry{{ID: 1, Mood: 1}},
		},
		{
			name:            "multiple entries",
			expectedEntries: []domain.MoodEntry{{ID: 1, Mood: 1}, {ID: 2, Mood: 2}, {ID: 3, Mood: 3}, {ID: 4, Mood: 4}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			server := NewHTTPServer(&mockProvider{entries: tt.expectedEntries})

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
func TestCreateEntry(t *testing.T) {

	date, err := time.Parse("2006-01-02 15:04", "2026-11-27 14:30")
	if err != nil {
		t.Fatal("error preparing date")
	}

	tests := []struct {
		name            string
		body            string
		expectedCode    int
		expectedEntries []domain.MoodEntry
	}{
		{
			name:            "successful entry creation",
			body:            `{"mood": 5, "date": "2026-11-27 14:30"}`,
			expectedEntries: []domain.MoodEntry{{Mood: 5, Date: date, Tags: []string{}}},
			// adding tags assertion is necessary -> reflect.DeepEqual
			// distinguishes nil and empty Tags list
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid request body",
			body:         `invalid json`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProvider{}
			server := NewHTTPServer(mock)

			req, err := http.NewRequest(http.MethodPost, "/api/entries", strings.NewReader(tt.body))
			if err != nil {
				t.Fatal("creating request failed")
			}

			rec := httptest.NewRecorder()
			server.createEntry(rec, req)

			if rec.Code != tt.expectedCode {
				t.Errorf("expected status %v, got %v", tt.expectedCode, rec.Code)
			}

			if len(mock.entries) != len(tt.expectedEntries) {
				t.Fatalf("expected %d entries, got %d", len(tt.expectedEntries), len(mock.entries))
			}

			if !reflect.DeepEqual(mock.entries, tt.expectedEntries) {
				t.Errorf("expected entries %v, got %v", tt.expectedEntries, mock.entries)
			}
			t.Logf("expected: %#v", tt.expectedEntries)
			t.Logf("got:      %#v", mock.entries)
		})
	}
}

package httpserver

import (
	"net/http"
	"net/http/httptest"
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

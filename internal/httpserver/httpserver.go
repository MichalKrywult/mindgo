package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

type HTTPServer struct {
	provider entryProvider
}

type entryProvider interface {
	GetEntries() []domain.MoodEntry
	AddEntry(domain.MoodEntry) error
}

func NewHTTPServer(provider entryProvider) *HTTPServer {
	return &HTTPServer{provider: provider}
}

func (server *HTTPServer) Run() error {

	fmt.Println("Server works!")
	return http.ListenAndServe(":8080", server.routes())
}

func (server *HTTPServer) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", server.greetUser)
	mux.HandleFunc("GET /api/entries", server.getEntries)
	mux.HandleFunc("POST /api/entries", server.createEntry)

	return mux
}

type createEntryRequest struct {
	Mood *int     `json:"mood"` // pointer so it can be compared to nil
	Note string   `json:"note"`
	Tags []string `json:"tags"`
	Date string   `json:"date"`
}

func (server *HTTPServer) greetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Hello from Mindgo!")
}

func (server *HTTPServer) createEntry(w http.ResponseWriter, r *http.Request) {
	var request createEntryRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "placeholder", http.StatusBadRequest)
		return
	}

	if request.Mood == nil {
		http.Error(w, "mood value is obligatory", http.StatusBadRequest)
		return
	}

	var date time.Time

	if request.Date == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02 15:04", request.Date)
		if err != nil {
			http.Error(w, "Invalid date", http.StatusBadRequest)
			return
		}
	}

	entry, err := domain.NewMoodEntry(*request.Mood, date, request.Note, request.Tags)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}

	err = server.provider.AddEntry(entry)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (server *HTTPServer) getEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	entries := server.provider.GetEntries()
	err := json.NewEncoder(w).Encode(entries) // encode to json, NewEncoder implements RepsonseWriter
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

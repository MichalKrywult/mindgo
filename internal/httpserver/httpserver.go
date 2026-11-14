package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MichalKrywult/mindgo/internal/domain"
)

type HTTPServer struct {
	provider entryProvider
}

type entryProvider interface {
	GetEntries() []domain.MoodEntry
}

func NewHTTPServer(provider entryProvider) *HTTPServer {
	return &HTTPServer{provider: provider}
}

func (server *HTTPServer) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.greetUser)

	mux.HandleFunc("/api/entries", server.getEntries)

	fmt.Println("Server works!")
	return http.ListenAndServe(":8080", mux)
}

func (server *HTTPServer) greetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Hello from Mindgo!")
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

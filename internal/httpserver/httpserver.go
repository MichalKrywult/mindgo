package httpserver

import (
	"fmt"
	"net/http"

	"github.com/MichalKrywult/mindgo/internal/tracker"
)

type HTTPServer struct {
	moodTracker *tracker.MoodTracker
}

func NewHTTPServer(moodTracker *tracker.MoodTracker) *HTTPServer {
	return &HTTPServer{moodTracker: moodTracker}
}

func (server *HTTPServer) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", greetUser)

	fmt.Println("Server works!")
	return http.ListenAndServe(":8080", mux)
}

func greetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from Mindgo!")
}

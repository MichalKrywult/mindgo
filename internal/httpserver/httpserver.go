package httpserver

import (
	"fmt"
	"net/http"
)

func greetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from Mindgo!")
}

func ServeHTTP() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", greetUser)

	return http.ListenAndServe(":8080", mux)
}

package httpserver

import (
	"fmt"
	"net/http"
)

func ServeHTTP() error {

	mux := http.NewServeMux()

	mux.HandleFunc("/", greetUser)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return err
	}

	return nil
}

func greetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from Mindgo!")
}

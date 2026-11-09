package httpserver

import (
	"fmt"
	"net/http"
)

func ServeHTTP() error {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Mindgo!")
	})

	err := http.ListenAndServe(":8080", nil)
	if err!=nil{
		return err
	}

	return nil
}

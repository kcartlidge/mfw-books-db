package main

import (
	"fmt"
	"net/http"
)

// HomeHandler is the handler for the home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body><h1>MFW Books DB</h1></body></html>")
}

// NotFoundHandler is the handler for 404 errors
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<html><body><h1>404 Not Found</h1></body></html>")
}

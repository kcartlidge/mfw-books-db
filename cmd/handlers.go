package main

import (
	"fmt"
	"net/http"
)

// HomeHandler is the handler for the home page
func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new template manager
	templates, err := NewTemplates()
	if err != nil {
		http.Error(w, "Error loading templates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Load the books from the JSON file
	books := LoadFile(s.Filename)

	// Create the template data
	data := TemplateData{
		Title:    "All Books",
		Filename: s.Filename,
		Content:  books,
	}

	// Render the template
	if err := templates.Render(w, "home.go.html", data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// NotFoundHandler is the handler for 404 errors
func (s *Server) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<html><body><h1>404 Not Found</h1></body></html>")
}

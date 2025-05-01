package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
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

	// Get the sort field and direction from cookies
	sortField, err := s.CookieHandler.GetCookie(r, "mfw-sort-details")
	if err == nil {
		sortDirection, _ := s.CookieHandler.GetCookie(r, "mfw-sort-direction")
		descending := sortDirection == "desc"

		// Apply the appropriate sort based on the field
		switch strings.ToLower(sortField) {
		case "isbn":
			SortBooksByISBN(books, descending)
		case "status":
			SortBooksByStatus(books, descending)
		case "title":
			SortBooksByTitle(books, descending)
		case "author":
			SortBooksByAuthor(books, descending)
		case "series":
			SortBooksBySeries(books, descending)
		case "rating":
			SortBooksByRating(books, descending)
		case "genre":
			SortBooksByGenre(books, descending)
		}
	}

	// Create the template data
	data := TemplateData{
		Title:     "All Books",
		Filename:  s.Filename,
		Content:   books,
		SortField: sortField,
	}

	// Render the template
	if err := templates.Render(w, "home.go.html", data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// SortHandler handles sorting requests
func (s *Server) SortHandler(w http.ResponseWriter, r *http.Request) {
	// Get the sort field from the URL
	vars := mux.Vars(r)
	newSortField := vars["field"]

	// Get the current sort field and direction from cookies
	currentSortField, _ := s.CookieHandler.GetCookie(r, "mfw-sort-details")
	currentSortDirection, _ := s.CookieHandler.GetCookie(r, "mfw-sort-direction")

	// Determine the new sort direction
	var newSortDirection string
	if newSortField == currentSortField {
		// Invert the current direction
		if currentSortDirection == "asc" {
			newSortDirection = "desc"
		} else {
			newSortDirection = "asc"
		}
	} else {
		// Use default direction for the new field
		if newSortField == "rating" {
			newSortDirection = "desc"
		} else {
			newSortDirection = "asc"
		}
	}

	// Set both cookies
	err := s.CookieHandler.SetCookie(w, "mfw-sort-details", newSortField, 86400) // 24 hours
	if err != nil {
		http.Error(w, "Error setting sort field cookie: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.CookieHandler.SetCookie(w, "mfw-sort-direction", newSortDirection, 86400) // 24 hours
	if err != nil {
		http.Error(w, "Error setting sort direction cookie: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// NotFoundHandler is the handler for 404 errors
func (s *Server) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<html><body><h1>404 Not Found</h1></body></html>")
}

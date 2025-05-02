package main

import (
	"fmt"
	"net/http"
	"strconv"
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
	if err := templates.Render(w, "home", data); err != nil {
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

// EditHandler handles the book edit page
func (s *Server) EditHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ISBN from the URL
	vars := mux.Vars(r)
	isbn := vars["isbn"]

	// Load the books from the JSON file
	books := LoadFile(s.Filename)

	// Find the book with the matching ISBN
	var book *Book
	for i := range books {
		if books[i].ISBN == isbn {
			book = &books[i]
			break
		}
	}

	// Create a map to store unique series
	seriesMap := make(map[string]bool)
	for _, b := range books {
		if b.Series != "" {
			seriesMap[b.Series] = true
		}
	}

	// Convert map to slice and sort
	series := make([]string, 0, len(seriesMap))
	for s := range seriesMap {
		series = append(series, s)
	}
	SortStrings(series, false) // Use existing sort function

	// Create a new template manager
	templates, err := NewTemplates()
	if err != nil {
		http.Error(w, "Error loading templates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the template data
	data := TemplateData{
		Title:    "Edit Book",
		Filename: s.Filename,
		Content:  book,
		Series:   series,
	}

	// Render the template
	if err := templates.Render(w, "edit", data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// SaveHandler handles saving book edits
func (s *Server) SaveHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ISBN from the URL
	vars := mux.Vars(r)
	isbn := vars["isbn"]

	// Parse the form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	title := r.FormValue("title")
	authorSort := r.FormValue("authorSort")
	if title == "" || authorSort == "" {
		http.Error(w, "Title and Author Sort are required", http.StatusBadRequest)
		return
	}

	// Load the books from the JSON file
	books := LoadFile(s.Filename)

	// Find the book with the matching ISBN
	var bookIndex int = -1
	for i := range books {
		if books[i].ISBN == isbn {
			bookIndex = i
			break
		}
	}

	if bookIndex == -1 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Verify the book's ID matches the hidden ID field
	hiddenID := r.FormValue("id")
	if books[bookIndex].ID != hiddenID {
		http.Error(w, "Book ID mismatch", http.StatusBadRequest)
		return
	}

	// Update only the allowed fields
	books[bookIndex].Title = title
	books[bookIndex].AuthorSort = splitAndTrim(authorSort)
	books[bookIndex].Genre = splitAndTrim(r.FormValue("genres"))
	books[bookIndex].Series = r.FormValue("series")
	books[bookIndex].Sequence = r.FormValue("sequence")
	books[bookIndex].Status = r.FormValue("status")
	if len(books[bookIndex].Status) > 0 {
		books[bookIndex].StatusIcon = string(books[bookIndex].Status[0]) // First character of status
	}
	books[bookIndex].Notes = r.FormValue("notes")

	// Parse rating
	ratingStr := r.FormValue("rating")
	if ratingStr != "" {
		rating, err := strconv.Atoi(ratingStr)
		if err != nil || rating < 0 || rating > 5 {
			http.Error(w, "Rating must be a whole number between 0 and 5", http.StatusBadRequest)
			return
		}
		books[bookIndex].Rating = rating
	} else {
		books[bookIndex].Rating = 0
	}

	// Save the updated books
	if err := SaveFile(s.Filename, books); err != nil {
		http.Error(w, "Error saving file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to the home page
	http.Redirect(w, r, "/#b_"+isbn, http.StatusSeeOther)
}

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"

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

	// Default title and filter
	title := "All Books"
	var filter BookFilter

	// Get the filter from cookie
	filterName, err := s.CookieHandler.GetCookie(r, "mfw-filter")
	if err == nil {
		// Apply the appropriate filter
		switch strings.ToLower(filterName) {
		case "all":
			filter = GetPopulatedAllBooksFilter(books)
		case "reading":
			filter = GetPopulatedReadingFilter(books)
		case "next":
			filter = GetPopulatedNextFilter(books)
		case "done":
			filter = GetPopulatedDoneFilter(books)
		case "other":
			filter = GetPopulatedOtherFilter(books)
		}
		books = filter.Books
		title = filter.Name
	}

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
		Title:     title,
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

// FilterHandler handles filter requests
func (s *Server) FilterHandler(w http.ResponseWriter, r *http.Request) {
	// Get the filter name from the URL
	vars := mux.Vars(r)
	filterName := vars["filter"]

	// Set the filter cookie
	err := s.CookieHandler.SetCookie(w, "mfw-filter", filterName, 86400) // 24 hours
	if err != nil {
		http.Error(w, "Error setting filter cookie: "+err.Error(), http.StatusInternalServerError)
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

	// Create a map to store unique genres
	genreMap := make(map[string]bool)
	for _, b := range books {
		if len(b.Genre) > 0 && b.Genre[0] != "" {
			genreMap[b.Genre[0]] = true
		}
		if len(b.Genre) > 1 && b.Genre[1] != "" {
			genreMap[b.Genre[1]] = true
		}
	}

	// Convert map to slice and sort
	genres := make([]string, 0, len(genreMap))
	for g := range genreMap {
		genres = append(genres, g)
	}
	SortStrings(genres, false)

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
		Genres:   genres,
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
	books[bookIndex].Title = strings.TrimSpace(title)
	books[bookIndex].AuthorSort = splitAndTrim(authorSort)
	books[bookIndex].Genre[0] = cleanGenre(r.FormValue("genre1"))
	books[bookIndex].Genre[1] = cleanGenre(r.FormValue("genre2"))
	books[bookIndex].Series = strings.TrimSpace(r.FormValue("series"))
	books[bookIndex].Sequence = strings.TrimSpace(r.FormValue("sequence"))
	books[bookIndex].Status = strings.TrimSpace(r.FormValue("status"))
	books[bookIndex].Notes = strings.TrimSpace(r.FormValue("notes"))
	if len(books[bookIndex].Status) > 0 {
		books[bookIndex].StatusIcon = string(books[bookIndex].Status[0]) // First character of status
	}

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

// capitalizeWords capitalizes the first letter of each word in a string
func capitalizeWords(s string) string {
	// List of words to preserve as-is
	preserve := map[string]bool{
		"SF": true,
	}

	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			if preserve[strings.ToUpper(word)] {
				words[i] = word
			} else {
				words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
			}
		}
	}
	return strings.Join(words, " ")
}

// cleanGenre removes non-alphanumeric characters from the start and end of a genre string
func cleanGenre(genre string) string {
	// First trim whitespace
	genre = strings.TrimSpace(genre)
	// Remove trailing non-alphanumeric characters
	for len(genre) > 0 && !unicode.IsLetter(rune(genre[len(genre)-1])) && !unicode.IsNumber(rune(genre[len(genre)-1])) {
		genre = genre[:len(genre)-1]
	}
	// Remove leading non-alphanumeric characters
	for len(genre) > 0 && !unicode.IsLetter(rune(genre[0])) && !unicode.IsNumber(rune(genre[0])) {
		genre = genre[1:]
	}
	return capitalizeWords(genre)
}

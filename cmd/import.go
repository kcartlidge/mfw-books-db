package main

import (
	"fmt"
	"strings"
	"time"
)

// ProcessISBNs takes a slice of ISBNs and queries Google Books for each one
func ProcessISBNs(isbns []string, existingBooks map[string]Book) {
	// Create a grid to track all books
	grid := NewGrid([]string{"ISBN", "New?", "Title", "Author", "Error"})

	// Track counts
	var newCount, matchedCount, errorCount int

	for _, isbn := range isbns {
		// Check if we already have this book
		if existingBook, exists := existingBooks[isbn]; exists {
			grid.AddRow(
				isbn,
				"No",
				existingBook.Title,
				existingBook.Authors,
				"",
			)
			matchedCount++
			continue
		}

		// Get the book from Google Books
		gb, err := GetBookByISBN(isbn)
		if err != nil {
			grid.AddRow(
				isbn,
				"Error",
				"",
				"",
				err.Error(),
			)
			errorCount++
			continue
		}

		// Map to our Book model
		book := mapGoogleBook(isbn, gb)

		// Add to the grid
		grid.AddRow(
			isbn,
			"Yes",
			book.Title,
			book.Authors,
			"",
		)
		newCount++
	}

	// Print the grid
	fmt.Println(grid)
	fmt.Println()

	// Print summary
	fmt.Printf("Starting with %d books in the database.\n", len(existingBooks))
	fmt.Printf("Processed %d ISBNs: %d new books added, %d matched existing books, %d errors.\n",
		len(isbns), newCount, matchedCount, errorCount)
	fmt.Printf("Final total: %d books.\n", len(existingBooks)+newCount)
	fmt.Println()
}

// mapGoogleBook converts a GoogleBook to our Book model
func mapGoogleBook(isbn string, gb *GoogleBook) Book {
	// Join authors with commas for the Authors field
	authors := strings.Join(gb.Authors, ", ")

	// Create author sort string (last name, first name)
	authorSort := ""
	if len(gb.Authors) > 0 {
		parts := strings.Split(gb.Authors[0], " ")
		if len(parts) > 1 {
			authorSort = fmt.Sprintf("%s, %s", parts[len(parts)-1], strings.Join(parts[:len(parts)-1], " "))
		} else {
			authorSort = gb.Authors[0]
		}
	}

	return Book{
		ISBN:        isbn,
		Title:       gb.Title,
		Authors:     authors,
		AuthorSort:  authorSort,
		Genre:       gb.Categories,
		Link:        gb.Link,
		IsException: false,
		Status:      "U - Unread",
		StatusIcon:  "U",
		ModifiedUtc: time.Now().UTC().Format(time.RFC3339),
		// Additional fields from Google Books
		PublishedDate: gb.PublishedDate,
		Publisher:     gb.Publisher,
		PageCount:     gb.PageCount,
		Language:      gb.Language,
		Description:   gb.Description,
	}
}

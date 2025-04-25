package main

import (
	"fmt"
	"strings"
	"time"
)

// ProcessISBNs takes a slice of ISBNs and queries Google Books for each one
func ProcessISBNs(isbns []string, books []Book) []Book {
	// Create a grid to track all books
	grid := NewGrid([]string{"ISBN", "New?", "Title", "Author", "Error"})
	grid.SetShowNumbers(true)

	// Track counts
	var newCount, matchedCount, errorCount int
	originalCount := len(books)

	for _, isbn := range isbns {
		// Check if we already have this book
		found := false
		for _, book := range books {
			if book.ISBN == isbn {
				grid.AddRow(
					isbn,
					"-",
					book.Title,
					book.GetAuthorDisplay(),
					book.ExceptionReason,
				)
				matchedCount++
				found = true
				break
			}
		}
		if found {
			continue
		}

		// Get the book from Google Books
		gb, err := GetBookByISBN(isbn)
		if err != nil {
			// Create a book with just the ISBN and error information
			book := Book{
				ISBN:            isbn,
				IsException:     true,
				ExceptionReason: err.Error(),
				ModifiedUtc:     time.Now().UTC().Format(time.RFC3339),
			}
			books = append(books, book)
			grid.AddRow(
				isbn,
				"Error",
				"",
				"",
				err.Error(),
			)
			errorCount++ // Only count new errors
			continue
		}

		// Map to our Book model
		book := mapGoogleBook(isbn, gb)

		// Add to the grid and the books slice
		grid.AddRow(
			isbn,
			"Yes",
			book.Title,
			book.GetAuthorDisplay(),
			"",
		)
		books = append(books, book)
		newCount++
	}

	// Print the grid
	fmt.Println(grid)
	fmt.Println()

	// Print summary
	fmt.Printf("Started with %d books in the database.\n", originalCount)
	fmt.Printf("%d added, %d matched, and %d new errors.\n",
		newCount, matchedCount, errorCount)
	fmt.Printf("Ended with %d books in the database.\n", len(books))
	fmt.Println()

	return books
}

// mapGoogleBook converts a GoogleBook to our Book model
func mapGoogleBook(isbn string, gb *GoogleBook) Book {
	// Create author sort strings (last name, first name) for each author
	authorSorts := make([]string, len(gb.Authors))
	for i, author := range gb.Authors {
		parts := strings.Split(author, " ")
		if len(parts) > 1 {
			authorSorts[i] = fmt.Sprintf("%s, %s", parts[len(parts)-1], strings.Join(parts[:len(parts)-1], " "))
		} else {
			authorSorts[i] = author
		}
	}

	return Book{
		ISBN:          isbn,
		Title:         gb.Title,
		Authors:       gb.Authors,
		AuthorSort:    authorSorts,
		Genre:         gb.Categories,
		Link:          gb.Link,
		IsException:   false,
		Status:        "Unread",
		StatusIcon:    "U",
		ModifiedUtc:   time.Now().UTC().Format(time.RFC3339),
		PublishedDate: gb.PublishedDate,
		Publisher:     gb.Publisher,
		PageCount:     gb.PageCount,
		Language:      gb.Language,
		Description:   gb.Description,
	}
}

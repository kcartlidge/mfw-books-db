package main

import (
	"fmt"
	"strings"
	"time"
)

// ProcessISBNs takes a slice of ISBNs and queries Google Books for each one
func ProcessISBNs(isbns []string, books []Book, errorsCleared bool, singleHit bool) []Book {
	// Create a grid to track all books
	grid := NewGrid([]string{"ISBN", "NEW?", "TITLE", "AUTHORS", "ERROR"})
	grid.SetShowNumbers(true)

	// Track counts
	var newCount, matchedCount, errorCount int
	originalCount := len(books)

	fmt.Print("Processing:")
	for i, isbn := range isbns {
		if i%5 == 0 {
			fmt.Printf(" %d", i)
		}

		// Look up the book
		book, found, err := lookupBook(isbn, books, singleHit)
		if found {
			grid.AddRow(
				isbn,
				"-",
				book.Title,
				book.GetAuthorSortDisplay(),
				book.ExceptionReason,
			)
			matchedCount++
			continue
		}

		if err != nil {
			grid.AddRow(
				isbn,
				"Error",
				"",
				"",
				err.Error(),
			)
			errorCount++ // Only count new errors
		} else {
			grid.AddRow(
				isbn,
				"Yes",
				book.Title,
				book.GetAuthorSortDisplay(),
				"",
			)
			newCount++
		}
		books = append(books, book)
	}
	if len(isbns)%5 != 0 {
		fmt.Printf(" %d", len(isbns))
	}
	fmt.Println()
	fmt.Println()

	// Print the grid
	fmt.Println(grid)
	fmt.Println()

	// Print summary
	if errorsCleared {
		fmt.Printf("Started with %d books in the database (after clearing any errors).\n", originalCount)
	} else {
		fmt.Printf("Started with %d books in the database.\n", originalCount)
	}
	fmt.Printf("%d added, %d matched, and %d new errors.\n",
		newCount, matchedCount, errorCount)
	fmt.Printf("Ended with %d books in the database.\n", len(books))
	fmt.Println()

	return books
}

// lookupBook checks if a book exists and if not, looks it up in Google Books
func lookupBook(isbn string, books []Book, singleHit bool) (Book, bool, error) {
	// Check if we already have this book
	for _, book := range books {
		if book.ISBN == isbn {
			return book, true, nil
		}
	}

	// Get the book from Google Books
	gb, err := GetBookByISBN(isbn, singleHit)
	if err != nil {
		// Create a book with just the ISBN and error information
		book := Book{
			ID:              "",
			ISBN:            isbn,
			IsException:     true,
			ExceptionReason: err.Error(),
			ModifiedUtc:     time.Now().UTC().Format(time.RFC3339),
		}
		return book, false, err
	}

	// Map to our Book model
	book := mapGoogleBook(isbn, gb)
	return book, false, nil
}

// mapGoogleBook converts a GoogleBook to our Book model
func mapGoogleBook(isbn string, gb *GoogleBook) Book {
	return Book{
		ID:            gb.ID,
		ISBN:          isbn,
		Title:         fixTitle(gb.Title),
		Authors:       gb.Authors,
		AuthorSort:    fixAuthorSorts(gb.Authors),
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

// fixTitle moves "The " from the start to the end of the title
func fixTitle(title string) string {
	if strings.HasPrefix(title, "The ") {
		return fmt.Sprintf("%s, the", title[4:])
	}
	return title
}

// fixAuthorSorts creates author sort strings (last name, first name) for each author
// This is not internationalised, so it only works for English
// It also handles initials (with or without periods)
func fixAuthorSorts(authors []string) []string {
	sorts := []string{}

	// Process each author
	for _, author := range authors {
		// Get all the author segments
		segments := []string{}
		for _, p := range strings.Split(author, " ") {
			if len(p) > 0 {
				if strings.Contains(p, ".") {
					// Initials should be capitalised and the period removed
					segments = append(segments, strings.ToUpper(strings.TrimSpace(strings.ReplaceAll(p, ".", ""))))
				} else {
					segments = append(segments, strings.TrimSpace(p))
				}
			}
		}

		// Generate the resulting name
		if len(segments) == 1 {
			// Single segment, so just use the name
			sorts = append(sorts, segments[0])
		} else {
			// Start with the last name
			name := segments[len(segments)-1] + ","

			// Multiple segments, so we need to figure out the initials
			last := ""
			for _, segment := range segments[:len(segments)-1] {
				if len(last) != 1 {
					name += " "
				}
				// Capitalise single letter initials
				if len(segment) == 1 {
					segment = strings.ToUpper(segment)
				}
				name += segment
				last = segment
			}

			// Add the name to the list
			sorts = append(sorts, name)
		}
	}

	return sorts
}

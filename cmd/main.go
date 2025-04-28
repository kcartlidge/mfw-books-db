package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	fmt.Println()
	fmt.Println()
	fmt.Println("MFW BOOKS DATABASE v1.0.0")

	// Parse command line arguments
	parser := NewArgsParser()
	parser.AddArgument("file", "JSON file containing book data", "", true)
	parser.AddArgument("isbns", "Text file containing ISBNs to process", "", false)
	parser.AddFlag("clear-errors", "Removes errored ISBNs so they retry")
	parser.ShowUsage()
	parser.Parse(os.Args[1:])

	// Check for errors
	if parser.HasErrors() {
		parser.PrintErrors()
		check(errors.New("errors found when parsing command line arguments"))
	}

	// Get the arguments
	parser.ShowProvided()
	jsonFile := parser.GetArgument("file")
	clearErrors := parser.GetFlag("clear-errors")

	// Load the books from the JSON file
	fmt.Println()
	fmt.Println()
	fmt.Println("Loading books from", jsonFile)
	books := LoadFile(jsonFile)
	fmt.Printf("Found %d book(s) in the database\n", len(books))
	fmt.Println()

	// Clear errors if requested
	if clearErrors {
		fmt.Println("Clearing errored ISBNs so they are retried")
		removed, err := ClearErroredBooks(jsonFile)
		if err != nil {
			fmt.Println()
			fmt.Println("ERROR clearing errored ISBNs")
			check(err)
		}
		if removed > 0 {
			fmt.Printf("Removed %d errored ISBNs\n", removed)
		}
		// Reload the books after clearing errors
		books = LoadFile(jsonFile)
		fmt.Printf("There are now %d book(s) in the database\n", len(books))
		fmt.Println()
	}

	// Load the ISBNs from the text file
	if parser.HasArgument("isbns") {
		isbnsFile := parser.GetArgument("isbns")
		fmt.Println("Loading ISBNs from", isbnsFile)
		isbns := LoadISBNs(isbnsFile)
		fmt.Printf("Found %d ISBN(s) to consider for processing\n", len(isbns))
		fmt.Println("Only new ISBNs will be processed")
		fmt.Println()

		// Process the ISBNs
		books = ProcessISBNs(isbns, books, clearErrors)

		// Save the updated books
		// Only save if we have more books than we started with
		if len(books) > len(LoadFile(jsonFile)) {
			if err := SaveFile(jsonFile, books); err != nil {
				fmt.Println()
				fmt.Println("ERROR saving file")
				check(err)
			}
			fmt.Println("Saved books to", jsonFile)
			fmt.Println()
		}
	}

	fmt.Println("Done.")
	fmt.Println()
	fmt.Println()
}

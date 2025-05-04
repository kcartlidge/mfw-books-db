package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	fmt.Println()
	fmt.Println()
	fmt.Println("MFW BOOKS DATABASE v1.0.0")

	// Parse command line arguments
	parser := NewArgsParser()
	parser.AddArgument("file", "JSON file containing book data", "", true)
	parser.AddArgument("isbns", "Text file containing ISBNs to process", "", false)
	parser.AddArgument("serve", "Local web server port for viewing the database", "", false)
	parser.AddFlag("clear-errors", "Removes errored ISBNs so they retry")
	parser.AddFlag("single-hit", "Only call the API once per ISBN (result quality varies)")
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
	singleHit := parser.GetFlag("single-hit")

	// Load the books from the JSON file
	fmt.Println()
	fmt.Println()
	fmt.Println("Loading books from", jsonFile)
	books := LoadFile(jsonFile)
	fmt.Printf("Found %d book(s) in the database\n", len(books))
	fmt.Println()

	// Create backup at launch
	if err := createBackupFile(jsonFile); err != nil {
		fmt.Printf("Error creating launch backup: %v\n", err)
	}

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
		if singleHit {
			fmt.Println("Single hit mode is enabled (only call the API once per ISBN)")
			fmt.Println("The initial call is by ISBN and gets book data including an ID")
			fmt.Println("The fetched genre, publisher, and page count are not always accurate")
			fmt.Println("The optional second call is by ID and often gets more details")
			fmt.Println()
		}

		// Process the ISBNs
		books = ProcessISBNs(isbns, books, clearErrors, singleHit)

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

	// Start the server
	if parser.HasArgument("serve") {
		port := parser.GetArgument("serve")
		portInt, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println("ERROR converting port to int")
			check(err)
		}

		// Get the absolute path of the books file
		absPath, err := filepath.Abs(jsonFile)
		if err != nil {
			fmt.Println("ERROR getting absolute path of books file")
			check(err)
		}

		server, err := NewServer(portInt, absPath)
		if err != nil {
			fmt.Println("ERROR creating server")
			check(err)
		}
		server.Start()
	}

	// Create backup at exit
	if err := createBackupFile(jsonFile); err != nil {
		fmt.Printf("Error creating exit backup: %v\n", err)
	}

	fmt.Println("Done.")
	fmt.Println()
	fmt.Println()
}

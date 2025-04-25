package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println()
	fmt.Println()
	fmt.Println("MFW BOOKS DATABASE v1.0.0")
	fmt.Println()

	// Check if a filename was provided as an argument
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename as an argument")
		fmt.Println()
		os.Exit(1)
	}

	// Load the books from the JSON file
	jsonFile := os.Args[1]
	fmt.Println("Loading books from", jsonFile)
	books := LoadFile(jsonFile)
	fmt.Printf("Found %d book(s) in the database\n", len(books))
	fmt.Println()

	// Load the ISBNs from the text file
	isbnsFile := filepath.Join(filepath.Dir(jsonFile), "isbns.txt")
	fmt.Println("Loading ISBNs from", isbnsFile)
	isbns := LoadISBNs(isbnsFile)
	fmt.Printf("Found %d ISBN(s) to consider for processing\n", len(isbns))
	fmt.Println("Only new ISBNs will be processed")
	fmt.Println()

	// Process the ISBNs
	books = ProcessISBNs(isbns, books)

	// Save the updated books
	// Only save if we have more books than we started with
	if len(books) > len(LoadFile(jsonFile)) {
		if err := SaveFile(jsonFile, books); err != nil {
			fmt.Println()
			fmt.Println("ERROR saving file")
			fmt.Println(err.Error())
			fmt.Println()
			os.Exit(1)
		}
		fmt.Println("Saved books to", jsonFile)
		fmt.Println()
	}
	fmt.Println("Done.")
	fmt.Println()
	fmt.Println()
}

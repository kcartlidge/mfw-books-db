package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println()
	fmt.Println("MFW BOOKS DATABASE v0.0.1")
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
	fmt.Printf("Found %d ISBN(s) to process\n", len(isbns))
	fmt.Println()
}

// check is a helper function to check for errors and exit the program if an error occurs
func check(err error) {
	if err != nil {
		fmt.Println()
		fmt.Println("ERROR")
		fmt.Println(err.Error())
		fmt.Println()
		os.Exit(1)
	}
}

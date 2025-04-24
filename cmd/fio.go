package main

import (
	"encoding/json"
	"os"
	"strings"
)

// LoadFile attempts to load books from a file, returning an empty slice if the file doesn't exist
func LoadFile(filename string) []Book {
	exists, f, err := CheckFileExists(filename)
	check(err)
	if !exists {
		return []Book{}
	}
	defer f.Close()

	// Decode the JSON file into a slice of Book structs
	var books []Book
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&books); err != nil {
		return []Book{}
	}
	return books
}

// LoadISBNs reads ISBNs from a text file, one per line
func LoadISBNs(filename string) []string {
	exists, f, err := CheckFileExists(filename)
	check(err)
	if !exists {
		return []string{}
	}
	defer f.Close()

	content, err := os.ReadFile(filename)
	check(err)

	// Replace CRLF with LF, then split
	content = []byte(strings.ReplaceAll(string(content), "\r\n", "\n"))
	lines := strings.Split(string(content), "\n")

	// Trim each line and add it to a slice if it's not empty
	var isbns []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			isbns = append(isbns, line)
		}
	}

	return isbns
}

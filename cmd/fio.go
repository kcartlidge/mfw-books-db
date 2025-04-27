package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// SaveFile saves books to a JSON file and creates a daily backup
func SaveFile(filename string, books []Book) error {
	// Create backup filename with today's date
	backupDir := filepath.Join(filepath.Dir(filename), "backups")
	backupName := time.Now().Format("2006-01-02") + " " + filepath.Base(filename)
	backupPath := filepath.Join(backupDir, backupName)

	// If a book file already exists, create a backup
	if found, _, _ := CheckFileExists(filename); found {
		// Create backup directory if it doesn't exist
		err := os.MkdirAll(backupDir, 0755)
		if err != nil {
			return err
		}

		// Open existing file
		src, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer src.Close()

		// Create backup file
		dst, err := os.Create(backupPath)
		if err != nil {
			return err
		}
		defer dst.Close()
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
		fmt.Printf("Created backup file %s\n", backupPath)
	}

	// Save new file
	mainFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer mainFile.Close()
	encoder := json.NewEncoder(mainFile)
	encoder.SetIndent("", "  ")
	return encoder.Encode(books)
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

// ClearErroredBooks removes books marked as exceptions from the file
func ClearErroredBooks(filename string) (int, error) {
	// Load the current books
	books := LoadFile(filename)
	originalCount := len(books)

	// Remove errored books in-place
	i := 0
	for _, book := range books {
		if !book.IsException {
			books[i] = book
			i++
		}
	}
	books = books[:i]

	// Save if we removed any books
	if len(books) < originalCount {
		if err := SaveFile(filename, books); err != nil {
			return 0, err
		}
		return originalCount - len(books), nil
	}

	return 0, nil
}

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

	content, err := os.ReadFile(filename)
	check(err)

	var books []Book
	if err := json.Unmarshal(content, &books); err != nil {
		check(err)
	}

	// Validate ratings
	for i := range books {
		if books[i].Rating < 0 || books[i].Rating > 5 {
			books[i].Rating = 0
		}
	}

	return books
}

// SaveFile saves books to a JSON file and creates a daily backup
func SaveFile(filename string, books []Book) error {
	// Sort the books
	SortBooksByTitle(books, false)

	// Ensure array fields are never null and copy authorSort to authors if needed
	for i := range books {
		if books[i].Authors == nil {
			books[i].Authors = []string{}
		}
		if books[i].AuthorSort == nil {
			books[i].AuthorSort = []string{}
		}
		if books[i].Genre == nil {
			books[i].Genre = []string{}
		}
		// If authors is empty but authorSort has values, copy them
		if len(books[i].Authors) == 0 && len(books[i].AuthorSort) > 0 {
			books[i].Authors = books[i].AuthorSort
		}
	}

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
	json, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		check(err)
	}
	err = os.WriteFile(filename, json, 0644)
	if err != nil {
		check(err)
	}
	return nil
}

// LoadISBNs reads ISBNs from a text file, one per line
func LoadISBNs(filename string) []string {
	exists, f, err := CheckFileExists(filename)
	check(err)
	if !exists {
		check(fmt.Errorf("file not found: %s", filename))
	}
	defer f.Close()

	// Read the file into a string and split it into lines
	content := readFileNormalisedToLF(filename)
	lines := strings.Split(content, "\n")

	// Trim each line and add it to a slice if it's not empty
	lineNumber := 0
	var isbns []string
	for _, line := range lines {
		lineNumber++
		line = strings.TrimSpace(strings.ReplaceAll(line, " ", ""))
		if line != "" {
			if len(line) < 7 || len(line) > 13 {
				check(fmt.Errorf("invalid-looking ISBN on line %d: %s", lineNumber, line))
			}

			// Commented this out pending further thought as some ISBNs are not numeric
			// in that they include X etc (eg '033026656X')
			// if _, err := strconv.Atoi(line); err != nil {
			// 	check(fmt.Errorf("non-numeric ISBN on line %d: %s", lineNumber, line))
			// }

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

// CheckFileExists checks if a file exists (and returns a handle to it)
func CheckFileExists(filename string) (bool, *os.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, f, nil
}

// readFileNormalisedToLF reads a file and returns it with all CRs removed
func readFileNormalisedToLF(filename string) string {
	// Sanity check
	exists, f, err := CheckFileExists(filename)
	check(err)
	if !exists {
		return ""
	}
	defer f.Close()

	// Read the file into a byte slice
	bytes, err := io.ReadAll(f)
	check(err)
	text := string(bytes)

	// Normalise the file to LF
	// Looks odd but that's because some ISBN scanners (eg Alfa)
	// create text files with "\r" delimiters (not "\r\n" or "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.ReplaceAll(text, "\n\n", "\n")

	// Return the text
	return text
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

	content, err := os.ReadFile(filename)
	check(err)

	var books []Book
	if err := json.Unmarshal(content, &books); err != nil {
		check(err)
	}

	// Validate ratings and ensure Genre array has exactly 2 elements
	for i := range books {
		if books[i].Rating < 0 || books[i].Rating > 5 {
			books[i].Rating = 0
		}
		// Ensure Genre array has exactly 2 entries
		if books[i].Genre == nil {
			books[i].Genre = []string{"", ""}
		}
		books[i].Genre = append(books[i].Genre, "", "")[:2]
	}

	return books
}

// SaveFile saves books to a JSON file
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
		// Create new genre array of size 2 and copy up to 2 entries
		oldGenres := books[i].Genre
		books[i].Genre = make([]string, 2)
		copy(books[i].Genre, oldGenres)

		// If authors is empty but authorSort has values, copy them
		if len(books[i].Authors) == 0 && len(books[i].AuthorSort) > 0 {
			books[i].Authors = books[i].AuthorSort
		}
		// Ensure Genre array has exactly 2 entries
		if len(books[i].Genre) < 2 {
			// Pad with empty strings if needed
			books[i].Genre = append(books[i].Genre, make([]string, 2-len(books[i].Genre))...)
		} else if len(books[i].Genre) > 2 {
			// Truncate to 2 entries
			books[i].Genre = books[i].Genre[:2]
		}

		// Trim string fields
		books[i].Title = strings.TrimSpace(books[i].Title)
		books[i].Series = strings.TrimSpace(books[i].Series)
		books[i].Sequence = strings.TrimSpace(books[i].Sequence)
		books[i].Status = strings.TrimSpace(books[i].Status)
		books[i].Notes = strings.TrimSpace(books[i].Notes)
	}

	// Save file
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

// createBackupFile creates a backup of the source file
func createBackupFile(sourcePath string) error {
	// Check if file exists
	found, f, err := CheckFileExists(sourcePath)
	if err != nil {
		return err
	}
	if !found {
		return nil // No file to backup
	}
	defer f.Close()

	// Get file info for last modified time
	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}

	// Create backup filename with last modified date
	backupDir := filepath.Join(filepath.Dir(sourcePath), "backups")
	backupName := fileInfo.ModTime().Format("2006-01-02") + " " + filepath.Base(sourcePath)
	backupPath := filepath.Join(backupDir, backupName)

	// Create backup directory if it doesn't exist
	err = os.MkdirAll(backupDir, 0755)
	if err != nil {
		return err
	}

	// Open existing file
	src, err := os.Open(sourcePath)
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

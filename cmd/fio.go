package main

import (
	"encoding/json"
)

// LoadFile attempts to load books from a file, returning an empty slice if the file doesn't exist
func LoadFile(filename string) []Book {
	exists, f, err := CheckFileExists(filename)
	check(err)
	if !exists {
		return []Book{}
	}
	defer f.Close()

	var books []Book
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&books); err != nil {
		return []Book{}
	}
	return books
}

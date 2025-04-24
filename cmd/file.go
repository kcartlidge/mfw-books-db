package main

import (
	"os"
)

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

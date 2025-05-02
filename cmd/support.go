package main

import (
	"fmt"
	"os"
	"strings"
)

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

// joinWithAmpersand joins a slice of strings with an ampersand
func joinWithAmpersand(items []string) string {
	return strings.Join(items, " & ")
}

// splitAndTrim splits a string on ampersands and trims each segment
func splitAndTrim(s string) []string {
	if s == "" {
		return []string{}
	}
	segments := strings.Split(s, "&")
	result := make([]string, 0, len(segments))
	for _, segment := range segments {
		segment = strings.TrimSpace(segment)
		if segment != "" {
			result = append(result, segment)
		}
	}
	return result
}

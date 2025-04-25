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

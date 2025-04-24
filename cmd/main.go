package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println()
	fmt.Println("MFW BOOKS DATABASE v0.0.1")
	fmt.Println()

	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename as an argument")
		fmt.Println()
		os.Exit(1)
	}

	filename := os.Args[1]
	fmt.Println("Loading books from", filename)
	books := LoadFile(filename)
	fmt.Printf("Found %d book(s) in the database\n", len(books))
	fmt.Println()
}

func check(err error) {
	if err != nil {
		fmt.Println()
		fmt.Println("ERROR")
		fmt.Println(err.Error())
		fmt.Println()
		os.Exit(1)
	}
}

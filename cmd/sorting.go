package main

import "sort"

// SortBooksByTitleAuthorISBN sorts books by title, then author, then ISBN
func SortBooksByTitleAuthorISBN(books []Book) {
	sort.Slice(books, func(i, j int) bool {
		// First compare titles
		if books[i].Title != books[j].Title {
			return books[i].Title < books[j].Title
		}

		// If titles are equal, compare authors
		authorI := books[i].GetAuthorDisplay()
		authorJ := books[j].GetAuthorDisplay()
		if authorI != authorJ {
			return authorI < authorJ
		}

		// If authors are equal, compare ISBNs
		return books[i].ISBN < books[j].ISBN
	})
}

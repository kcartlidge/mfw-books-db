package main

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
)

// SortBooksByISBN sorts books by ISBN, then author, then series/sequence
func SortBooksByISBN(books []Book, descending bool) {
	sortBooksByFallbackOrder(books)
	slices.SortStableFunc(books, func(a, b Book) int {
		result := strings.Compare(a.ISBN, b.ISBN)
		if descending {
			return -result
		}
		return result
	})
}

// SortBooksByStatus sorts books by status, then author, then series/sequence
func SortBooksByStatus(books []Book, descending bool) {
	sortBooksByFallbackOrder(books)
	slices.SortStableFunc(books, func(a, b Book) int {
		// Get status letters
		aStatus := a.GetStatusLetter()
		bStatus := b.GetStatusLetter()

		// If either book has no status, sort it to the end
		if aStatus == "-" && bStatus == "-" {
			return 0
		} else if aStatus == "-" {
			return 1 // a has no status, sort it after b
		} else if bStatus == "-" {
			return -1 // b has no status, sort it after a
		}

		// Priority order: C and N first, then others alphabetically
		if aStatus == "C" || aStatus == "N" {
			if bStatus == "C" || bStatus == "N" {
				// Both are priority statuses, compare them
				result := strings.Compare(aStatus, bStatus)
				if descending {
					return -result
				}
				return result
			}
			return -1 // a is priority, sort it before b
		}
		if bStatus == "C" || bStatus == "N" {
			return 1 // b is priority, sort it before a
		}

		// Neither is priority, compare alphabetically
		result := strings.Compare(aStatus, bStatus)
		if descending {
			return -result
		}
		return result
	})
}

// SortBooksByTitle sorts books by title, then author, then series/sequence
func SortBooksByTitle(books []Book, descending bool) {
	sortBooksByFallbackOrder(books)
	slices.SortStableFunc(books, func(a, b Book) int {
		result := strings.Compare(a.Title, b.Title)
		if descending {
			return -result
		}
		return result
	})
}

// SortBooksByAuthor sorts books by author, then series/sequence
func SortBooksByAuthor(books []Book, descending bool) {
	sortBooksByFallbackOrder(books)
	slices.SortStableFunc(books, func(a, b Book) int {
		result := strings.Compare(a.GetFirstAuthorSort(), b.GetFirstAuthorSort())
		if descending {
			return -result
		}
		return result
	})
}

// SortBooksBySeries sorts books by series, then author, then sequence
func SortBooksBySeries(books []Book, descending bool) {
	sortBooksByFallbackOrderSupportingDescendingSeries(books, descending)
}

// SortBooksByRating sorts books by rating, then author, then series/sequence
func SortBooksByRating(books []Book, descending bool) {
	sortBooksByFallbackOrder(books)
	slices.SortStableFunc(books, func(a, b Book) int {
		// If either book has no rating, sort it to the end
		if a.Rating == 0 && b.Rating == 0 {
			return 0
		} else if a.Rating == 0 {
			return 1 // a has no rating, sort it after b
		} else if b.Rating == 0 {
			return -1 // b has no rating, sort it after a
		}
		// Both have ratings, sort higher ratings first
		result := cmp.Compare(b.Rating, a.Rating)
		if descending {
			return -result
		}
		return result
	})
}

// SortBooksByGenre sorts books by first genre, then author, then series/sequence
func SortBooksByGenre(books []Book, descending bool) {
	sortBooksByFallbackOrder(books)
	slices.SortStableFunc(books, func(a, b Book) int {
		// If either book has no genre, sort it to the end
		if len(a.Genre) == 0 && len(b.Genre) == 0 {
			return 0
		} else if len(a.Genre) == 0 {
			return 1 // a has no genre, sort it after b
		} else if len(b.Genre) == 0 {
			return -1 // b has no genre, sort it after a
		}
		// Both have genres, compare them
		result := strings.Compare(a.GetFirstGenre(), b.GetFirstGenre())
		if descending {
			return -result
		}
		return result
	})
}

// SortStrings sorts a slice of strings alphabetically
func SortStrings(strings []string, descending bool) {
	slices.SortFunc(strings, func(a, b string) int {
		result := cmp.Compare(a, b)
		if descending {
			return -result
		}
		return result
	})
}

// sortBooksByFallbackOrder sorts books by Series, then Sequence, then Author, then Title
func sortBooksByFallbackOrder(books []Book) {
	sortBooksByFallbackOrderSupportingDescendingSeries(books, false)
}

// sortBooksByFallbackOrder sorts books by Series, then Sequence, then Author, then Title
func sortBooksByFallbackOrderSupportingDescendingSeries(books []Book, descending bool) {
	slices.SortStableFunc(books, func(a, b Book) int {
		// If either book has no series, sort it to the end
		if a.Series == "" && b.Series == "" {
			// Both have no series, continue with other comparisons
		} else if a.Series == "" {
			return 1 // a has no series, sort it after b
		} else if b.Series == "" {
			return -1 // b has no series, sort it after a
		} else {
			// Both have series, compare them
			seriesResult := strings.Compare(a.Series, b.Series)
			if descending {
				seriesResult = -seriesResult
			}
			if seriesResult != 0 {
				return seriesResult
			}
		}

		// Compare sequence
		if sequenceResult := compareSequence(a, b); sequenceResult != 0 {
			return sequenceResult
		}

		// Compare author
		authorResult := strings.Compare(a.GetFirstAuthorSort(), b.GetFirstAuthorSort())
		if authorResult != 0 {
			return authorResult
		}

		// Compare title
		return strings.Compare(a.Title, b.Title)
	})
}

// compareSequence compares two books by their sequence
func compareSequence(a, b Book) int {
	// If either sequence is empty, sort it to the end
	if a.Sequence == "" && b.Sequence == "" {
		return 0
	} else if a.Sequence == "" {
		return 1 // a has no sequence, sort it after b
	} else if b.Sequence == "" {
		return -1 // b has no sequence, sort it after a
	}

	// Try to parse both sequences as numbers
	aNum, aErr := strconv.Atoi(a.Sequence)
	bNum, bErr := strconv.Atoi(b.Sequence)

	// If both are pure numbers, compare them numerically
	if aErr == nil && bErr == nil {
		return cmp.Compare(aNum, bNum)
	}

	// Split into number and text parts
	aNumStr := ""
	aText := ""
	for i, c := range a.Sequence {
		if c >= '0' && c <= '9' {
			aNumStr += string(c)
		} else {
			aText = a.Sequence[i:]
			break
		}
	}

	bNumStr := ""
	bText := ""
	for i, c := range b.Sequence {
		if c >= '0' && c <= '9' {
			bNumStr += string(c)
		} else {
			bText = b.Sequence[i:]
			break
		}
	}

	// Convert numeric parts
	aNum, _ = strconv.Atoi(aNumStr)
	bNum, _ = strconv.Atoi(bNumStr)

	// Compare numeric parts first
	if result := cmp.Compare(aNum, bNum); result != 0 {
		return result
	}

	// If numeric parts are equal, pure numbers sort before alphanumeric
	if aErr == nil {
		return -1 // a is pure number, sort it before b
	}
	if bErr == nil {
		return 1 // b is pure number, sort it before a
	}

	// Both are alphanumeric with same number, compare text parts
	return strings.Compare(aText, bText)
}

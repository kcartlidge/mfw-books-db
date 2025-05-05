package main

// BookFilter represents a filtered view of books
type BookFilter struct {
	// Name is the display name of the filter (e.g. "All Books", "Unread Books", etc.)
	Name string

	// Books is the filtered collection of books
	Books []Book

	// Populate takes a source collection of books and populates this filter's Books
	// collection by applying its specific filtering logic
	Populate func(source []Book)
}

// GetPopulatedAllBooksFilter returns a BookFilter that contains all books
func GetPopulatedAllBooksFilter(books []Book) BookFilter {
	var filter BookFilter
	filter = BookFilter{
		Name: "All Books",
		Populate: func(source []Book) {
			// Simply copy all books from the source
			filter.Books = make([]Book, len(source))
			copy(filter.Books, source)
		},
	}
	// Populate the filter with the provided books
	filter.Populate(books)
	return filter
}

// GetPopulatedReadingFilter returns a BookFilter that contains only books currently being read
func GetPopulatedReadingFilter(books []Book) BookFilter {
	var filter BookFilter
	filter = BookFilter{
		Name: "Reading",
		Populate: func(source []Book) {
			filter.Books = filterByStatusIcon(source, "C")
		},
	}
	// Populate the filter with the provided books
	filter.Populate(books)
	return filter
}

// GetPopulatedNextFilter returns a BookFilter that contains only books marked as next to read
func GetPopulatedNextFilter(books []Book) BookFilter {
	var filter BookFilter
	filter = BookFilter{
		Name: "Next",
		Populate: func(source []Book) {
			filter.Books = filterByStatusIcon(source, "N")
		},
	}
	// Populate the filter with the provided books
	filter.Populate(books)
	return filter
}

// GetPopulatedDoneFilter returns a BookFilter that contains only completed books
func GetPopulatedDoneFilter(books []Book) BookFilter {
	var filter BookFilter
	filter = BookFilter{
		Name: "Done",
		Populate: func(source []Book) {
			filter.Books = filterByStatusIcon(source, "R", "A")
		},
	}
	// Populate the filter with the provided books
	filter.Populate(books)
	return filter
}

// GetPopulatedOtherFilter returns a BookFilter that contains books with other statuses
func GetPopulatedOtherFilter(books []Book) BookFilter {
	var filter BookFilter
	filter = BookFilter{
		Name: "Other",
		Populate: func(source []Book) {
			// Get all unique status icons from the source
			statusMap := make(map[string]bool)
			for _, book := range source {
				if book.StatusIcon != "" {
					statusMap[book.StatusIcon] = true
				}
			}
			// Remove the main status icons
			delete(statusMap, "C") // Current
			delete(statusMap, "N") // Next
			delete(statusMap, "R") // Read
			delete(statusMap, "A") // Abandoned
			// Convert remaining status icons to slice
			otherStatuses := make([]string, 0, len(statusMap))
			for status := range statusMap {
				otherStatuses = append(otherStatuses, status)
			}
			// Use the private function to filter by the remaining status icons
			filter.Books = filterByStatusIcon(source, otherStatuses...)
		},
	}
	// Populate the filter with the provided books
	filter.Populate(books)
	return filter
}

// filterByStatusIcon returns a slice of books that match any of the provided status icons
func filterByStatusIcon(source []Book, statusIcons ...string) []Book {
	// Create a slice to hold the filtered books
	filtered := make([]Book, 0)
	// Add only books with matching StatusIcon
	for _, book := range source {
		for _, icon := range statusIcons {
			if book.StatusIcon == icon {
				filtered = append(filtered, book)
				break
			}
		}
	}
	return filtered
}

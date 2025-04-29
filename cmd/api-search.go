package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

// GoogleBook represents the book data we get from Google Books API
type GoogleBook struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	Authors             []string `json:"authors"`
	Link                string   `json:"selfLink"`
	Categories          []string `json:"categories"`
	PublishedDate       string   `json:"publishedDate"`
	Description         string   `json:"description"`
	PageCount           int      `json:"pageCount"`
	Language            string   `json:"language"`
	Publisher           string   `json:"publisher"`
	IndustryIdentifiers []struct {
		Type       string `json:"type"`
		Identifier string `json:"identifier"`
	} `json:"industryIdentifiers"`
}

// GetBookByISBN queries Google Books API for a book by ISBN
func GetBookByISBN(isbn string, singleHit bool) (*GoogleBook, error) {
	baseURL := "https://www.googleapis.com/books/v1/volumes"
	query := fmt.Sprintf("isbn:%s", isbn)

	params := url.Values{}
	params.Add("q", query)

	// Rate limiting: 333ms delay between requests (3 requests per second)
	time.Sleep(333 * time.Millisecond)

	// Up to 3 retries for failed requests
	var lastErr error
	for retry := 0; retry < 3; retry++ {
		// Make the request to the Google Books API
		resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
		if err != nil {
			lastErr = err
			// 2-second delay between retries
			time.Sleep(2 * time.Second)
			continue
		}

		// Handle rate limit errors (HTTP 429)
		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			lastErr = errors.New("rate limit exceeded")
			// 2-second delay between retries
			time.Sleep(2 * time.Second)
			continue
		}

		// Decode the response into a struct
		var result struct {
			Items []struct {
				ID         string     `json:"id"`
				SelfLink   string     `json:"selfLink"`
				VolumeInfo GoogleBook `json:"volumeInfo"`
			} `json:"items"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			lastErr = err
			// 2-second delay between retries
			time.Sleep(2 * time.Second)
			continue
		}
		resp.Body.Close()

		// Check if the response contains any items
		if len(result.Items) == 0 {
			return nil, errors.New("no book found")
		}

		// Get the first item from the response
		book := result.Items[0].VolumeInfo
		book.ID = result.Items[0].ID
		book.Link = result.Items[0].SelfLink

		// Get the further details for the book if we are not using single hit calls
		// Swallow any error as this is not critical
		if !singleHit {
			err = getFurtherBookDetailsForGoogleBook(&book)
			if err != nil {
				fmt.Printf("Error getting further book details for %s: %s\n", book.Title, err.Error())
			}
		}

		// Return the book
		return &book, nil
	}

	return nil, lastErr
}

// getFurtherBookDetailsForGoogleBook re-queries Google Books API
// by book ID as this often returns more accurate genre information
// and also offers a more accurate publisher
func getFurtherBookDetailsForGoogleBook(book *GoogleBook) error {
	// Get the further details for the book
	resp, err := http.Get(book.Link)
	if err != nil {
		return err
	}

	// Decode the response into a struct
	var result struct {
		ID         string     `json:"id"`
		SelfLink   string     `json:"selfLink"`
		VolumeInfo GoogleBook `json:"volumeInfo"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return err
	}
	resp.Body.Close()

	// Check if the response contains any items
	if result.ID != book.ID {
		return errors.New("no book found")
	}

	// Update the book with the extra details
	mergeGenres(book, result.VolumeInfo.Categories)
	if strings.TrimSpace(result.VolumeInfo.Publisher) != "" {
		book.Publisher = result.VolumeInfo.Publisher
	}
	if result.VolumeInfo.PageCount > 0 {
		book.PageCount = result.VolumeInfo.PageCount
	}

	return nil
}

// mergeGenres merges the new genres with the existing ones
// It also attempts to de-dupe and simplify the list
func mergeGenres(book *GoogleBook, newGenres []string) {
	// New ones get priority as they are more reliable
	mergedGenres := append(newGenres, book.Categories...)
	book.Categories = []string{}

	// De-dupe genres, where we allow up to 2 levels of
	// genre before dropping the remainder
	for _, genre := range mergedGenres {
		bits := strings.Split(genre, "/")
		for i := 0; i < len(bits); i++ {
			bits[i] = strings.TrimSpace(bits[i])
		}
		newGenre := bits[0]
		if len(bits) > 1 {
			newGenre = fmt.Sprintf("%s, %s", bits[0], bits[1])
		}
		if !slices.Contains(book.Categories, newGenre) {
			book.Categories = append(book.Categories, newGenre)
		}
	}
}

// MapGoogleBookToBook maps a Google Book API response to our Book struct
func MapGoogleBookToBook(gb *GoogleBook) *Book {
	book := &Book{
		ID:              gb.ID,
		ISBN:            getISBNFromIdentifiers(gb.IndustryIdentifiers),
		Title:           gb.Title,
		Authors:         gb.Authors,
		Genre:           gb.Categories,
		Link:            gb.Link,
		PublishedDate:   gb.PublishedDate,
		Publisher:       gb.Publisher,
		PageCount:       gb.PageCount,
		Language:        gb.Language,
		Description:     gb.Description,
		AuthorSort:      []string{},
		Series:          "",
		Sequence:        "",
		Status:          "",
		Rating:          0,
		Notes:           "",
		StatusIcon:      "",
		ModifiedUtc:     "",
		IsException:     false,
		ExceptionReason: "",
	}

	return book
}

// getISBNFromIdentifiers extracts the ISBN from industry identifiers
func getISBNFromIdentifiers(identifiers []struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}) string {
	for _, id := range identifiers {
		if id.Type == "ISBN_13" || id.Type == "ISBN_10" {
			return id.Identifier
		}
	}
	return ""
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// GoogleBook represents the book data we get from Google Books API
type GoogleBook struct {
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
func GetBookByISBN(isbn string) (*GoogleBook, error) {
	baseURL := "https://www.googleapis.com/books/v1/volumes"
	query := fmt.Sprintf("isbn:%s", isbn)

	params := url.Values{}
	params.Add("q", query)

	// Make the request to the Google Books API
	resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the response into a struct
	var result struct {
		Items []struct {
			VolumeInfo GoogleBook `json:"volumeInfo"`
			SelfLink   string     `json:"selfLink"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Check if the response contains any items
	if len(result.Items) == 0 {
		return nil, errors.New("no book found")
	}

	// Get the first item from the response
	book := result.Items[0].VolumeInfo
	book.Link = result.Items[0].SelfLink

	// Return the book
	return &book, nil
}

// MapGoogleBookToBook maps a Google Book API response to our Book struct
func MapGoogleBookToBook(gb *GoogleBook) *Book {
	book := &Book{
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

package main

import (
	"encoding/json"
	"fmt"
)

// Book represents a book in the database
type Book struct {
	ISBN            string   `json:"isbn"`
	Title           string   `json:"title"`
	Authors         []string `json:"authors"`
	Genre           []string `json:"genre"`
	Link            string   `json:"link"`
	PublishedDate   string   `json:"publishedDate"`
	Publisher       string   `json:"publisher"`
	PageCount       int      `json:"pageCount"`
	Language        string   `json:"language"`
	Description     string   `json:"description"`
	AuthorSort      []string `json:"authorSort"`
	Series          string   `json:"series"`
	Sequence        string   `json:"sequence"`
	Status          string   `json:"status"`
	Rating          int      `json:"rating"`
	Notes           string   `json:"notes"`
	StatusIcon      string   `json:"statusIcon"`
	ModifiedUtc     string   `json:"modifiedUtc"`
	IsException     bool     `json:"isException"`
	ExceptionReason string   `json:"exceptionReason"`
}

// GetAuthorDisplay returns a formatted string for displaying authors
func (b *Book) GetAuthorDisplay() string {
	return b.getDisplayString(b.Authors)
}

// GetAuthorSortDisplay returns a formatted string for displaying author sorts
func (b *Book) GetAuthorSortDisplay() string {
	return b.getDisplayString(b.AuthorSort)
}

// GetGenreDisplay returns a formatted string for displaying genres
func (b *Book) GetGenreDisplay() string {
	return b.getDisplayString(b.Genre)
}

// MarshalJSON implements json.Marshaler for Book
func (b *Book) MarshalJSON() ([]byte, error) {
	type Alias Book
	return json.Marshal(&struct {
		*Alias
		SeriesSort string `json:"seriesSort"`
	}{
		Alias:      (*Alias)(b),
		SeriesSort: b.ComputeSeriesSort(),
	})
}

// Print displays the book's information
func (b *Book) Print() {
	grid := NewGrid([]string{"Property", "Value"})
	grid.SetShowHeaders(false)

	grid.AddRow("ISBN:", b.ISBN)
	grid.AddRow("Title:", b.Title)
	grid.AddRow("Authors:", joinWithAmpersand(b.Authors))
	grid.AddRow("Genres:", joinWithAmpersand(b.Genre))
	grid.AddRow("Published:", b.PublishedDate)
	grid.AddRow("Publisher:", b.Publisher)
	grid.AddRow("Pages:", fmt.Sprintf("%d", b.PageCount))
	grid.AddRow("Language:", b.Language)
	grid.AddRow("Description:", b.Description)
	grid.AddRow("Author Sort:", joinWithAmpersand(b.AuthorSort))
	grid.AddRow("Series:", b.Series)
	grid.AddRow("Sequence:", b.Sequence)
	grid.AddRow("Link:", b.Link)
	grid.AddRow("Status:", b.Status)
	grid.AddRow("Rating:", fmt.Sprintf("%d", b.Rating))
	grid.AddRow("Notes:", b.Notes)
	grid.AddRow("Status Icon:", b.StatusIcon)
	grid.AddRow("Series Sort:", b.ComputeSeriesSort())
	grid.AddRow("Modified:", b.ModifiedUtc)
	grid.AddRow("Exception:", fmt.Sprintf("%v", b.IsException))
	grid.AddRow("Exception Reason:", b.ExceptionReason)

	fmt.Println(grid)
}

// ComputeSeriesSort returns the computed series sort value
func (b *Book) ComputeSeriesSort() string {
	if b.Series == "" {
		return ""
	}
	if b.Sequence == "" {
		return b.Series
	}
	return fmt.Sprintf("%s [%s]", b.Series, b.Sequence)
}

// getDisplayString returns a formatted string for displaying array items
func (b *Book) getDisplayString(items []string) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}
	return fmt.Sprintf("%s [+%d]", items[0], len(items)-1)
}

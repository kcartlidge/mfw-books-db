package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Book represents a book in the database
type Book struct {
	ISBN            string   `json:"isbn"`
	Title           string   `json:"title"`
	Authors         string   `json:"authors"`
	Genre           []string `json:"genre"`
	Link            string   `json:"link"`
	PublishedDate   string   `json:"publishedDate"`
	Publisher       string   `json:"publisher"`
	PageCount       int      `json:"pageCount"`
	Language        string   `json:"language"`
	Description     string   `json:"description"`
	AuthorSort      string   `json:"authorSort"`
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
	grid.AddRow("Authors:", b.Authors)
	grid.AddRow("Genres:", strings.Join(b.Genre, ", "))
	grid.AddRow("Published:", b.PublishedDate)
	grid.AddRow("Publisher:", b.Publisher)
	grid.AddRow("Pages:", fmt.Sprintf("%d", b.PageCount))
	grid.AddRow("Language:", b.Language)
	grid.AddRow("Description:", b.Description)
	grid.AddRow("Author Sort:", b.AuthorSort)
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

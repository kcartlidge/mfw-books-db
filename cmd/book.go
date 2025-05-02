package main

import (
	"fmt"
	"html/template"
	"strings"
)

// Book represents a book in the database
type Book struct {
	ID              string   `json:"id"`
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

// GetSeriesSort returns the computed series sort value
func (b *Book) GetSeriesSort() string {
	if b.Series == "" {
		return ""
	}
	if b.Sequence == "" {
		return b.Series
	}
	return fmt.Sprintf("%s [%s]", b.Series, b.Sequence)
}

// GetStatusLetter returns the status icon letter for the book
func (b *Book) GetStatusLetter() string {
	if b.StatusIcon == "" {
		return "-"
	}
	return b.StatusIcon
}

// GetFirstAuthorSort returns the first AuthorSort of the book
func (b *Book) GetFirstAuthorSort() string {
	if len(b.AuthorSort) == 0 {
		return ""
	}
	return b.AuthorSort[0]
}

// GetFirstGenre returns the first Genre of the book
func (b *Book) GetFirstGenre() string {
	if len(b.Genre) == 0 {
		return ""
	}
	return b.Genre[0]
}

// GetGenresForEdit returns a formatted string for editing as ... & ... & ...
func (b *Book) GetGenresForEdit() string {
	return b.getLines(b.Genre, " & ")
}

// GetAuthorDisplay returns a formatted string for displaying authors
func (b *Book) GetAuthorDisplay() string {
	return b.getDisplayString(b.Authors)
}

// GetAuthorSortDisplay returns a formatted string for displaying author sorts as ...[+n]
func (b *Book) GetAuthorSortDisplay() string {
	return b.getDisplayString(b.AuthorSort)
}

// GetAuthorsForEdit returns a formatted string for editing as ... & ... & ...
func (b *Book) GetAuthorsForEdit() string {
	return b.getLines(b.Authors, " & ")
}

// GetAuthorSortForEdit returns a formatted string for editing as ... & ... & ...
func (b *Book) GetAuthorSortForEdit() string {
	return b.getLines(b.AuthorSort, " & ")
}

// GetAuthorSortHtmlDisplay returns a formatted string for displaying author sorts as lines
func (b *Book) GetAuthorSortHtmlDisplay() template.HTML {
	return template.HTML(b.getHtmlLines(b.AuthorSort))
}

// GetGenreDisplay returns a formatted string for displaying genres
func (b *Book) GetGenreDisplay() string {
	return b.getDisplayString(b.Genre)
}

// GetGenreHtmlDisplay returns a formatted string for displaying genres
func (b *Book) GetGenreHtmlDisplay() template.HTML {
	return template.HTML(b.getHtmlLines(b.Genre))
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
	grid.AddRow("Author Sort:", joinWithAmpersand(b.AuthorSort))
	grid.AddRow("Series:", b.Series)
	grid.AddRow("Sequence:", b.Sequence)
	grid.AddRow("Link:", b.Link)
	grid.AddRow("Status:", b.Status)
	grid.AddRow("Status Icon:", b.StatusIcon)
	grid.AddRow("Rating:", fmt.Sprintf("%d", b.Rating))
	grid.AddRow("Series Sort:", b.GetSeriesSort())
	grid.AddRow("Description:", b.Description)
	grid.AddRow("Notes:", b.Notes)
	grid.AddRow("Exception:", fmt.Sprintf("%v", b.IsException))
	grid.AddRow("Exception Reason:", b.ExceptionReason)
	grid.AddRow("Modified:", b.ModifiedUtc)

	fmt.Println(grid)
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

// getHtmlLines returns HTML for the array items as a list of lines
func (b *Book) getHtmlLines(items []string) string {
	return b.getLines(items, "<br/>")
}

// getLines returns a string for the array items as a list of lines with a join string
func (b *Book) getLines(items []string, joinWith string) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}
	return strings.Join(items, joinWith)
}

// GetLinkGoodreads returns the link for the book on Goodreads
func (b *Book) GetLinkGoodreads() string {
	return fmt.Sprintf("https://www.goodreads.com/search?q=%s", b.ISBN)
}

// GetLinkGoogleBooksJson returns the link for fetching the book details as a JSON string
func (b *Book) GetLinkGoogleBooksJson() string {
	return b.Link
}

// GetLinkGoogleBooksView returns the link for the book as a HTML page on Google Books
func (b *Book) GetLinkGoogleBooksView() string {
	return fmt.Sprintf("https://books.google.com/books?id=%s&dq=isbn:%s", b.ID, b.ISBN)
}

// GetLinkOpenLibrary returns the link for the book on OpenLibrary
func (b *Book) GetLinkOpenLibrary() string {
	return fmt.Sprintf("https://openlibrary.org/isbn/%s", b.ISBN)
}

// GetLinkLibraryThing returns the link for the book on LibraryThing
func (b *Book) GetLinkLibraryThing() string {
	return fmt.Sprintf("https://www.librarything.com/search.php?search=%s", b.ISBN)
}

// GetLinkWaterstones returns the link for the book on Waterstones
func (b *Book) GetLinkWaterstones() string {
	return fmt.Sprintf("https://www.waterstones.com/index/search?term=%s", b.ISBN)
}

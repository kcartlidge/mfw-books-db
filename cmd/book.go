package main

// Book represents a book in the database
type Book struct {
	ISBN            string   `json:"isbn"`
	Title           string   `json:"title"`
	Authors         string   `json:"authors"`
	AuthorSort      string   `json:"authorSort"`
	Series          *string  `json:"series"`
	Sequence        *string  `json:"sequence"`
	Genre           []string `json:"genre"`
	Link            string   `json:"link"`
	IsException     bool     `json:"isException"`
	ExceptionReason string   `json:"exceptionReason"`
	ModifiedUtc     string   `json:"modifiedUtc"`
	Status          string   `json:"status"`
	Rating          *int     `json:"rating"`
	Notes           *string  `json:"notes"`
	StatusIcon      string   `json:"statusIcon"`
	SeriesSort      string   `json:"seriesSort"`
}

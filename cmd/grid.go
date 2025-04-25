package main

import (
	"fmt"
	"strings"
)

const (
	ColumnPadding = 3
	MaxCellWidth  = 50
	LineChar      = "="
)

// Grid represents a table of text content with fixed columns
type Grid struct {
	headers     []string
	rows        [][]string
	widths      []int
	showHeaders bool
	showNumbers bool
	numberWidth int
}

// NewGrid creates a new Grid with the given column headers
func NewGrid(headers []string) *Grid {
	g := &Grid{
		headers:     headers,
		widths:      make([]int, len(headers)),
		showHeaders: true,  // Default to showing headers
		showNumbers: false, // Default to not showing row numbers
		numberWidth: 0,     // Will be calculated when needed
	}
	// Initialize widths based on headers
	for i, header := range headers {
		g.widths[i] = min(len(header), MaxCellWidth)
	}
	return g
}

// SetShowHeaders sets whether headers should be displayed
func (g *Grid) SetShowHeaders(show bool) {
	g.showHeaders = show
}

// SetShowNumbers sets whether row numbers should be displayed
func (g *Grid) SetShowNumbers(show bool) {
	g.showNumbers = show
}

// AddRow adds a row of cells to the grid
func (g *Grid) AddRow(cells ...string) {
	// Ensure we have enough cells
	row := make([]string, len(g.headers))
	copy(row, cells)

	// Update column widths and truncate cells
	for i, cell := range row {
		row[i] = truncateCell(cell)
		if len(row[i]) > g.widths[i] {
			g.widths[i] = min(len(row[i]), MaxCellWidth)
		}
	}

	g.rows = append(g.rows, row)
}

// String returns the grid as a formatted string
func (g *Grid) String() string {
	var sb strings.Builder

	// Calculate number width if showing numbers
	if g.showNumbers {
		// Calculate width needed for the largest row number
		maxNum := len(g.rows)
		g.numberWidth = len(fmt.Sprintf("%d", maxNum)) + ColumnPadding
	}

	if g.showHeaders {
		// Print headers
		if g.showNumbers {
			sb.WriteString(strings.Repeat(" ", g.numberWidth))
		}
		for i, header := range g.headers {
			if i > 0 {
				sb.WriteString(strings.Repeat(" ", ColumnPadding))
			}
			fmt.Fprintf(&sb, "%-*s", g.widths[i], truncateCell(header))
		}
		sb.WriteString("\n")

		// Print separator
		if g.showNumbers {
			sb.WriteString(strings.Repeat(LineChar, g.numberWidth))
		}
		for i, width := range g.widths {
			if i > 0 {
				sb.WriteString(strings.Repeat(LineChar, ColumnPadding))
			}
			sb.WriteString(strings.Repeat(LineChar, width))
		}
		sb.WriteString("\n")
	}

	// Print rows
	for rowNum, row := range g.rows {
		if g.showNumbers {
			fmt.Fprintf(&sb, "%*d%s", g.numberWidth-ColumnPadding, rowNum+1, strings.Repeat(" ", ColumnPadding))
		}
		for i, cell := range row {
			if i > 0 {
				sb.WriteString(strings.Repeat(" ", ColumnPadding))
			}
			fmt.Fprintf(&sb, "%-*s", g.widths[i], cell)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// truncateCell truncates a cell's content if it exceeds MaxCellWidth, respecting word boundaries
func truncateCell(s string) string {
	if len(s) <= MaxCellWidth {
		return s
	}

	// Find the last space before MaxCellWidth-4
	truncPoint := MaxCellWidth - 4
	lastSpace := strings.LastIndex(s[:truncPoint+1], " ")

	// If we found a space and removing the partial word wouldn't remove too much text
	if lastSpace > 0 && lastSpace > truncPoint/2 {
		return s[:lastSpace] + " ..."
	}

	// Otherwise truncate at the original point
	return s[:truncPoint] + " ..."
}

package main

import (
	"fmt"
	"strings"
)

const (
	ColumnPadding = 3
	MaxCellWidth  = 50
)

// Grid represents a table of text content with fixed columns
type Grid struct {
	headers     []string
	rows        [][]string
	widths      []int
	showHeaders bool
}

// NewGrid creates a new Grid with the given column headers
func NewGrid(headers []string) *Grid {
	g := &Grid{
		headers:     headers,
		widths:      make([]int, len(headers)),
		showHeaders: true, // Default to showing headers
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

	if g.showHeaders {
		// Print headers
		for i, header := range g.headers {
			if i > 0 {
				sb.WriteString(strings.Repeat(" ", ColumnPadding))
			}
			fmt.Fprintf(&sb, "%-*s", g.widths[i], truncateCell(header))
		}
		sb.WriteString("\n")

		// Print separator
		for i, width := range g.widths {
			if i > 0 {
				sb.WriteString(strings.Repeat("-", ColumnPadding))
			}
			sb.WriteString(strings.Repeat("-", width))
		}
		sb.WriteString("\n")
	}

	// Print rows
	for _, row := range g.rows {
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

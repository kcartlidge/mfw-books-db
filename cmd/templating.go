package main

import (
	"errors"
	"html/template"
	"io"
	"path/filepath"
)

// TemplateData represents the data passed to templates
type TemplateData struct {
	Title     string
	Filename  string
	Content   interface{}
	SortField string
}

// Templates holds all our templates
type Templates struct {
	top   *template.Template
	base  *template.Template
	pages map[string]*template.Template
}

// Global template manager
var templates *Templates

// NewTemplates creates a new template manager
func NewTemplates() (*Templates, error) {
	// If we already have templates, return them
	if templates != nil {
		return templates, nil
	}

	t := &Templates{
		pages: make(map[string]*template.Template),
	}

	// Load the top template
	top, err := template.ParseFiles("templates/top.go.html")
	if err != nil {
		return nil, err
	}
	t.top = top

	// Load the base template
	base, err := template.ParseFiles("templates/base.go.html")
	if err != nil {
		return nil, err
	}
	t.base = base

	// Load all page templates
	pages, err := filepath.Glob("templates/pages/*.go.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Create a new template for each page
		tmpl := template.New(filepath.Base(page))
		tmpl.Funcs(template.FuncMap{
			"For": func(count int) []int {
				var Items []int
				for i := range count {
					Items = append(Items, i)
				}
				return Items
			},
		})

		// Parse all templates together
		combined, err := tmpl.ParseFiles(
			"templates/top.go.html",
			"templates/base.go.html",
			page,
		)
		if err != nil {
			return nil, err
		}

		// Store the combined template
		t.pages[filepath.Base(page)] = combined
	}

	// Store the templates globally
	templates = t
	return t, nil
}

// Render renders a template with the given data
func (t *Templates) Render(w io.Writer, name string, data TemplateData) error {
	tmpl, ok := t.pages[name]
	if !ok {
		return errors.New("template not found: " + name)
	}

	// Execute the home template, which will include top and base as needed
	return tmpl.ExecuteTemplate(w, "home", data)
}

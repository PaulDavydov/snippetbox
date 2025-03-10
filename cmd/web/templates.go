package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.pauldvyd.net/internal/models"
)

// Holds dynamic data we want to pass into HTML templates
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// use filepath.Glob() to get slice of all filepaths that math pattern
	// "./html/pages/*.tmpl"
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// ignore the index, get the actual value
	for _, page := range pages {
		name := filepath.Base(page)

		// parse base tmpl into template set
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// parse partials folder into template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Parse the page template into the template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add template set to the map, using the anme of the page as the key
		cache[name] = ts
	}

	return cache, nil
}

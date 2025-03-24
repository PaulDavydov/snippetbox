package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.pauldvyd.net/internal/models"
)

// Holds dynamic data we want to pass into HTML templates
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
}

// custom function used within our templates. Can have as many parameters as we want
// but can only return one value. Time formatting is done by using January 02 2006 15:04
// as the default. The time package will then adjust properly to your current time
func currDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Map that acts as a global variable. Stores the names and lookup of our
// custom template functions
var functions = template.FuncMap{
	"currDate": currDate,
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

		// register the FuncMap into the template set before calling any ParseFiles() methods.
		// parse base tmpl into template set
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
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

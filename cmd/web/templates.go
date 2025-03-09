package main

import "snippetbox.pauldvyd.net/internal/models"

// Holds dynamic data we want to pass into HTML templates
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

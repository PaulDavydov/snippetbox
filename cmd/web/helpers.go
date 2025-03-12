package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// writes a error message, then a stack trace to the errorLog,
// then sends a generic 500 Internal Server Error
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// sends a specific status code and corresponding description to the user
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound helper, sends a 404 not found response
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// render template pages
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Grab the correct template from the cache, based on the page name
	ts, err := app.templateCache[page]
	if !err {
		err := fmt.Errorf("the template %s does not exits", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	// Write template to the buffer
	terr := ts.ExecuteTemplate(buf, "base", data)
	if terr != nil {
		app.serverError(w, terr)
		return
	}

	// write out any provided HTTP status code
	w.WriteHeader(status)

	// write contents of the buffer to the http.ResponseWriter
	buf.WriteTo(w)
}

// TemplateData helper that returns a pointer to a templateData struct initialized
// with the current year. View as a constructor
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

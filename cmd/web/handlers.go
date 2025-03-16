package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.pauldvyd.net/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Initiate snippets, grabbing latest snippets created
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Call newTemplateData() helper to get templateData struct containing the current year
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// pass data to render method
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// httprouter is parsing a request, the values of the named parameters will be stored
	// in the request context.
	params := httprouter.ParamsFromContext(r.Context())

	// use the httprouter ByName method to the value of a specific parameter, for us it is id
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// use the SnippetModel get method to retrieve the data for a specific record
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// call newTemplateData method to create a templateData struct containing current year
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet...."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	// pass data into SnippetModel.Insert() method
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// redirect user to relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) snippetDotCheck(w http.ResponseWriter, r *http.Request) {
	app.clientError(w, http.StatusTeapot)
}

package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.pauldvyd.net/internal/models"
	"snippetbox.pauldvyd.net/internal/validator"
)

// represents the form data and validation errors for the form fields
// All fields are exported, due to the capital letters. Must be exported for the
// html/template package to read
// use third party package to decode the form data.
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

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
	data := app.newTemplateData(r)

	// Initialize a new createSnippetForm struct and pass it into the template
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	// Call the decode method of the decoder, passing in the current request and a pointer
	// to our snippetCreateForm struct. This will fill in our struct with the relevant values
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the embedded Validator type methods to call the checks on the form data
	// add the errors to the Validator map via the CheckField() method
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// use the isValid() method to see if any checks failed. If they did, re-render the template with the data
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	// pass data into SnippetModel.Insert() method
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// redirect user to relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetDotCheck(w http.ResponseWriter, r *http.Request) {
	app.clientError(w, http.StatusTeapot)
}

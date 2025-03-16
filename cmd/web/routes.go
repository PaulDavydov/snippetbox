package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// route app calls through middleware.secureHeaders()
func (app *application) routes() http.Handler {
	// handle URLs and corresponding app functions
	router := httprouter.New()

	// wrap our notFound() helper and assign it a custom handler for the 404 Not Found
	// responses. Can set set for other custom handlers as well.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.methodNotAllowed(w)
	})

	// have router handle static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// router now maps the appropriate method, patterns, and handlers
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)
	router.HandlerFunc(http.MethodGet, "/.env", app.snippetDotCheck)
	router.HandlerFunc(http.MethodGet, "/snippet/.env", app.snippetDotCheck)

	// use the alice package to chain the middleware in a more readable way
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// wrap router in the middleware and return it as normal
	return standard.Then(router)
}

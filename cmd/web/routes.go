package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// route app calls through middleware.secureHeaders()
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// use the alice package to chain the middleware in a more readable way
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// passes servermux as the 'next' parameter for the middleware.secureHeaders()
	return standard.Then(mux)
}

package main

import "net/http"

// route app calls through middleware.secureHeaders()
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// passes servermux as the 'next' parameter for the middleware.secureHeaders()
	return secureHeaders(mux)
}

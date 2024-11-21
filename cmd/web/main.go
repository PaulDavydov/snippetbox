package main
import (
  "log"
  "net/http"
  "flag"
)

func main() {

  // Hardcodes the command line flag with the name addr.
  addr := flag.String("addr", ":4000", "HTTP network address")

  flag.Parse()

  mux := http.NewServeMux()
  
  // Serves files out of the static directory
  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  mux.HandleFunc("/", home)
  mux.HandleFunc("/snippet/view", snippetView)
  mux.HandleFunc("/snippet/create", snippetCreate)

  // Add address passed in as argument to be printed here
  log.Printf("Starting server on %s", *addr)
  err := http.ListenAndServe(*addr, mux)
  log.Fatal(err)
}

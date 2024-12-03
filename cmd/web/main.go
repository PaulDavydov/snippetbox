package main
import (
  "log"
  "net/http"
  "flag"
  "os"
)

func main() {

  // Hardcodes the command line flag with the name addr.
  addr := flag.String("addr", ":4000", "HTTP network address")

  flag.Parse()

  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  mux := http.NewServeMux()
  
  // Serves files out of the static directory
  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  mux.HandleFunc("/", home)
  mux.HandleFunc("/snippet/view", snippetView)
  mux.HandleFunc("/snippet/create", snippetCreate)

  // Add address passed in as argument to be printed here
  infoLog.Printf("Starting server on %s", *addr)
  err := http.ListenAndServe(*addr, mux)
  errorLog.Fatal(err)
}

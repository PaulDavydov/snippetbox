package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"snippetbox.pauldvyd.net/internal/models"

	"github.com/jackc/pgx/v5"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {
	// Setup custom logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Load env variables
	envErr := godotenv.Load()
	if envErr != nil {
		errorLog.Fatal(envErr)
	}

	// Hardcodes the command line flag and build db conenction URL
	dbConnStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_URL"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	dsn := flag.String("dsn", dbConnStr, "PostgreSQL data source name")
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	// connect to postgreSQL db
	db, dberr := openDB(*dsn)
	if dberr != nil {
		errorLog.Fatal(dberr)
	}
	defer db.Close(context.Background())

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{Conn: db},
	}

	// Initialize http.Server struct and add ErrorLog field
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Add address passed in as argument to be printed here
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}

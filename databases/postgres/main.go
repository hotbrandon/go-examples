package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN is not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Println("Successfully connected to the database!")

	http.HandleFunc("/actors/count", actorCountHandler(db))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}

}

// actorCountHandler returns an http.HandlerFunc that closes over the db pool.
func actorCountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rowCount := 0
		// This query might fail if the DB connection is lost at runtime.
		err := db.QueryRow("SELECT COUNT(*) FROM actor").Scan(&rowCount)
		if err != nil {
			// Log the detailed error for server-side observability.
			log.Printf("ERROR: could not query database: %v", err)

			// Return a generic error to the client.
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Number of actors: %d\n", rowCount)
	}
}

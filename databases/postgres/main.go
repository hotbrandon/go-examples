package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	r.GET("/actors/count", actorCountHandler(db))

	r.Run("0.0.0.0:80800")

}

// actorCountHandler returns an http.HandlerFunc that closes over the db pool.
func actorCountHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rowCount := 0
		// This query might fail if the DB connection is lost at runtime.
		err := db.QueryRow("SELECT COUNT(*) FROM actor").Scan(&rowCount)
		if err != nil {
			// Log the detailed error for server-side observability.
			log.Printf("ERROR: could not query database: %v", err)

			// Return a generic error to the client.
			c.String(500, "Internal Server Error")
			return
		}
		c.String(200, "Number of actors: %d\n", rowCount)
	}

}

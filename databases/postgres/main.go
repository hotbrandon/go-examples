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
	r.GET("/actors/", actorListHandler(db))

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

func actorListHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type Actor struct {
			ActorID   int    `json:"actor_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}

		rows, err := db.Query("SELECT actor_id, first_name, last_name FROM actor")
		if err != nil {
			log.Println("Error querying database:", err)
			c.String(500, "Internal Server Error")
			return
		}
		defer rows.Close()

		actors := []Actor{}
		for rows.Next() {
			var actor Actor
			if err := rows.Scan(&actor.ActorID, &actor.FirstName, &actor.LastName); err != nil {
				log.Println("Error scanning row:", err)
				c.String(500, "Internal Server Error")
				return
			}
			actors = append(actors, actor)
		}

		if err = rows.Err(); err != nil {
			log.Println("Error during rows iteration:", err)
			c.String(500, "Internal Server Error")
			return
		}
		c.JSON(200, actors)
	}
}

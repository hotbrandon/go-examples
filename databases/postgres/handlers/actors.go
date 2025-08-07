package handlers

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ActorCountHandler returns an http.HandlerFunc that closes over the db pool.
func ActorCountHandler(db *sql.DB) gin.HandlerFunc {
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

func ActorListHandler(db *sql.DB) gin.HandlerFunc {
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

		file, err := os.Create("actors.csv")
		if err != nil {
			log.Println("Error creating CSV file:", err)
			c.String(500, "Internal Server Error")
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		headers := []string{"Actor ID", "First Name", "Last Name"}
		if err := writer.Write(headers); err != nil {
			log.Println("Error writing headers:", err)
			c.String(500, "Internal Server Error")
			return
		}

		actors := []Actor{}
		for rows.Next() {
			var actor Actor
			if err := rows.Scan(&actor.ActorID, &actor.FirstName, &actor.LastName); err != nil {
				log.Println("Error scanning row:", err)
				c.String(500, "Internal Server Error")
				return
			}

			if err := writer.Write([]string{
				strconv.Itoa(actor.ActorID),
				actor.FirstName,
				actor.LastName,
			}); err != nil {
				log.Println("Error writing row to CSV:", err)
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

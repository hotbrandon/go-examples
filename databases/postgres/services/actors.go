package services

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type ActorService struct {
	db *sql.DB
}

func NewActorService(db *sql.DB) *ActorService {
	return &ActorService{
		db: db,
	}
}

func (s *ActorService) ActorCount() (int, error) {
	rowCount := 0
	// This query might fail if the DB connection is lost at runtime.
	err := s.db.QueryRow("SELECT COUNT(*) FROM actor").Scan(&rowCount)
	if err != nil {
		// Log the detailed error for server-side observability.
		log.Printf("ERROR: could not query database: %v", err)
		return 0, err
	}
	return rowCount, nil
}

type Actor struct {
	ActorID   int    `json:"actor_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (s *ActorService) ActorList() ([]Actor, error) {
	rows, err := s.db.Query("SELECT actor_id, first_name, last_name FROM actor")
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, err
	}
	defer rows.Close()

	file, err := os.Create("actors.csv")
	if err != nil {
		log.Println("Error creating CSV file:", err)
		return nil, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Actor ID", "First Name", "Last Name"}
	if err := writer.Write(headers); err != nil {
		log.Println("Error writing headers:", err)
		return nil, err
	}

	actors := []Actor{}
	for rows.Next() {
		var actor Actor
		if err := rows.Scan(&actor.ActorID, &actor.FirstName, &actor.LastName); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		if err := writer.Write([]string{
			strconv.Itoa(actor.ActorID),
			actor.FirstName,
			actor.LastName,
		}); err != nil {
			log.Println("Error writing row to CSV:", err)
			return nil, err
		}

		actors = append(actors, actor)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error during rows iteration:", err)
		return nil, err
	}

	return actors, nil
}

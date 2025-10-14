package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// 1. Postgres connection string
	connStr := "host=pi.hole port=5433 user=lisa password=lisa1234 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// 2. Ping to verify connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// 3. Insert now() into table tt
	var insertedTimeString string
	query := `INSERT INTO tt (ct) VALUES (now()) RETURNING ct`
	err = db.QueryRow(query).Scan(&insertedTimeString)
	if err != nil {
		log.Fatal("Insert failed:", err)
	}

	fmt.Println("Inserted UTC timestamp:", insertedTimeString)

	// 3. Generate current time in Go (Taipei local time)
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		log.Fatal("Failed to load location:", err)
	}
	now := time.Now().In(loc)

	// 4. Insert timestamp into Postgres using parameter
	query = `INSERT INTO tt (ct) VALUES ($1) RETURNING ct`
	err = db.QueryRow(query, now).Scan(&insertedTimeString)
	if err != nil {
		log.Fatal("Insert failed:", err)
	}
	fmt.Println("value of now():", now)
	fmt.Println("Inserted LOCAL timestamp:", insertedTimeString)
}

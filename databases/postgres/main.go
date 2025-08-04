package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://lisa:lisa1234@192.168.88.13:5433/dvdrental?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	rowCount := 0
	err = db.QueryRow("SELECT COUNT(*) FROM actor").Scan(&rowCount)
	if err != nil {
		fmt.Println("Error querying the database:", err)
		return
	}
	fmt.Println("Number of actors:", rowCount)
}

package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()

	dsn := os.Getenv("POSTGRES_DSN")
	db, err := sql.Open("postgres", dsn)
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

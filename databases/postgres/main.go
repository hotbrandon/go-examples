package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"postgres-demo/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("cannot find .env file: %v", err)
	}

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
	r.GET("/actors/count", handlers.ActorCountHandler(db))
	r.GET("/actors", handlers.ActorListHandler(db))

	r.Run("0.0.0.0:8080")

}

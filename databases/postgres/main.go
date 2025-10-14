package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"postgres-demo/handlers"
	"postgres-demo/services"

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

	fmt.Println("Successfully connected to the database!")

	r := gin.Default()

	actorService := services.NewActorService(db)
	actorHandler := handlers.NewActorHandler(actorService)

	r.GET("/actors/count", actorHandler.ActorCountHandler())
	r.GET("/actors", actorHandler.ActorListHandler())

	r.Run("0.0.0.0:8080")

}

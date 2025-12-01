package main

import (
	"gin-logging/middleware"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(logger))

	r.GET("/ping", func(c *gin.Context) {
		log := c.MustGet("logger").(*slog.Logger)
		log.Info("handling ping request")
		c.JSON(200, gin.H{"message": "pong"})
	})

	logger.Info("starting server", "port", 8080)

	if err := r.Run(":8080"); err != nil {
		logger.Error("server error", "error", err)
	}
}

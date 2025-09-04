package handlers

import (
	"log"
	"postgres-demo/services"

	"github.com/gin-gonic/gin"
)

type ActorHandler struct {
	svc *services.ActorService
}

func NewActorHandler(svc *services.ActorService) *ActorHandler {
	return &ActorHandler{
		svc: svc,
	}
}

// ActorCountHandler returns an http.HandlerFunc that closes over the db pool.
func (h *ActorHandler) ActorCountHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		rowCount, err := h.svc.ActorCount()
		if err != nil {
			log.Printf("ERROR: could not query database: %v", err)
			c.String(500, "Internal Server Error")
			return
		}
		c.String(200, "Number of actors: %d\n", rowCount)
	}

}

func (h *ActorHandler) ActorListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		actors, err := h.svc.ActorList()
		if err != nil {
			log.Println("Error querying database:", err)
			c.String(500, "Internal Server Error")
			return
		}
		c.JSON(200, actors)
	}
}

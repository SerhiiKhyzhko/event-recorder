package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func NewRouter(eventHandler *EventHandler, statsHandler *StatsHandler) http.Handler {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, //allows for everyone
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/events", eventHandler.CreateEvent)
	router.GET("/events", eventHandler.GetEvents)

	router.GET("/stats", statsHandler.GetStats)

	return router
}
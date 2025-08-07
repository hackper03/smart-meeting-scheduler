package routes

import (
	"smart-meeting-scheduler/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// API routes
	api := r.Group("/api/v1")
	{
		api.POST("/schedule", handlers.ScheduleMeeting)
		api.GET("/users/:userId/calendar", handlers.GetUserCalendar)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})
}

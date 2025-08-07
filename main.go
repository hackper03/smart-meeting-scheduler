package main

import (
	"os"
	"smart-meeting-scheduler/database"
	"smart-meeting-scheduler/logger"
	"smart-meeting-scheduler/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.NewLogger(logger.INFO, os.Stdout)
	database.Connect()
	r := gin.Default()
	log.Info("Starting the Gin server...")
	routes.RegisterRoutes(r)
	log.Info("Smart Meeting Scheduler API starting on port 8080...")
	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server: %v", err)
	}

}

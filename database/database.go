package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"smart-meeting-scheduler/logger"
	"smart-meeting-scheduler/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	log := logger.NewLogger(logger.INFO, os.Stdout)
	log.Info("Connecting to database...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: %v", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("SSL_MODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: %v", err)
	}

	// Automigrate
	err = DB.AutoMigrate(&models.User{}, &models.Event{}, &models.Meeting{})
	if err != nil {
		log.Fatal("Failed to migrate database: %v", err)
	}

	seedData()
}
func seedData() {
	// Check if data already exists
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return // Data already seeded
	}

	// Create users
	users := []models.User{
		{ID: "user1", Name: "Alice Johnson"},
		{ID: "user2", Name: "Bob Smith"},
		{ID: "user3", Name: "Charlie Brown"},
		{ID: "user4", Name: "Diana Ross"},
	}

	for _, user := range users {
		DB.Create(&user)
	}

	// Create some initial events
	baseTime := time.Date(2024, 9, 2, 9, 0, 0, 0, time.UTC)

	events := []models.Event{
		// Alice's events
		{UserID: "user1", Title: "Morning Standup", StartTime: baseTime, EndTime: baseTime.Add(30 * time.Minute)},
		{UserID: "user1", Title: "Client Call", StartTime: baseTime.Add(2 * time.Hour), EndTime: baseTime.Add(3 * time.Hour)},
		{UserID: "user1", Title: "Lunch Break", StartTime: baseTime.Add(3 * time.Hour), EndTime: baseTime.Add(4 * time.Hour)},

		// Bob's events
		{UserID: "user2", Title: "Team Meeting", StartTime: baseTime.Add(1 * time.Hour), EndTime: baseTime.Add(2 * time.Hour)},
		{UserID: "user2", Title: "Code Review", StartTime: baseTime.Add(4 * time.Hour), EndTime: baseTime.Add(5 * time.Hour)},

		// Charlie's events
		{UserID: "user3", Title: "Design Review", StartTime: baseTime.Add(30 * time.Minute), EndTime: baseTime.Add(90 * time.Minute)},
		{UserID: "user3", Title: "1:1 Meeting", StartTime: baseTime.Add(5 * time.Hour), EndTime: baseTime.Add(6 * time.Hour)},
	}

	for _, event := range events {
		DB.Create(&event)
	}

	log.Println("Database seeded with initial data")
}

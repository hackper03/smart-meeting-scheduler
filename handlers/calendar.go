package handlers

import (
	"net/http"
	"smart-meeting-scheduler/database"
	"smart-meeting-scheduler/errors"
	"smart-meeting-scheduler/models"
	"smart-meeting-scheduler/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUserCalendar(c *gin.Context) {
	userID := c.Param("userId")

	// Get query parameters
	startParam := c.Query("start")
	endParam := c.Query("end")

	if startParam == "" || endParam == "" {
		utils.SendError(c, errors.New(http.StatusBadRequest, "start and end query parameters are required"))
		return
	}

	// Parse time parameters
	startTime, err := time.Parse(time.RFC3339, startParam)
	if err != nil {
		utils.SendError(c, errors.New(http.StatusBadRequest, "Invalid start time format. Use RFC3339 (ISO 8601)"))
		return
	}

	endTime, err := time.Parse(time.RFC3339, endParam)
	if err != nil {
		utils.SendError(c, errors.New(http.StatusBadRequest, "Invalid end time format. Use RFC3339 (ISO 8601)"))
		return
	}

	// Fetch events
	var events []models.Event
	err = database.DB.Where("user_id = ? AND start_time < ? AND end_time > ?",
		userID, endTime, startTime).Find(&events).Error

	if err != nil {
		utils.SendError(c, errors.New(http.StatusInternalServerError, "Error fetching calendar events"))
		return
	}

	// Transform to response format
	var response []gin.H
	for _, event := range events {
		response = append(response, gin.H{
			"title":     event.Title,
			"startTime": event.StartTime.Format("2006-01-02T15:04:05Z"),
			"endTime":   event.EndTime.Format("2006-01-02T15:04:05Z"),
		})
	}

	utils.SendSuccess(c, http.StatusOK, response)
}

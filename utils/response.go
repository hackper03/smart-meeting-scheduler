package utils

import (
	"smart-meeting-scheduler/errors"

	"github.com/gin-gonic/gin"
)

// SendError sends a standardized error response.
func SendError(c *gin.Context, err *errors.CustomError) {
	c.JSON(err.Code, gin.H{
		"error": err.Message,
	})
}

// SendSuccess sends a success response.
func SendSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"data": data,
	})
}

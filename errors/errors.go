package errors

import "fmt"

// CustomError defines a structured error with a code and message.
type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *CustomError) Error() string {
	return fmt.Sprintf("StatusCode: %d, Message: %s", e.Code, e.Message)
}

// New creates a new CustomError.
func New(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

// Predefined errors
var (
	ErrInvalidRequestFormat = New(400, "Invalid request format")
	ErrEmptyParticipantIDs  = New(400, "Participant IDs cannot be empty")
	ErrFindingSlot          = New(500, "Error finding available slot")
	ErrNoAvailableSlot      = New(409, "No available time slot found for all participants")
	ErrCreatingMeeting      = New(500, "Error creating meeting")
)

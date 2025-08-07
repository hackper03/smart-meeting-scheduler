package handlers

import (
	"net/http"
	"smart-meeting-scheduler/errors"
	"smart-meeting-scheduler/models"
	"smart-meeting-scheduler/services"
	"smart-meeting-scheduler/utils"

	"github.com/gin-gonic/gin"
)

func ScheduleMeeting(c *gin.Context) {
	var req services.ScheduleRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, errors.ErrInvalidRequestFormat)
		return
	}

	// Validate participants
	if len(req.ParticipantIDs) == 0 {
		utils.SendError(c, errors.ErrEmptyParticipantIDs)
		return
	}

	// Find optimal slot
	slot, err := services.FindOptimalSlot(req)
	if err != nil {
		utils.SendError(c, errors.ErrFindingSlot)
		return
	}

	if slot == nil {
		utils.SendError(c, errors.ErrNoAvailableSlot)
		return
	}

	// Create meeting
	meeting, err := services.CreateMeeting(*slot, req.ParticipantIDs, "New Meeting")
	if err != nil {
		utils.SendError(c, errors.ErrCreatingMeeting)
		return
	}

	// Prepare response
	participantIDs, _ := meeting.GetParticipantIDs()
	response := models.MeetingResponse{
		MeetingID:      meeting.ID,
		Title:          meeting.Title,
		ParticipantIDs: participantIDs,
		StartTime:      meeting.StartTime.Format("2006-01-02T15:04:05Z"),
		EndTime:        meeting.EndTime.Format("2006-01-02T15:04:05Z"),
	}

	utils.SendSuccess(c, http.StatusCreated, response)
}

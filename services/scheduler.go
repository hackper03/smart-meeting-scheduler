package services

import (
	"os"
	"smart-meeting-scheduler/database"
	"smart-meeting-scheduler/logger"
	"smart-meeting-scheduler/models"
	"sort"
	"time"
)

type TimeSlot struct {
	Start time.Time
	End   time.Time
	Score float64
}

type ScheduleRequest struct {
	ParticipantIDs  []string  `json:"participantIds" binding:"required"`
	DurationMinutes int       `json:"durationMinutes" binding:"required,min=1"`
	TimeRange       TimeRange `json:"timeRange" binding:"required"`
}

type TimeRange struct {
	Start string `json:"start" binding:"required"`
	End   string `json:"end" binding:"required"`
}

func FindOptimalSlot(req ScheduleRequest) (*TimeSlot, error) {
	log := logger.NewLogger(logger.INFO, os.Stdout)

	startTime, err := time.Parse(time.RFC3339, req.TimeRange.Start)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(time.RFC3339, req.TimeRange.End)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(req.DurationMinutes) * time.Minute

	var events []models.Event
	err = database.DB.Where("user_id IN ? AND start_time < ? AND end_time > ?",
		req.ParticipantIDs, endTime, startTime).Find(&events).Error
	if err != nil {
		return nil, err
	}

	availableSlots := findAvailableSlots(startTime, endTime, duration, events)
	if len(availableSlots) == 0 {
		log.Error("No available slots found")
		return nil, nil
	}

	for i := range availableSlots {
		availableSlots[i].Score = ScoreSlot(availableSlots[i], events)
	}

	log.Info("Finding optimal slot with %d available slots", len(availableSlots))
	sort.Slice(availableSlots, func(i, j int) bool {
		if availableSlots[i].Score == availableSlots[j].Score {
			return availableSlots[i].Start.Before(availableSlots[j].Start)
		}
		return availableSlots[i].Score > availableSlots[j].Score
	})

	return &availableSlots[0], nil
}

func findAvailableSlots(start, end time.Time, duration time.Duration, events []models.Event) []TimeSlot {
	log := logger.NewLogger(logger.INFO, os.Stdout)
	log.Debug("Finding available slots between %s and %s for duration %s", start, end, duration)
	var slots []TimeSlot

	// Create a timeline of all busy periods
	type busyPeriod struct {
		start time.Time
		end   time.Time
	}

	var busyPeriods []busyPeriod
	log.Debug("Collecting busy periods from events")
	for _, event := range events {
		busyPeriods = append(busyPeriods, busyPeriod{
			start: event.StartTime,
			end:   event.EndTime,
		})
	}

	// Sort busy periods by start time
	log.Debug("Sorting busy periods")
	sort.Slice(busyPeriods, func(i, j int) bool {
		return busyPeriods[i].start.Before(busyPeriods[j].start)
	})

	// Merge overlapping periods
	log.Debug("Merging overlapping busy periods")
	var mergedPeriods []busyPeriod
	for _, period := range busyPeriods {
		if len(mergedPeriods) == 0 || mergedPeriods[len(mergedPeriods)-1].end.Before(period.start) {
			mergedPeriods = append(mergedPeriods, period)
		} else {
			// Extend the last period if there's overlap
			if period.end.After(mergedPeriods[len(mergedPeriods)-1].end) {
				mergedPeriods[len(mergedPeriods)-1].end = period.end
			}
		}
	}

	// Find gaps that can accommodate the meeting
	current := start

	for _, busy := range mergedPeriods {
		if busy.start.After(current) {
			// There's a gap before this busy period
			gapEnd := busy.start
			if gapEnd.Sub(current) >= duration {
				// This gap can accommodate the meeting
				slotEnd := current.Add(duration)
				if slotEnd.Before(gapEnd) || slotEnd.Equal(gapEnd) {
					slots = append(slots, TimeSlot{
						Start: current,
						End:   slotEnd,
					})
				}
			}
		}

		if busy.end.After(current) {
			current = busy.end
		}
	}

	// Check if there's time after the last busy period
	for end.Sub(current) >= duration {
		slotEnd := current.Add(duration)
		if slotEnd.After(end) {
			break
		}
		slots = append(slots, TimeSlot{
			Start: current,
			End:   slotEnd,
		})
		current = current.Add(duration)
	}

	return slots
}

func ScoreSlot(slot TimeSlot, events []models.Event) float64 {
	log := logger.NewLogger(logger.INFO, os.Stdout)
	log.Debug("Scoring slot from %s to %s", slot.Start, slot.End)
	score := 100.0 // Base score

	// Prefer earlier slots (higher score for earlier times)
	hour := slot.Start.Hour()
	if hour >= 9 && hour <= 12 {
		score += 20 // Morning preference
	} else if hour >= 13 && hour <= 17 {
		score += 10 // Afternoon is okay
	} else {
		score -= 30 // Penalize early morning or late hours
	}

	// Working hours bonus
	if hour >= 9 && hour < 17 {
		score += 15
	} else {
		score -= 25
	}

	// Buffer time scoring
	bufferTime := 15 * time.Minute

	for _, event := range events {
		// Check gap before slot
		gapBefore := slot.Start.Sub(event.EndTime)
		if gapBefore > 0 && gapBefore < bufferTime {
			score -= 20 // Penalize tight gaps
		} else if gapBefore > 0 && gapBefore < 30*time.Minute {
			score -= 10 // Small penalty for small gaps
		}

		// Check gap after slot
		gapAfter := event.StartTime.Sub(slot.End)
		if gapAfter > 0 && gapAfter < bufferTime {
			score -= 20
		} else if gapAfter > 0 && gapAfter < 30*time.Minute {
			score -= 10
		}

		// Bonus for back-to-back meetings (no awkward gaps)
		if gapBefore == 0 || gapAfter == 0 {
			score += 5
		}
	}

	// Prefer slots that don't fragment the day
	// Check if this slot creates small unusable gaps
	for _, event := range events {
		timeBetween := event.StartTime.Sub(slot.End)
		if timeBetween > 0 && timeBetween < 30*time.Minute {
			score -= 15 // Penalize creating small gaps
		}

		timeBefore := slot.Start.Sub(event.EndTime)
		if timeBefore > 0 && timeBefore < 30*time.Minute {
			score -= 15
		}
	}

	return score
}

func CreateMeeting(slot TimeSlot, participantIDs []string, title string) (*models.Meeting, error) {
	// Create meeting record
	meeting := models.Meeting{
		Title:     title,
		StartTime: slot.Start,
		EndTime:   slot.End,
	}

	err := meeting.SetParticipantIDs(participantIDs)
	if err != nil {
		return nil, err
	}

	err = database.DB.Create(&meeting).Error
	if err != nil {
		return nil, err
	}

	// Create calendar events for each participant
	for _, participantID := range participantIDs {
		event := models.Event{
			UserID:    participantID,
			Title:     title,
			StartTime: slot.Start,
			EndTime:   slot.End,
		}

		err = database.DB.Create(&event).Error
		if err != nil {
			return nil, err
		}
	}

	return &meeting, nil
}

package event

import (
	"learn/internal/feature/speaker"
	"learn/internal/feature/venue"
	"time"
)

type SpeakerInput struct {
	SpeakerID    uint   `json:"speaker_id" binding:"required"`
	SessionTitle string `json:"session_title"`
}

type CreateEventInput struct {
	VenueID        uint           `json:"venue_id" binding:"required"`
	Name           string         `json:"name" binding:"required"`
	Description    string         `json:"description"`
	Date           time.Time      `json:"date" binding:"required"`
	Time           time.Time      `json:"time" binding:"required"`
	Status         EventStatus    `json:"status,omitempty"`
	SalesStartDate time.Time      `json:"sales_start_date,omitempty"`
	SalesEndDate   time.Time      `json:"sales_end_date,omitempty"`
	Speakers       []SpeakerInput `json:"speakers"`
}

type EventSpeakerResponse struct {
	Speaker      speaker.Speaker `json:"speaker"`
	SessionTitle string          `json:"session_title"`
}

type EventResponse struct {
	ID             uint                   `json:"id"`
	Venue          venue.Venue            `json:"venue"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Date           time.Time              `json:"date"`
	Time           time.Time              `json:"time"`
	Status         EventStatus            `json:"status"`
	SalesStartDate time.Time              `json:"sales_start_date"`
	SalesEndDate   time.Time              `json:"sales_end_date"`
	EventSpeakers  []EventSpeakerResponse `json:"speakers"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

func ToEventResponse(event Event) EventResponse {
	var eventSpeakerResponses []EventSpeakerResponse
	for _, es := range event.EventSpeakers {
		eventSpeakerResponses = append(eventSpeakerResponses, EventSpeakerResponse{
			Speaker:      es.Speaker,
			SessionTitle: es.SessionTitle,
		})
	}

	return EventResponse{
		ID:             event.ID,
		Venue:          event.Venue,
		Name:           event.Name,
		Description:    event.Description,
		Date:           event.Date,
		Time:           event.Time,
		Status:         event.Status,
		SalesStartDate: event.SalesStartDate,
		SalesEndDate:   event.SalesEndDate,
		EventSpeakers:  eventSpeakerResponses,
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
	}
}

func ToEventResponses(events []Event) []EventResponse {
	var responses []EventResponse
	for _, event := range events {
		responses = append(responses, ToEventResponse(event))
	}
	return responses
}
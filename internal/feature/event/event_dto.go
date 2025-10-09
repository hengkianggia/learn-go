package event

import (
	"learn/internal/feature/venue"
	"time"
)

type CreateEventInput struct {
	VenueID        uint        `json:"venue_id" binding:"required"`
	Name           string      `json:"name" binding:"required"`
	Description    string      `json:"description"`
	Date           time.Time   `json:"date" binding:"required"`
	Time           time.Time   `json:"time" binding:"required"`
	Status         EventStatus `json:"status,omitempty"`
	SalesStartDate time.Time   `json:"sales_start_date,omitempty"`
	SalesEndDate   time.Time   `json:"sales_end_date,omitempty"`
}

type EventResponse struct {
	ID             uint        `json:"id"`
	Venue          venue.Venue `json:"venue"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Date           time.Time   `json:"date"`
	Time           time.Time   `json:"time"`
	Status         EventStatus `json:"status"`
	SalesStartDate time.Time   `json:"sales_start_date"`
	SalesEndDate   time.Time   `json:"sales_end_date"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

func ToEventResponse(event Event) EventResponse {
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

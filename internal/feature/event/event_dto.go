package event

import "time"

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

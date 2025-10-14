package event

import (
	"learn/internal/feature/guest"
	"learn/internal/feature/venue"
	"time"
)

type GuestInput struct {
	GuestID      uint   `json:"guest_id" binding:"required"`
	SessionTitle string `json:"session_title"`
}

type PriceInput struct {
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"required"`
}

type CreateEventInput struct {
	VenueID        uint         `json:"venue_id" binding:"required"`
	Name           string       `json:"name" binding:"required"`
	Description    string       `json:"description"`
	Date           time.Time    `json:"date" binding:"required"`
	Time           time.Time    `json:"time" binding:"required"`
	Status         EventStatus  `json:"status,omitempty"`
	SalesStartDate time.Time    `json:"sales_start_date,omitempty"`
	SalesEndDate   time.Time    `json:"sales_end_date,omitempty"`
	Guests         []GuestInput `json:"guests"`
	Prices         []PriceInput `json:"prices"`
}

type EventGuestResponse struct {
	Guest        guest.GuestResponse `json:"guest"`
	SessionTitle string              `json:"session_title"`
}

type EventPriceResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type EventResponse struct {
	ID             uint                 `json:"id"`
	Slug           string               `json:"slug"`
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	Date           time.Time            `json:"date"`
	Time           time.Time            `json:"time"`
	Status         EventStatus          `json:"status"`
	SalesStartDate time.Time            `json:"sales_start_date"`
	SalesEndDate   time.Time            `json:"sales_end_date"`
	Venue          venue.VenueResponse  `json:"venue"`
	EventGuests    []EventGuestResponse `json:"guests"`
	Prices         []EventPriceResponse `json:"prices"`
}

func ToEventPriceResponse(price EventPrice) EventPriceResponse {
	return EventPriceResponse{
		ID:    price.ID,
		Name:  price.Name,
		Price: price.Price,
	}
}

func ToEventResponse(event Event) EventResponse {
	var eventGuestResponses []EventGuestResponse
	for _, eg := range event.EventGuests {
		eventGuestResponses = append(eventGuestResponses, EventGuestResponse{
			Guest:        guest.ToGuestResponse(eg.Guest),
			SessionTitle: eg.SessionTitle,
		})
	}

	var eventPriceResponses []EventPriceResponse
	for _, p := range event.Prices {
		eventPriceResponses = append(eventPriceResponses, ToEventPriceResponse(p))
	}

	return EventResponse{
		ID:             event.ID,
		Slug:           event.Slug,
		Name:           event.Name,
		Description:    event.Description,
		Date:           event.Date,
		Time:           event.Time,
		Status:         event.Status,
		SalesStartDate: event.SalesStartDate,
		SalesEndDate:   event.SalesEndDate,
		Venue:          venue.ToVenueResponse(event.Venue),
		EventGuests:    eventGuestResponses,
		Prices:         eventPriceResponses,
	}
}

func ToEventResponses(events []Event) []EventResponse {
	var responses []EventResponse
	for _, event := range events {
		responses = append(responses, ToEventResponse(event))
	}
	return responses
}
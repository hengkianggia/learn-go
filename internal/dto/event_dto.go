package dto

import (
	"learn/internal/model"
	"time"
)

type GuestInput struct {
	GuestID      uint   `json:"guest_id" binding:"required"`
	SessionTitle string `json:"session_title"`
}

type PriceInput struct {
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"required"`
	Quota int    `json:"quota" binding:"required"`
}

type CreateEventInput struct {
	VenueID        uint              `json:"venue_id" binding:"required"`
	Name           string            `json:"name" binding:"required"`
	Description    string            `json:"description"`
	EventStartAt   time.Time         `json:"event_start_at" binding:"required"`
	Status         model.EventStatus `json:"status,omitempty"`
	SalesStartDate time.Time         `json:"sales_start_date,omitempty"`
	SalesEndDate   time.Time         `json:"sales_end_date,omitempty"`
	Guests         []GuestInput      `json:"guests"`
	Prices         []PriceInput      `json:"prices"`
}

type UpdateEventInput struct {
	Name           *string            `json:"name,omitempty"`
	Description    *string            `json:"description,omitempty"`
	EventStartAt   *time.Time         `json:"event_start_at,omitempty"`
	Time           *time.Time         `json:"time,omitempty"`
	Status         *model.EventStatus `json:"status,omitempty"`
	SalesStartDate *time.Time         `json:"sales_start_date,omitempty"`
	SalesEndDate   *time.Time         `json:"sales_end_date,omitempty"`
}

type EventGuestResponse struct {
	Guest        GuestResponse `json:"guest"`
	SessionTitle string        `json:"session_title"`
}

type EventPriceResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Quota int    `json:"quota"`
}

type EventResponseBase struct {
	ID             uint                 `json:"id"`
	Slug           string               `json:"slug"`
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	EventStartAt   time.Time            `json:"event_start_at"`
	Status         model.EventStatus    `json:"status"`
	SalesStartDate time.Time            `json:"sales_start_date"`
	SalesEndDate   time.Time            `json:"sales_end_date"`
	EventGuests    []EventGuestResponse `json:"guests"`
	Prices         []EventPriceResponse `json:"prices"`
}

type EventResponse struct {
	EventResponseBase
	Venue VenueResponse `json:"venue"`
}

type EventResponseByVenue struct {
	EventResponseBase
}

func ToEventPriceResponse(price model.EventPrice) EventPriceResponse {
	return EventPriceResponse{
		ID:    price.ID,
		Name:  price.Name,
		Price: int(price.Price),
		Quota: price.Quota,
	}
}

func toEventResponseBase(event model.Event) EventResponseBase {
	var eventGuestResponses []EventGuestResponse
	for _, eg := range event.EventGuests {
		eventGuestResponses = append(eventGuestResponses, EventGuestResponse{
			Guest:        ToGuestResponse(eg.Guest),
			SessionTitle: eg.SessionTitle,
		})
	}

	var eventPriceResponses []EventPriceResponse
	for _, p := range event.Prices {
		eventPriceResponses = append(eventPriceResponses, ToEventPriceResponse(p))
	}

	return EventResponseBase{
		ID:             event.ID,
		Slug:           event.Slug,
		Name:           event.Name,
		Description:    event.Description,
		EventStartAt:   event.EventStartAt,
		Status:         event.Status,
		SalesStartDate: event.SalesStartDate,
		SalesEndDate:   event.SalesEndDate,
		EventGuests:    eventGuestResponses,
		Prices:         eventPriceResponses,
	}
}

func ToEventResponse(event model.Event) EventResponse {
	return EventResponse{
		EventResponseBase: toEventResponseBase(event),
		Venue:             ToVenueResponse(event.Venue),
	}
}

func ToEventResponses(events []model.Event) []EventResponse {
	var responses []EventResponse
	for _, event := range events {
		responses = append(responses, ToEventResponse(event))
	}
	return responses
}

func ToEventResponseByVenue(event model.Event) EventResponseByVenue {
	return EventResponseByVenue{
		EventResponseBase: toEventResponseBase(event),
	}
}

func ToEventResponsesByVenue(events []model.Event) []EventResponseByVenue {
	var responses []EventResponseByVenue
	for _, event := range events {
		responses = append(responses, ToEventResponseByVenue(event))
	}
	return responses
}

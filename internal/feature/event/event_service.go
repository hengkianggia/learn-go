package event

import (
	"errors"
	"learn/internal/feature/venue"
	"log/slog"
)

type EventService interface {
	CreateEvent(input CreateEventInput) (*Event, error)
}

type eventService struct {
	eventRepo EventRepository
	venueRepo venue.VenueRepository
	logger    *slog.Logger
}

func NewEventService(eventRepo EventRepository, venueRepo venue.VenueRepository, logger *slog.Logger) EventService {
	return &eventService{eventRepo: eventRepo, venueRepo: venueRepo, logger: logger}
}

func (s *eventService) CreateEvent(input CreateEventInput) (*Event, error) {
	_, err := s.venueRepo.GetVenueByID(input.VenueID)
	if err != nil {
		return nil, errors.New("venue not found")
	}

	event := Event{
		VenueID:       input.VenueID,
		Name:          input.Name,
		Description:   input.Description,
		Date:          input.Date,
		Time:          input.Time,
		Status:        input.Status,
		SalesStartDate: input.SalesStartDate,
		SalesEndDate:  input.SalesEndDate,
	}

	if event.Status == "" {
		event.Status = Draft
	}

	err = s.eventRepo.CreateEvent(&event)
	if err != nil {
		s.logger.Error("failed to create event", slog.String("error", err.Error()))
		return nil, err
	}

	// After creating the event, fetch the venue data and populate the Venue field
	createdEvent, err := s.eventRepo.GetEventByID(event.ID)
	if err != nil {
		s.logger.Error("failed to get created event", slog.String("error", err.Error()))
		return nil, err
	}

	return createdEvent, nil
}
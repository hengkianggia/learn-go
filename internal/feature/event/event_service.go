package event

import (
	"errors"
	"learn/internal/feature/guest"
	"learn/internal/feature/venue"
	"log/slog"
)

type EventService interface {
	CreateEvent(input CreateEventInput) (*Event, error)
}

type eventService struct {
	eventRepo EventRepository
	venueRepo venue.VenueRepository
	guestRepo guest.GuestRepository
	logger    *slog.Logger
}

func NewEventService(eventRepo EventRepository, venueRepo venue.VenueRepository, guestRepo guest.GuestRepository, logger *slog.Logger) EventService {
	return &eventService{eventRepo: eventRepo, venueRepo: venueRepo, guestRepo: guestRepo, logger: logger}
}

func (s *eventService) CreateEvent(input CreateEventInput) (*Event, error) {
	// Check if venue exists
	_, err := s.venueRepo.GetVenueByID(input.VenueID)
	if err != nil {
		return nil, errors.New("venue not found")
	}

	// Create the event first
	event := Event{
		VenueID:        input.VenueID,
		Name:           input.Name,
		Description:    input.Description,
		Date:           input.Date,
		Time:           input.Time,
		Status:         input.Status,
		SalesStartDate: input.SalesStartDate,
		SalesEndDate:   input.SalesEndDate,
	}

	if event.Status == "" {
		event.Status = Draft
	}

	err = s.eventRepo.CreateEvent(&event)
	if err != nil {
		s.logger.Error("failed to create event", slog.String("error", err.Error()))
		return nil, err
	}

	// Now, handle the guests
	if len(input.Guests) > 0 {
		var eventGuests []EventGuest
		for _, guestInput := range input.Guests {
			// Check if guest exists
			_, err := s.guestRepo.GetGuestByID(guestInput.GuestID)
			if err != nil {
				return nil, errors.New("one or more guests not found")
			}

			eventGuests = append(eventGuests, EventGuest{
				EventID:      event.ID,
				GuestID:    guestInput.GuestID,
				SessionTitle: guestInput.SessionTitle,
			})
		}

		// Create the associations in the join table
		if err := s.eventRepo.CreateEventGuests(eventGuests); err != nil {
			s.logger.Error("failed to create event guests", slog.String("error", err.Error()))
			return nil, err
		}
	}

	// Fetch the created event with all associations
	createdEvent, err := s.eventRepo.GetEventByID(event.ID)
	if err != nil {
		s.logger.Error("failed to get created event", slog.String("error", err.Error()))
		return nil, err
	}

	return createdEvent, nil
}
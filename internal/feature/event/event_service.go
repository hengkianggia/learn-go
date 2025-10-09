package event

import (
	"errors"
	"learn/internal/feature/speaker"
	"learn/internal/feature/venue"
	"log/slog"
)

type EventService interface {
	CreateEvent(input CreateEventInput) (*Event, error)
}

type eventService struct {
	eventRepo   EventRepository
	venueRepo   venue.VenueRepository
	speakerRepo speaker.SpeakerRepository
	logger      *slog.Logger
}

func NewEventService(eventRepo EventRepository, venueRepo venue.VenueRepository, speakerRepo speaker.SpeakerRepository, logger *slog.Logger) EventService {
	return &eventService{eventRepo: eventRepo, venueRepo: venueRepo, speakerRepo: speakerRepo, logger: logger}
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

	// Now, handle the speakers
	if len(input.Speakers) > 0 {
		var eventSpeakers []EventSpeaker
		for _, speakerInput := range input.Speakers {
			// Check if speaker exists
			_, err := s.speakerRepo.GetSpeakerByID(speakerInput.SpeakerID)
			if err != nil {
				return nil, errors.New("one or more speakers not found")
			}

			eventSpeakers = append(eventSpeakers, EventSpeaker{
				EventID:      event.ID,
				SpeakerID:    speakerInput.SpeakerID,
				SessionTitle: speakerInput.SessionTitle,
			})
		}

		// Create the associations in the join table
		if err := s.eventRepo.CreateEventSpeakers(eventSpeakers); err != nil {
			s.logger.Error("failed to create event speakers", slog.String("error", err.Error()))
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

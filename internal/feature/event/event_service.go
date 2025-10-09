package event

import "log/slog"

type EventService interface {
	CreateEvent(input CreateEventInput) (*Event, error)
}

type eventService struct {
	eventRepo EventRepository
	logger    *slog.Logger
}

func NewEventService(eventRepo EventRepository, logger *slog.Logger) EventService {
	return &eventService{eventRepo: eventRepo, logger: logger}
}

func (s *eventService) CreateEvent(input CreateEventInput) (*Event, error) {
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

	err := s.eventRepo.CreateEvent(&event)
	if err != nil {
		s.logger.Error("failed to create event", slog.String("error", err.Error()))
		return nil, err
	}

	return &event, nil
}

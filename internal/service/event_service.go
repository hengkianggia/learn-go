package service

import (
	"errors"
	"fmt"
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/slug"
	"learn/internal/repository"
	"log/slog"

	"gorm.io/gorm"
)

type EventService interface {
	CreateEvent(input dto.CreateEventInput) (*model.Event, error)
	GetEventBySlug(slug string) (*model.Event, error)
	GetEventsByGuestSlug(guestSlug string) ([]model.Event, error)
}

type eventService struct {
	eventRepo repository.EventRepository
	venueRepo repository.VenueRepository
	guestRepo repository.GuestRepository
	logger    *slog.Logger
}

func NewEventService(eventRepo repository.EventRepository, venueRepo repository.VenueRepository, guestRepo repository.GuestRepository, logger *slog.Logger) EventService {
	return &eventService{eventRepo: eventRepo, venueRepo: venueRepo, guestRepo: guestRepo, logger: logger}
}

func (s *eventService) CreateEvent(input dto.CreateEventInput) (*model.Event, error) {
	// Check if venue exists
	_, err := s.venueRepo.GetVenueByID(input.VenueID)
	if err != nil {
		return nil, errors.New("venue not found")
	}

	baseSlug := slug.GenerateSlug(input.Name)
	uniqueSlug := baseSlug
	count := 1

	for {
		_, err := s.eventRepo.FindBySlug(uniqueSlug)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				break
			}
			s.logger.Error("failed to check for existing slug", slog.String("error", err.Error()))
			return nil, err
		}
		uniqueSlug = fmt.Sprintf("%s-%d", baseSlug, count)
		count++
	}

	// Create the event first
	event := model.Event{
		VenueID:        input.VenueID,
		Name:           input.Name,
		Slug:           uniqueSlug,
		Description:    input.Description,
		Date:           input.Date,
		Time:           input.Time,
		Status:         input.Status,
		SalesStartDate: input.SalesStartDate,
		SalesEndDate:   input.SalesEndDate,
	}

	if event.Status == "" {
		event.Status = model.Draft
	}

	err = s.eventRepo.CreateEvent(&event)
	if err != nil {
		s.logger.Error("failed to create event", slog.String("error", err.Error()))
		return nil, err
	}

	// Now, handle the prices
	if len(input.Prices) > 0 {
		var eventPrices []model.EventPrice
		for _, priceInput := range input.Prices {
			eventPrices = append(eventPrices, model.EventPrice{
				EventID: event.ID,
				Name:    priceInput.Name,
				Price:   priceInput.Price,
			})
		}

		// Create the prices
		if err := s.eventRepo.CreateEventPrices(eventPrices); err != nil {
			s.logger.Error("failed to create event prices", slog.String("error", err.Error()))
			return nil, err
		}
	}

	// Now, handle the guests
	if len(input.Guests) > 0 {
		var eventGuests []model.EventGuest
		for _, guestInput := range input.Guests {
			// Check if guest exists
			_, err := s.guestRepo.GetGuestByID(guestInput.GuestID)
			if err != nil {
				return nil, errors.New("one or more guests not found")
			}

			eventGuests = append(eventGuests, model.EventGuest{
				EventID:      event.ID,
				GuestID:      guestInput.GuestID,
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

func (s *eventService) GetEventBySlug(slug string) (*model.Event, error) {
	return s.eventRepo.FindBySlug(slug)
}

func (s *eventService) GetEventsByGuestSlug(guestSlug string) ([]model.Event, error) {
	return s.eventRepo.GetEventsByGuestSlug(guestSlug)
}

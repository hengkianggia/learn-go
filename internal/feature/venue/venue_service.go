package venue

import "log/slog"

type VenueService interface {
	CreateVenue(input CreateVenueInput) (*Venue, error)
}

type venueService struct {
	venueRepo VenueRepository
	logger    *slog.Logger
}

func NewVenueService(venueRepo VenueRepository, logger *slog.Logger) VenueService {
	return &venueService{venueRepo: venueRepo, logger: logger}
}

func (s *venueService) CreateVenue(input CreateVenueInput) (*Venue, error) {
	venue := Venue{
		Name:      input.Name,
		Address:   input.Address,
		City:      input.City,
		State:     input.State,
		ZipCode:   input.ZipCode,
		Capacity:  input.Capacity,
		IsActive:  input.IsActive,
	}

	err := s.venueRepo.CreateVenue(&venue)
	if err != nil {
		s.logger.Error("failed to create venue", slog.String("error", err.Error()))
		return nil, err
	}

	return &venue, nil
}

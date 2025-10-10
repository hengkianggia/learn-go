package venue

import (
	"errors"
	"fmt"
	"learn/internal/pkg/slug"
	"log/slog"

	"gorm.io/gorm"
)

type VenueService interface {
	CreateVenue(input CreateVenueInput) (*Venue, error)
	GetVenueBySlug(slug string) (*Venue, error)
}

type venueService struct {
	venueRepo VenueRepository
	logger    *slog.Logger
}

func NewVenueService(venueRepo VenueRepository, logger *slog.Logger) VenueService {
	return &venueService{venueRepo: venueRepo, logger: logger}
}

func (s *venueService) CreateVenue(input CreateVenueInput) (*Venue, error) {
	baseSlug := slug.GenerateSlug(input.Name)
	uniqueSlug := baseSlug
	count := 1

	for {
		_, err := s.venueRepo.FindBySlug(uniqueSlug)
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

	venue := Venue{
		Name:      input.Name,
		Slug:      uniqueSlug,
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

func (s *venueService) GetVenueBySlug(slug string) (*Venue, error) {
	return s.venueRepo.FindBySlug(slug)
}
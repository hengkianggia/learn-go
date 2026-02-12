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

type venueService struct {
	venueRepo repository.VenueRepository
	logger    *slog.Logger
}

type VenueService interface {
	CreateVenue(input dto.CreateVenueInput) (*model.Venue, error)
	GetVenueBySlug(slug string) (*model.Venue, error)
	UpdateVenue(slug string, input dto.UpdateVenueInput) (*model.Venue, error)
}

func NewVenueService(venueRepo repository.VenueRepository, logger *slog.Logger) VenueService {
	return &venueService{venueRepo: venueRepo, logger: logger}
}

func (s *venueService) CreateVenue(input dto.CreateVenueInput) (*model.Venue, error) {
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

	venue := model.Venue{
		Name:     input.Name,
		Slug:     uniqueSlug,
		Address:  input.Address,
		City:     input.City,
		State:    input.State,
		ZipCode:  input.ZipCode,
		Capacity: input.Capacity,
		IsActive: input.IsActive,
		Country:  input.Country,
	}

	err := s.venueRepo.CreateVenue(&venue)
	if err != nil {
		s.logger.Error("failed to create venue", slog.String("error", err.Error()))
		return nil, err
	}

	return &venue, nil
}

func (s *venueService) GetVenueBySlug(slug string) (*model.Venue, error) {
	return s.venueRepo.FindBySlug(slug)
}

func (s *venueService) UpdateVenue(slug string, input dto.UpdateVenueInput) (*model.Venue, error) {
	venue, err := s.venueRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		venue.Name = *input.Name
	}
	if input.Address != nil {
		venue.Address = *input.Address
	}
	if input.City != nil {
		venue.City = *input.City
	}
	if input.State != nil {
		venue.State = *input.State
	}
	if input.ZipCode != nil {
		venue.ZipCode = *input.ZipCode
	}
	if input.Capacity != nil {
		venue.Capacity = *input.Capacity
	}
	if input.IsActive != nil {
		venue.IsActive = *input.IsActive
	}
	if input.Country != nil {
		venue.Country = *input.Country
	}

	if err := s.venueRepo.UpdateVenue(venue); err != nil {
		s.logger.Error("failed to update venue", slog.String("error", err.Error()))
		return nil, err
	}

	return venue, nil
}

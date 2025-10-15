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

type GuestService interface {
	CreateGuest(input dto.CreateGuestInput) (*model.Guest, error)
	GetGuestBySlug(slug string) (*model.Guest, error)
	UpdateGuest(slug string, input dto.UpdateGuestInput) (*model.Guest, error)
}

type guestService struct {
	guestRepo repository.GuestRepository
	logger    *slog.Logger
}

func NewGuestService(guestRepo repository.GuestRepository, logger *slog.Logger) GuestService {
	return &guestService{guestRepo: guestRepo, logger: logger}
}

func (s *guestService) CreateGuest(input dto.CreateGuestInput) (*model.Guest, error) {
	baseSlug := slug.GenerateSlug(input.Name)
	uniqueSlug := baseSlug
	count := 1

	for {
		_, err := s.guestRepo.FindBySlug(uniqueSlug)
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

	guest := model.Guest{
		Name: input.Name,
		Slug: uniqueSlug,
		Bio:  input.Bio,
	}

	err := s.guestRepo.CreateGuest(&guest)
	if err != nil {
		s.logger.Error("failed to create guest", slog.String("error", err.Error()))
		return nil, err
	}

	return &guest, nil
}

func (s *guestService) GetGuestBySlug(slug string) (*model.Guest, error) {
	return s.guestRepo.FindBySlug(slug)
}

func (s *guestService) UpdateGuest(slug string, input dto.UpdateGuestInput) (*model.Guest, error) {
	guest, err := s.guestRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		guest.Name = *input.Name
	}
	if input.Bio != nil {
		guest.Bio = *input.Bio
	}

	if err := s.guestRepo.UpdateGuest(guest); err != nil {
		s.logger.Error("failed to update guest", slog.String("error", err.Error()))
		return nil, err
	}

	return guest, nil
}

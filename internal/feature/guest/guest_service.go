package guest

import (
	"errors"
	"fmt"
	"learn/internal/pkg/slug"
	"log/slog"

	"gorm.io/gorm"
)

type GuestService interface {
	CreateGuest(input CreateGuestInput) (*Guest, error)
	GetGuestBySlug(slug string) (*Guest, error)
}

type guestService struct {
	guestRepo GuestRepository
	logger    *slog.Logger
}

func NewGuestService(guestRepo GuestRepository, logger *slog.Logger) GuestService {
	return &guestService{guestRepo: guestRepo, logger: logger}
}

func (s *guestService) CreateGuest(input CreateGuestInput) (*Guest, error) {
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

	guest := Guest{
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

func (s *guestService) GetGuestBySlug(slug string) (*Guest, error) {
	return s.guestRepo.FindBySlug(slug)
}
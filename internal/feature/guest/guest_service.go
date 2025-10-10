package guest

import "log/slog"

type GuestService interface {
	CreateGuest(input CreateGuestInput) (*Guest, error)
}

type guestService struct {
	guestRepo GuestRepository
	logger    *slog.Logger
}

func NewGuestService(guestRepo GuestRepository, logger *slog.Logger) GuestService {
	return &guestService{guestRepo: guestRepo, logger: logger}
}

func (s *guestService) CreateGuest(input CreateGuestInput) (*Guest, error) {
	guest := Guest{
		Name: input.Name,
		Bio:  input.Bio,
	}

	err := s.guestRepo.CreateGuest(&guest)
	if err != nil {
		s.logger.Error("failed to create guest", slog.String("error", err.Error()))
		return nil, err
	}

	return &guest, nil
}
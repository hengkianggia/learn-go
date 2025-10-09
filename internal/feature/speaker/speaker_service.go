package speaker

import "log/slog"

type SpeakerService interface {
	CreateSpeaker(input CreateSpeakerInput) (*Speaker, error)
}

type speakerService struct {
	speakerRepo SpeakerRepository
	logger      *slog.Logger
}

func NewSpeakerService(speakerRepo SpeakerRepository, logger *slog.Logger) SpeakerService {
	return &speakerService{speakerRepo: speakerRepo, logger: logger}
}

func (s *speakerService) CreateSpeaker(input CreateSpeakerInput) (*Speaker, error) {
	speaker := Speaker{
		Name: input.Name,
		Bio:  input.Bio,
	}

	err := s.speakerRepo.CreateSpeaker(&speaker)
	if err != nil {
		s.logger.Error("failed to create speaker", slog.String("error", err.Error()))
		return nil, err
	}

	return &speaker, nil
}

package service

import (
	"errors"
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/repository"
	"log/slog"
)

type AuthService interface {
	Register(input dto.RegisterInput) (*model.User, error)
	Login(input dto.LoginInput) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewAuthService(userRepo repository.UserRepository, logger *slog.Logger) AuthService {
	return &authService{userRepo: userRepo, logger: logger}
}

func (s *authService) Register(input dto.RegisterInput) (*model.User, error) {
	// Validate that passwords match
	if input.Password != input.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	user := model.User{
		Name:        input.Name,
		Email:       input.Email,
		Password:    input.Password, // Password will be hashed by the BeforeSave hook
		PhoneNumber: input.PhoneNumber,
		UserType:    input.UserType,
		IsVerified:  false,
	}

	if user.UserType == "" {
		user.UserType = model.Attendee
	}

	err := s.userRepo.Save(&user)
	if err != nil {
		s.logger.Error("failed to save user", slog.String("error", err.Error()))
		return nil, err
	}

	return &user, nil
}

func (s *authService) Login(input dto.LoginInput) (string, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := ValidatePassword(user.Password, input.Password); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := GenerateJWT(*user)
	if err != nil {
		s.logger.Error("failed to generate token", slog.String("error", err.Error()))
		return "", errors.New("could not generate token")
	}

	return token, nil
}

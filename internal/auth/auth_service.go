package auth

import (
	"errors"
	"log/slog"
)

type AuthService interface {
	Register(input RegisterInput) (*User, error)
	Login(input LoginInput) (string, error)
}

type authService struct {
	userRepo UserRepository
	logger   *slog.Logger
}

func NewAuthService(userRepo UserRepository, logger *slog.Logger) AuthService {
	return &authService{userRepo: userRepo, logger: logger}
}

func (s *authService) Register(input RegisterInput) (*User, error) {
	// Validate that passwords match
	if input.Password != input.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	user := User{
		Name:        input.Name,
		Email:       input.Email,
		Password:    input.Password, // Password will be hashed by the BeforeSave hook
		PhoneNumber: input.PhoneNumber,
		UserType:    input.UserType,
		IsVerified:  false,
	}

	if user.UserType == "" {
		user.UserType = Attendee
	}

	err := s.userRepo.Save(&user)
	if err != nil {
		s.logger.Error("failed to save user", slog.String("error", err.Error()))
		return nil, err
	}

	return &user, nil
}

func (s *authService) Login(input LoginInput) (string, error) {
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
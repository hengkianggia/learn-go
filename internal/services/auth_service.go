package services

import (
	"errors"
	"learn/internal/auth"
	"learn/internal/dto"
	"learn/internal/models"
	"learn/internal/repositories"
)

// AuthService

type AuthService interface {
	Register(input dto.RegisterInput) (*models.User, error)
	Login(input dto.LoginInput) (string, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(input dto.RegisterInput) (*models.User, error) {
	user := models.User{
		Username: input.Username,
		Password: input.Password,
	}

	err := s.userRepo.Save(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *authService) Login(input dto.LoginInput) (string, error) {
	user, err := s.userRepo.FindByUsername(input.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := auth.ValidatePassword(user.Password, input.Password); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateJWT(*user)
	if err != nil {
		return "", errors.New("could not generate token")
	}

	return token, nil
}

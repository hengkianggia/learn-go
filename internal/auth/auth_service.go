package auth

import (
	"errors"
)

// AuthService

type AuthService interface {
	Register(input RegisterInput) (*User, error)
	Login(input LoginInput) (string, error)
}

type authService struct {
	userRepo UserRepository
}

func NewAuthService(userRepo UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(input RegisterInput) (*User, error) {
	user := User{
		Username: input.Username,
		Password: input.Password,
	}

	err := s.userRepo.Save(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *authService) Login(input LoginInput) (string, error) {
	user, err := s.userRepo.FindByUsername(input.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := ValidatePassword(user.Password, input.Password); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := GenerateJWT(*user)
	if err != nil {
		return "", errors.New("could not generate token")
	}

	return token, nil
}
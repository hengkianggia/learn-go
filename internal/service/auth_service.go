package service

import (
	"context"
	"errors"
	"fmt"
	"learn/internal/config"
	"learn/internal/dto"
	"learn/internal/model"
	"learn/internal/pkg/random"
	"learn/internal/repository"
	"log/slog"
	"time"
)

type authService struct {
	userRepo     repository.UserRepository
	emailService EmailService
	logger       *slog.Logger
}

type AuthService interface {
	Register(input dto.RegisterInput) (*model.User, error)
	Login(input dto.LoginInput) (string, error)
	VerifyOTP(email string, otp string) error
}

func NewAuthService(userRepo repository.UserRepository, emailService EmailService, logger *slog.Logger) AuthService {
	return &authService{
		userRepo:     userRepo,
		emailService: emailService,
		logger:       logger,
	}
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

	// Generate OTP
	otp := random.StringWithCharset(6, "0123456789")

	// Store OTP in Redis
	otpKey := fmt.Sprintf("auth:otp:%s", user.Email)
	err = config.Rdb.Set(context.Background(), otpKey, otp, 5*time.Minute).Err()
	if err != nil {
		s.logger.Error("failed to save OTP to Redis", slog.String("error", err.Error()))
		// Note: User is already created. In a production system, we might want to handle this better (e.g. rollback or separate step).
		return nil, errors.New("failed to generate verification code")
	}

	// Send OTP Email
	err = s.emailService.SendOTP(user.Email, otp)
	if err != nil {
		s.logger.Error("failed to send OTP email", slog.String("error", err.Error()))
		// Proceeding even if email fails, as user can request resend (if implemented) or we can just log it.
		// For this task, we'll return the error to let the user know something went wrong.
		return nil, errors.New("failed to send verification email")
	}

	return &user, nil
}

func (s *authService) VerifyOTP(email string, otp string) error {
	otpKey := fmt.Sprintf("auth:otp:%s", email)
	storedOTP, err := config.Rdb.Get(context.Background(), otpKey).Result()
	if err != nil {
		return errors.New("invalid or expired OTP")
	}

	if storedOTP != otp {
		return errors.New("invalid OTP")
	}

	// OTP is valid, verify user
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	user.IsVerified = true
	err = s.userRepo.Save(user)
	if err != nil {
		return errors.New("failed to update user verification status")
	}

	// Delete OTP from Redis
	config.Rdb.Del(context.Background(), otpKey)

	return nil
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

package service

import (
	"errors"
	"learn/internal/model"
	"learn/internal/repository"
	"log/slog"
)

type adminService struct {
	userRepo     repository.UserRepository
	emailService EmailService
	logger       *slog.Logger
}

type AdminService interface {
	ApproveUser(userID uint) (*model.User, error)
	RejectUser(userID uint) error
	BlockUser(userID uint) (*model.User, error)
	UnblockUser(userID uint) (*model.User, error)
	DeleteUser(userID uint) error
}

func NewAdminService(userRepo repository.UserRepository, emailService EmailService, logger *slog.Logger) AdminService {
	return &adminService{
		userRepo:     userRepo,
		emailService: emailService,
		logger:       logger,
	}
}

func (s *adminService) ApproveUser(userID uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.UserType != model.Organizer {
		return nil, errors.New("only organizer accounts can be approved")
	}

	if user.IsApproved {
		return nil, errors.New("user is already approved")
	}

	err = s.userRepo.UpdateFields(userID, map[string]interface{}{"is_approved": true})
	if err != nil {
		s.logger.Error("failed to approve user", slog.String("error", err.Error()))
		return nil, errors.New("failed to approve user")
	}

	user.IsApproved = true

	s.logger.Info("organizer approved", slog.Uint64("user_id", uint64(userID)))
	return user, nil
}

func (s *adminService) RejectUser(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.UserType != model.Organizer {
		return errors.New("only organizer accounts can be rejected")
	}

	err = s.userRepo.Delete(userID)
	if err != nil {
		s.logger.Error("failed to reject user", slog.String("error", err.Error()))
		return errors.New("failed to reject user")
	}

	s.logger.Info("organizer rejected", slog.Uint64("user_id", uint64(userID)))
	return nil
}

func (s *adminService) BlockUser(userID uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.UserType == model.Administrator {
		return nil, errors.New("cannot block an administrator")
	}

	if user.IsBlocked {
		return nil, errors.New("user is already blocked")
	}

	err = s.userRepo.UpdateFields(userID, map[string]interface{}{"is_blocked": true})
	if err != nil {
		s.logger.Error("failed to block user", slog.String("error", err.Error()))
		return nil, errors.New("failed to block user")
	}

	user.IsBlocked = true

	s.logger.Info("user blocked", slog.Uint64("user_id", uint64(userID)))
	return user, nil
}

func (s *adminService) UnblockUser(userID uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsBlocked {
		return nil, errors.New("user is not blocked")
	}

	err = s.userRepo.UpdateFields(userID, map[string]interface{}{"is_blocked": false})
	if err != nil {
		s.logger.Error("failed to unblock user", slog.String("error", err.Error()))
		return nil, errors.New("failed to unblock user")
	}

	user.IsBlocked = false

	s.logger.Info("user unblocked", slog.Uint64("user_id", uint64(userID)))
	return user, nil
}

func (s *adminService) DeleteUser(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.UserType == model.Administrator {
		return errors.New("cannot delete an administrator")
	}

	err = s.userRepo.Delete(userID)
	if err != nil {
		s.logger.Error("failed to delete user", slog.String("error", err.Error()))
		return errors.New("failed to delete user")
	}

	s.logger.Info("user deleted", slog.Uint64("user_id", uint64(userID)))
	return nil
}

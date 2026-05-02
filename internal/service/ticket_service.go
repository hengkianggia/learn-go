package service

import (
	"errors"
	"learn/internal/dto"
	apperrors "learn/internal/errors"
	"learn/internal/model"
	"learn/internal/repository"
	"log/slog"
	"strings"

	"gorm.io/gorm"
)

type TicketService interface {
	CheckInTicket(input dto.CheckInTicketRequest, userID uint) (*dto.CheckInTicketResponse, error)
}

type ticketService struct {
	ticketRepo repository.TicketRepository
	logger     *slog.Logger
}

func NewTicketService(ticketRepo repository.TicketRepository, logger *slog.Logger) TicketService {
	return &ticketService{ticketRepo: ticketRepo, logger: logger}
}

func (s *ticketService) CheckInTicket(input dto.CheckInTicketRequest, userID uint) (*dto.CheckInTicketResponse, error) {
	ticketCode := strings.TrimSpace(input.TicketCode)
	if ticketCode == "" {
		return nil, apperrors.NewValidationError("ticket_code", "ticket code is required", input.TicketCode)
	}

	ticket, order, checkedIn, err := s.ticketRepo.CheckInTicketByCode(ticketCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewBusinessRuleError("ticket_exists", "ticket not found")
		}

		s.logger.Error("failed to check in ticket",
			slog.String("ticket_code", ticketCode),
			slog.Uint64("user_id", uint64(userID)),
			slog.String("error", err.Error()))
		return nil, apperrors.NewSystemError("check_in_ticket", err)
	}

	if order.Status != model.OrderPaid {
		return nil, apperrors.NewBusinessRuleError("ticket_order_status", "ticket order is not paid")
	}

	if !checkedIn {
		return nil, apperrors.NewBusinessRuleError("ticket_already_scanned", "ticket has already been checked in")
	}

	s.logger.Info("ticket checked in",
		slog.String("ticket_code", ticket.TicketCode),
		slog.Uint64("ticket_id", uint64(ticket.ID)),
		slog.Uint64("order_id", uint64(ticket.OrderID)),
		slog.Uint64("user_id", uint64(userID)))

	response := dto.ToCheckInTicketResponse(*ticket, *order)
	return &response, nil
}

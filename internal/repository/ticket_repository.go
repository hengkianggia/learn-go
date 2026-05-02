package repository

import (
	"learn/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketRepository interface {
	CreateTickets(tickets []model.Ticket) error
	CheckInTicketByCode(ticketCode string) (*model.Ticket, *model.Order, bool, error)
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) CreateTickets(tickets []model.Ticket) error {
	return r.db.Create(&tickets).Error
}

func (r *ticketRepository) CheckInTicketByCode(ticketCode string) (*model.Ticket, *model.Order, bool, error) {
	var ticket model.Ticket
	var order model.Order
	checkedIn := false

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("ticket_code = ?", ticketCode).
			First(&ticket).Error; err != nil {
			return err
		}

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&order, ticket.OrderID).Error; err != nil {
			return err
		}

		if order.Status != model.OrderPaid || ticket.IsScanned {
			return nil
		}

		if err := tx.Model(&ticket).Update("is_scanned", true).Error; err != nil {
			return err
		}

		ticket.IsScanned = true
		checkedIn = true
		return nil
	})
	if err != nil {
		return nil, nil, false, err
	}

	return &ticket, &order, checkedIn, nil
}

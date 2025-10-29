package repository

import (
	"learn/internal/model"

	"gorm.io/gorm"
)

type TicketRepository interface {
	CreateTickets(tickets []model.Ticket) error
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

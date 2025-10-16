package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "PENDING"
	OrderPaid      OrderStatus = "PAID"
	OrderCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	gorm.Model
	UserID       uint        `gorm:"not null"`
	User         User
	TotalPrice   float64     `gorm:"not null"`
	Status       OrderStatus `gorm:"not null;default:'PENDING'"`
	PaymentDue   time.Time
	Tickets      []Ticket
	Payments     []Payment
}

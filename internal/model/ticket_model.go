package model

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	EventPriceID uint   `gorm:"not null"`
	EventPrice   EventPrice
	OrderID      uint   `gorm:"not null"`
	Price        float64 `gorm:"not null"`
	Type         string `gorm:"not null"`
	SeatNumber   string
	TicketCode   string `gorm:"not null;unique"`
	IsScanned    bool   `gorm:"default:false"`
	OwnerName    string
	OwnerEmail   string
}

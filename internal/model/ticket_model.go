package model

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	EventPriceID uint   `gorm:"not null"`
	EventPrice   EventPrice
	OrderID      uint   `gorm:"not null"`
	Price        int64 `gorm:"not null"` // Price in smallest currency unit (e.g., cents)
	Type         string `gorm:"not null"`
	SeatNumber   string
	TicketCode   string `gorm:"not null;unique"`
	IsScanned    bool   `gorm:"default:false"`
	OwnerName    string
	OwnerEmail   string
}

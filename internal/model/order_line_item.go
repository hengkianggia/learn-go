package model

import "gorm.io/gorm"

type OrderLineItem struct {
	gorm.Model
	OrderID      uint  `gorm:"not null"`
	EventPriceID uint  `gorm:"not null"`
	Quantity     int   `gorm:"not null"`
	PricePerUnit int64 `gorm:"not null"` // Price per unit in smallest currency unit (e.g., cents)
	TotalPrice   int64 `gorm:"not null"` // Total price for this line item (PricePerUnit * Quantity)
}

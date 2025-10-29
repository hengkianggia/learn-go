package model

import "gorm.io/gorm"

type OrderLineItem struct {
	gorm.Model
	OrderID      uint  `gorm:"not null"`
	EventPriceID uint  `gorm:"not null"`
	Quantity     int   `gorm:"not null"`
	PricePerUnit int64 `gorm:"not null"` // Store the price at the time of order
}

package model

import (
	"time"

	"gorm.io/gorm"
)

// Payment represents the payment details for an order.
type Payment struct {
	gorm.Model
	OrderID       uint          `gorm:"unique;not null" json:"order_id"` // Foreign Key ke Order (1:1 relationship)
	PaymentMethod PaymentMethod `gorm:"type:varchar(50);not null" json:"payment_method"`
	TransactionID string        `gorm:"unique;not null" json:"transaction_id"` // ID from payment gateway
	// Amount        int64         `gorm:"not null" json:"amount"`                // Amount paid in smallest currency unit (e.g., cents)
	PaymentStatus PaymentStatus `gorm:"type:varchar(20);not null" json:"payment_status"`
	PaymentDate   time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"payment_date"`

	// Relationship
	Order Order `gorm:"foreignKey:OrderID"` // BelongsTo relationship with Order
}

package model

import (
	"time"

	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "PENDING"
	PaymentSuccess PaymentStatus = "SUCCESS"
	PaymentFailed  PaymentStatus = "FAILED"
)

type Payment struct {
	gorm.Model
	OrderID       uint          `gorm:"not null"`
	PaymentMethod string        `gorm:"not null"`
	TransactionID string        `gorm:"not null;unique"`
	Amount        float64       `gorm:"not null"`
	Status        PaymentStatus `gorm:"not null"`
	PaymentDate   time.Time
}

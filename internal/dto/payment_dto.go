package dto

import (
	"learn/internal/model"
	"time"
)

// PaymentResponse is the standardized response for a payment
type PaymentResponse struct {
	PaymentID     uint                `json:"payment_id"`
	OrderID       uint                `json:"order_id"`
	PaymentMethod model.PaymentMethod `json:"payment_method"`
	TransactionID string              `json:"transaction_id"`
	PaymentStatus model.PaymentStatus `json:"payment_status"`
	PaymentDate   time.Time           `json:"payment_date"`
}

// CreatePaymentRequest represents the request body for creating a new payment
type CreatePaymentRequest struct {
	OrderID       uint                `json:"order_id" binding:"required"`
	PaymentMethod model.PaymentMethod `json:"payment_method" binding:"required"`
	TransactionID string              `json:"transaction_id" binding:"required"`
}

// UpdatePaymentRequest represents the request body for updating an existing payment
type UpdatePaymentRequest struct {
	PaymentMethod *model.PaymentMethod `json:"payment_method"`
	TransactionID *string              `json:"transaction_id"`
	Amount        *int64               `json:"amount" binding:"min=1"`
	PaymentStatus *model.PaymentStatus `json:"payment_status"`
}

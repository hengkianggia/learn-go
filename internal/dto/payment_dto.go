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
	Amount        int64               `json:"amount"`
	PaymentStatus model.PaymentStatus `json:"payment_status"`
	PaymentDate   time.Time           `json:"payment_date"`

	// Midtrans Details
	PaymentURL           string `json:"payment_url,omitempty"`
	VirtualAccountNumber string `json:"virtual_account_number,omitempty"`
	BillKey              string `json:"bill_key,omitempty"`
	BillerCode           string `json:"biller_code,omitempty"`
	PaymentCode          string `json:"payment_code,omitempty"`
}

// CreatePaymentRequest represents the request body for creating a new payment
type CreatePaymentRequest struct {
	OrderID       uint                `json:"order_id" binding:"required"`
	PaymentMethod model.PaymentMethod `json:"payment_method" binding:"required"`
}

// UpdatePaymentRequest represents the request body for updating an existing payment
type UpdatePaymentRequest struct {
	PaymentMethod *model.PaymentMethod `json:"payment_method"`
	TransactionID *string              `json:"transaction_id"`
	Amount        *int64               `json:"amount" binding:"min=1"`
	PaymentStatus *model.PaymentStatus `json:"payment_status"`
}

// UpdatePaymentStatusRequest represents the request body for updating only the payment status
type UpdatePaymentStatusRequest struct {
	Status model.PaymentStatus `json:"status" binding:"required"`
}

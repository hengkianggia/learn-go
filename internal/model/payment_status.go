package model

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusSuccess PaymentStatus = "SUCCESS"
	PaymentStatusFailed  PaymentStatus = "FAILED"

	// Optional: Add more statuses as needed, e.g., for refunds or partial payments
	PaymentStatusRefunded PaymentStatus = "REFUNDED"
)

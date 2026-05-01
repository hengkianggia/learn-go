package model

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusSuccess PaymentStatus = "SUCCESS"
	PaymentStatusFailed  PaymentStatus = "FAILED"

	// Optional: Add more statuses as needed, e.g., for refunds or partial payments
	PaymentStatusRefunded PaymentStatus = "REFUNDED"
)

func IsTerminalPaymentStatus(status PaymentStatus) bool {
	return status == PaymentStatusSuccess || status == PaymentStatusFailed || status == PaymentStatusRefunded
}

func CanTransitionPaymentStatus(from, to PaymentStatus) bool {
	if from == to {
		return true
	}

	switch from {
	case PaymentStatusPending:
		return to == PaymentStatusSuccess || to == PaymentStatusFailed
	case PaymentStatusSuccess:
		return to == PaymentStatusRefunded
	case PaymentStatusFailed, PaymentStatusRefunded:
		return false
	default:
		return false
	}
}

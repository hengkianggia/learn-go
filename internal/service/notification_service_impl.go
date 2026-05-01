package service

import (
	"errors"
	apperrors "learn/internal/errors"
	"learn/internal/model"
	"log/slog"

	"gorm.io/gorm"
)

func (s *paymentService) HandleNotification(payload map[string]interface{}) error {
	// 1. Verify Signature
	valid, err := s.midtransGateway.VerifyPaymentNotification(payload)
	if err != nil {
		s.logger.Error("failed to verify notification signature", slog.String("error", err.Error()))
		return apperrors.NewSystemError("verify_notification", err)
	}
	if !valid {
		s.logger.Warn("invalid notification signature")
		return apperrors.NewBusinessRuleError("notification_signature", "invalid signature")
	}

	// 2. Get Order ID and Transaction Status
	orderIDStr, _ := payload["order_id"].(string)
	transactionStatus, _ := payload["transaction_status"].(string)
	fraudStatus, _ := payload["fraud_status"].(string)
	transactionID, _ := payload["transaction_id"].(string)
	if orderIDStr == "" || transactionStatus == "" || transactionID == "" {
		return apperrors.NewBusinessRuleError("notification_payload", "missing required notification fields")
	}

	s.logger.Info("processing midtrans notification",
		slog.String("order_id", orderIDStr),
		slog.String("transaction_status", transactionStatus),
		slog.String("fraud_status", fraudStatus),
	)

	// 3. Get Payment by Transaction ID
	payment, err := s.paymentRepository.GetPaymentByTransactionID(transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("payment not found for notification", slog.String("transaction_id", transactionID))
			return apperrors.NewBusinessRuleError("payment_not_found", "payment not found")
		}
		s.logger.Error("failed to get payment by transaction ID", slog.String("error", err.Error()))
		return apperrors.NewSystemError("get_payment_by_transaction_id", err)
	}

	// 4. Determine New Status
	var newStatus model.PaymentStatus

	switch transactionStatus {
	case "capture":
		if fraudStatus == "challenge" {
			newStatus = model.PaymentStatusPending // Admin need to review
		} else if fraudStatus == "accept" {
			newStatus = model.PaymentStatusSuccess
		}
	case "settlement":
		newStatus = model.PaymentStatusSuccess
	case "deny", "cancel", "expire":
		newStatus = model.PaymentStatusFailed
	case "refund", "partial_refund":
		newStatus = model.PaymentStatusRefunded
	case "pending":
		newStatus = model.PaymentStatusPending
	default:
		s.logger.Info("ignoring transaction status", slog.String("status", transactionStatus))
		return nil
	}

	// 5. Update Payment Status if changed. Duplicate notifications are intentionally idempotent.
	if payment.PaymentStatus == newStatus {
		s.logger.Info("ignoring duplicate midtrans notification",
			slog.Uint64("payment_id", uint64(payment.ID)),
			slog.String("status", string(newStatus)),
		)
		return nil
	}

	if !model.CanTransitionPaymentStatus(payment.PaymentStatus, newStatus) {
		s.logger.Warn("ignoring invalid payment status transition from notification",
			slog.Uint64("payment_id", uint64(payment.ID)),
			slog.String("from", string(payment.PaymentStatus)),
			slog.String("to", string(newStatus)),
		)
		return nil
	}

	_, err = s.UpdatePaymentStatus(payment.ID, newStatus)
	if err != nil {
		s.logger.Error("failed to update payment status from notification", slog.String("error", err.Error()))
		return err
	}

	return nil
}

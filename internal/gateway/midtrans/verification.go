package midtrans

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"learn/internal/config"
)

func (g *midtransGateway) VerifyPaymentNotification(payload map[string]interface{}) (bool, error) {
	orderID, exists := payload["order_id"].(string)
	if !exists {
		return false, errors.New("invalid notification payload: order_id missing")
	}
	statusCode, exists := payload["status_code"].(string)
	if !exists {
		return false, errors.New("invalid notification payload: status_code missing")
	}
	grossAmount, exists := payload["gross_amount"].(string)
	if !exists {
		return false, errors.New("invalid notification payload: gross_amount missing")
	}
	signatureKey, exists := payload["signature_key"].(string)
	if !exists {
		return false, errors.New("invalid notification payload: signature_key missing")
	}

	serverKey := config.AppConfig.MidtransServerKey
	input := orderID + statusCode + grossAmount + serverKey

	hash := sha512.Sum512([]byte(input))
	expectedSignature := hex.EncodeToString(hash[:])

	if signatureKey != expectedSignature {
		return false, nil
	}

	return true, nil
}

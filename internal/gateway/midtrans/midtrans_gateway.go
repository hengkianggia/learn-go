package midtrans

import (
	"errors"
	"learn/internal/config"
	"log/slog"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type midtransGateway struct {
	client coreapi.Client
	logger *slog.Logger
}

type MidtransGateway interface {
	ChargeBankTransfer(orderID string, amount int64, bank string) (*coreapi.ChargeResponse, error)
	ChargeGopay(orderID string, amount int64) (*coreapi.ChargeResponse, error)
	ChargeIndomaret(orderID string, amount int64, message string) (*coreapi.ChargeResponse, error)
	VerifyPaymentNotification(payload map[string]interface{}) (bool, error)
}

func NewMidtransGateway(logger *slog.Logger) MidtransGateway {
	c := coreapi.Client{}

	env := midtrans.Sandbox
	if config.AppConfig.MidtransEnv == "production" {
		env = midtrans.Production
	}

	c.New(config.AppConfig.MidtransServerKey, env)

	return &midtransGateway{
		client: c,
		logger: logger,
	}
}

func (g *midtransGateway) ChargeBankTransfer(orderID string, amount int64, bank string) (*coreapi.ChargeResponse, error) {
	req := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeBankTransfer,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		BankTransfer: &coreapi.BankTransferDetails{
			Bank: midtrans.Bank(bank),
		},
	}

	resp, err := g.client.ChargeTransaction(req)
	if err != nil {
		g.logger.Error("Midtrans Charge Error", slog.String("error", err.Message), slog.String("raw_error", err.Error()))
		return nil, errors.New("midtrans charge failed: " + err.Message)
	}

	return resp, nil
}

func (g *midtransGateway) ChargeGopay(orderID string, amount int64) (*coreapi.ChargeResponse, error) {
	req := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeGopay,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
	}

	resp, err := g.client.ChargeTransaction(req)
	if err != nil {
		g.logger.Error("Midtrans Charge Error", slog.String("error", err.Message))
		return nil, errors.New("midtrans charge failed: " + err.Message)
	}

	return resp, nil
}

func (g *midtransGateway) ChargeIndomaret(orderID string, amount int64, message string) (*coreapi.ChargeResponse, error) {
	req := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeConvenienceStore,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		ConvStore: &coreapi.ConvStoreDetails{
			Store:   "indomaret",
			Message: message,
		},
	}

	resp, err := g.client.ChargeTransaction(req)
	if err != nil {
		g.logger.Error("Midtrans Charge Error", slog.String("error", err.Message))
		return nil, errors.New("midtrans charge failed: " + err.Message)
	}

	return resp, nil
}

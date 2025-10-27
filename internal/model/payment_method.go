package model

type PaymentMethod string

const (
	PaymentMethodCreditCard     PaymentMethod = "CREDIT_CARD"
	PaymentMethodVirtualAccount PaymentMethod = "VIRTUAL_ACCOUNT"
	PaymentMethodPayPal         PaymentMethod = "PAYPAL"
	// Add other payment methods as needed
)

package model

type PaymentMethod string

const (
	PaymentMethodCreditCard      PaymentMethod = "CREDIT_CARD"
	PaymentMethodVirtualAccount  PaymentMethod = "VIRTUAL_ACCOUNT"
	PaymentMethodPayPal          PaymentMethod = "PAYPAL"
	PaymentMethodBankTransferBCA PaymentMethod = "BANK_TRANSFER_BCA"
	PaymentMethodBankTransferBNI PaymentMethod = "BANK_TRANSFER_BNI"
	PaymentMethodBankTransferBRI PaymentMethod = "BANK_TRANSFER_BRI"
	PaymentMethodGopay           PaymentMethod = "GOPAY"
	PaymentMethodIndomaret       PaymentMethod = "INDOMARET"
	// Add other payment methods as needed
)

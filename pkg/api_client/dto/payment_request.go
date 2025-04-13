package dto

type PaymentRequestDTO struct {
	CreatedBy       string  `json:"created_by"`
	CustomerID      string  `json:"customer_id"`
	PaymentTypeID   string  `json:"payment_type_id"`
	DstAccountID    string  `json:"dst_account_id"`
	Amount          float64 `json:"amount"`
	ChargeDailyFee  bool    `json:"charge_daily_fee"`
	ClearCredit     bool    `json:"clear_credit"`
	Descr           string  `json:"descr"`
	TransactionID   string  `json:"transaction_id"`
	PaymentSystemID string  `json:"payment_system_id"`
	Source          string  `json:"source"`
}

func NewApiClientPaymentRequestDTO(
	createdBy string,
	customerID string,
	paymentTypeID string,
	dstAccountID string,
	amount float64,
	chargeDailyFee bool,
	clearCredit bool,
	descr string,
	transactionID string,
	paymentSystemID string,
	source string,
) PaymentRequestDTO {
	return PaymentRequestDTO{
		CreatedBy:       createdBy,
		CustomerID:      customerID,
		PaymentTypeID:   paymentTypeID,
		DstAccountID:    dstAccountID,
		Amount:          amount,
		ChargeDailyFee:  chargeDailyFee,
		ClearCredit:     clearCredit,
		Descr:           descr,
		TransactionID:   transactionID,
		PaymentSystemID: paymentSystemID,
		Source:          source,
	}
}

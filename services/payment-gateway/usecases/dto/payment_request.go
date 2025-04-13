package dto

type PaymentRequestDTO struct {
	CustomerID      string
	CreatedBy       string
	PaymentTypeID   string
	DstAccountID    string
	Amount          float64
	Descr           string
	TransactionID   string
	PaymentSystemID string
	Source          string
}

func NewPaymentRequestDTO(
	customerID string,
	createdBy string,
	paymentTypeID string,
	dstAccountID string,
	amount float64,
	descr string,
	transactionID string,
	paymentSystemID string,
	source string,
) PaymentRequestDTO {
	return PaymentRequestDTO{
		CustomerID:      customerID,
		CreatedBy:       createdBy,
		PaymentTypeID:   paymentTypeID,
		DstAccountID:    dstAccountID,
		Amount:          amount,
		Descr:           descr,
		TransactionID:   transactionID,
		PaymentSystemID: paymentSystemID,
		Source:          source,
	}
}

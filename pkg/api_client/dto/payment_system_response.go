package dto

type PaymentSystemResponseDTO struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Token          string `json:"token"`
	Enabled        bool   `json:"enabled"`
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	DstAccountID   string `json:"dst_account_id"`
	PaymentTypeID  string `json:"payment_type_id"`
	TestMode       bool   `json:"test_mode"`
	TestCustomerID string `json:"test_customer_id"`
	PaymentDescr   string `json:"payment_descr"`
}

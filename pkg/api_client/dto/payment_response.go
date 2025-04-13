package dto

type PaymentResponseDTO struct {
	ID        string  `json:"id"`
	CreatedBy string  `json:"created_by"`
	Date      string  `json:"date"`
	Pid       int     `json:"pid"`
	Amount    float64 `json:"amount"`
}

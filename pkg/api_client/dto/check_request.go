package dto

type CheckRequestDTO struct {
	CustomerID string `json:"customer_id"`
}

func NewApiClientCheckRequestDTO(
	customerID string,
) CheckRequestDTO {
	return CheckRequestDTO{
		CustomerID: customerID,
	}
}

package dto

type CheckRequestDTO struct {
	CustomerID string
}

func NewCheckRequestDTO(
	customerID string,
) CheckRequestDTO {
	return CheckRequestDTO{
		CustomerID: customerID,
	}
}

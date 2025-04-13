package interfaces

import "code.evixo.ru/ncc/ncc-backend/pkg/api_client/dto"

type PaymentClient interface {
	Payment(request dto.PaymentRequestDTO) (*dto.PaymentResponseDTO, error)
	Check(request dto.CheckRequestDTO) (*dto.CheckResponseDTO, error)
	GetPaymentSystemByID(id string) (*dto.PaymentSystemResponseDTO, error)
}

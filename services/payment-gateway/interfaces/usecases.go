package interfaces

import (
	dto2 "code.evixo.ru/ncc/ncc-backend/pkg/api_client/dto"
	"code.evixo.ru/ncc/ncc-backend/services/payment-gateway/usecases/dto"
)

type PaymentUsecase interface {
	Execute(request dto.PaymentRequestDTO) (*dto2.PaymentResponseDTO, error)
}

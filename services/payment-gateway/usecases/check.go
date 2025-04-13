package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	dto2 "code.evixo.ru/ncc/ncc-backend/pkg/api_client/dto"
	"code.evixo.ru/ncc/ncc-backend/pkg/api_client/interfaces"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/payment-gateway/domain"
	"code.evixo.ru/ncc/ncc-backend/services/payment-gateway/usecases/dto"
	"fmt"
)

type CheckUsecase struct {
	cfg           *config.Config
	log           logger.Logger
	paymentClient interfaces.PaymentClient
}

func (uc CheckUsecase) Execute(request dto.CheckRequestDTO) error {
	checkRequest := dto2.NewApiClientCheckRequestDTO(request.CustomerID)
	response, err := uc.paymentClient.Check(checkRequest)
	if err != nil {
		return fmt.Errorf("can't check: %w", err)
	}

	if response == nil || !response.PaymentAllowed {
		return domain.ErrPaymentNotAllowed
	}

	return nil
}

func NewCheckUsecase(
	cfg *config.Config,
	log logger.Logger,
	paymentClient interfaces.PaymentClient,
) CheckUsecase {
	return CheckUsecase{
		cfg:           cfg,
		log:           log,
		paymentClient: paymentClient,
	}
}

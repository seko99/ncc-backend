package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	dto2 "code.evixo.ru/ncc/ncc-backend/pkg/api_client/dto"
	"code.evixo.ru/ncc/ncc-backend/pkg/api_client/interfaces"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/payment-gateway/usecases/dto"
	"fmt"
)

type PaymentUsecase struct {
	cfg           *config.Config
	log           logger.Logger
	paymentClient interfaces.PaymentClient
}

func (uc PaymentUsecase) Execute(request dto.PaymentRequestDTO) (*dto2.PaymentResponseDTO, error) {

	paymentRequest := dto2.NewApiClientPaymentRequestDTO(
		request.CreatedBy,
		request.CustomerID,
		request.PaymentTypeID,
		request.DstAccountID,
		request.Amount,
		uc.cfg.PaymentGatewayConfig.ChargeDailyFee,
		uc.cfg.PaymentGatewayConfig.ClearCredit,
		request.Descr,
		request.TransactionID,
		request.PaymentSystemID,
		request.Source,
	)
	payment, err := uc.paymentClient.Payment(paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("can't execute payment: %w", err)
	}

	return payment, err
}

func NewPaymentUsecase(
	cfg *config.Config,
	log logger.Logger,
	paymentClient interfaces.PaymentClient,
) PaymentUsecase {
	return PaymentUsecase{
		cfg:           cfg,
		log:           log,
		paymentClient: paymentClient,
	}
}

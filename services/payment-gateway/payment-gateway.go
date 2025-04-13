package PaymentGateway

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/payment-gateway/handlers"
	"github.com/gin-gonic/gin"
)

type PaymentGateway struct {
	cfg *config.Config
	log logger.Logger
}

type Message struct {
	Informing models2.InformingData
	Customer  models2.CustomerData
	Message   string
	Phone     string
}

func (s *PaymentGateway) Start() error {

	g := gin.New()

	citypayPaymentEndpoint := handlers.NewCitypayPayment(s.cfg, s.log)

	paymentsGroup := g.Group("/payments")
	paymentsGroup.GET("/citypay/:payment_system_id/:token", citypayPaymentEndpoint.Execute())

	err := g.Run(s.cfg.Listen)
	if err != nil {
		s.log.Error("can't listen: %v", err)
	}

	s.log.Info("Running PaymentGateway")

	select {}
}

func NewPaymentGateway(
	cfg *config.Config,
	log logger.Logger,
) *PaymentGateway {
	return &PaymentGateway{
		cfg: cfg,
		log: log,
	}
}

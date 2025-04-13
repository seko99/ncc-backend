package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	PaymentGateway "code.evixo.ru/ncc/ncc-backend/services/payment-gateway"
	"github.com/spf13/cobra"
)

var paymentGatewayCmd = &cobra.Command{
	Use:   "payment-gateway",
	Short: "Payment Gateway",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewConfig()
		if err != nil {
			panic(err)
		}
		log := zero.NewLogger()

		log.Info("Starting Payment Gateway...")
		paymentGateway := PaymentGateway.NewPaymentGateway(cfg, log)

		err = paymentGateway.Start()
		if err != nil {
			panic(err)
		}
	},
}

package utils

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	psql2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"code.evixo.ru/ncc/ncc-backend/services/utils"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"time"
)

var CompensateDuplicateFeesCmd = &cobra.Command{
	Use:   "compensate-duplicate-fees",
	Short: "Compensate duplicate fees",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewConfig()
		if err != nil {
			panic(err)
		}

		debugLevel := zerolog.InfoLevel

		log := zero.NewLogger(debugLevel)

		dryRunFlag, _ := cmd.Flags().GetBool("dry")

		if dryRunFlag {
			log.Info("DRY RUN - NO CHANGES IN DATABASE!")
		}

		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")

		periodStart, err := time.Parse("2006-01-02", start)
		if err != nil {
			panic(fmt.Sprintf("Wrong period start: %v", err))
		}

		periodEnd, err := time.Parse("2006-01-02", end)
		if err != nil {
			panic(fmt.Sprintf("Wrong period end: %v", err))
		}

		log.Info("Compensate duplicate fees from %v to %v", periodStart, periodEnd)

		storage := psqlstorage.NewStorage(cfg, log)
		err = storage.Connect()
		if err != nil {
			panic(fmt.Sprintf("can't connect to storage: %v", err))
		}

		customersRepo := psql2.NewCustomers(storage, nil)
		paymentsRepo := psql2.NewPayments(storage)
		feeRepo := psql2.NewFees(storage)

		feesUtilsService := utils.NewFeesUtils(cfg, log, customersRepo, paymentsRepo, feeRepo)

		err = feesUtilsService.CompensateDuplicateFees(periodStart, periodEnd, dryRunFlag)
		if err != nil {
			panic(fmt.Sprintf("error: %v", err))
		}

		log.Info("Success")
	},
}

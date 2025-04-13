package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/pkg/providers"
	psql2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"code.evixo.ru/ncc/ncc-backend/services/informings"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var informingsCmd = &cobra.Command{
	Use:   "informings",
	Short: "Informings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewConfig()
		if err != nil {
			panic(err)
		}

		debugLevel := zerolog.InfoLevel

		if *debugMode {
			debugLevel = zerolog.DebugLevel
		}
		log := zero.NewLogger(debugLevel)

		log.Info("Starting Informings")

		dryRunFlag, _ := cmd.Flags().GetBool("dry")

		if dryRunFlag {
			log.Info("DRY RUN - NO CHANGES IN DATABASE!")
		}

		storage := psqlstorage.NewStorage(cfg, log)
		err = storage.Connect()
		if err != nil {
			log.Error("can't connect to storage: %v", err)
			return
		}

		informingsRepo := psql2.NewInformings(storage)
		informingsTestCustomersRepo := psql2.NewInformingsTestCustomers(storage)
		informingLogRepo := psql2.NewInformingLog(storage)
		customerRepo := psql2.NewCustomers(storage, nil)

		phoenixSms := providers.NewPhoenixSms(cfg.Informings.SmsProvider, log)

		informings := informings.NewInformings(
			log,
			phoenixSms,
			informingsRepo,
			informingLogRepo,
			informingsTestCustomersRepo,
			customerRepo,
		)

		err = informings.Run(dryRunFlag)
		if err != nil {
			log.Error("can't run informings: %v", err)
			return
		}
	},
}

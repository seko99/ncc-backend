package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	psql2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"code.evixo.ru/ncc/ncc-backend/services/scores"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var scoresCmd = &cobra.Command{
	Use:   "scores",
	Short: "Scores processor",
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

		log.Info("Starting Scores Processor")

		dryRunFlag, _ := cmd.Flags().GetBool("dry")

		if dryRunFlag {
			log.Info("DRY RUN - NO CHANGES IN DATABASE!")
		}

		storage := psqlstorage.NewStorage(cfg, log)
		err = storage.Connect()
		if err != nil {
			panic(fmt.Sprintf("can't connect to storage: %v", err))
		}

		customersRepo := psql2.NewCustomers(storage, nil)
		paymentsRepo := psql2.NewPayments(storage)
		scoreRepo := psql2.NewScores(storage)

		scoresService := scores.NewScores(log, customersRepo, paymentsRepo, scoreRepo)
		err = scoresService.Process(dryRunFlag)
		if err != nil {
			log.Error("can't process scores: %v", err)
			return
		}
	},
}

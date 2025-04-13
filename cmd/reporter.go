package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	psql2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"code.evixo.ru/ncc/ncc-backend/services/reporter"
	"code.evixo.ru/ncc/ncc-backend/services/telegram"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var reporterCmd = &cobra.Command{
	Use:   "reporter",
	Short: "Reporter",
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

		log.Info("Starting Reporter")

		storage := psqlstorage.NewStorage(cfg, log)
		err = storage.Connect()
		if err != nil {
			log.Error("can't connect to storage: %v", err)
			return
		}

		broadcastEvents, err := events.NewEvents(cfg, log, uuid.NewString(), events.BroadcastEvents)
		if err != nil {
			log.Error("Can't init event system: %v", err)
			return
		}

		feeRepo := psql2.NewFees(storage)
		customerRepo := psql2.NewCustomers(storage, broadcastEvents)
		sessionsRepo := psql2.NewSessions(storage, broadcastEvents)
		snapshotRepo := psql2.NewSnapshots(storage)
		paymentsRepo := psql2.NewPayments(storage)
		scoresRepo := psql2.NewScores(storage)

		telegram := telegram.NewTelegram(cfg.Reporter.Telegram.Token)
		err = telegram.Connect()
		if err != nil {
			log.Error("can't connect to bot: %v", err)
			return
		}

		reporter := reporter.NewReporter(
			log,
			cfg.Reporter,
			telegram,
			feeRepo,
			customerRepo,
			sessionsRepo,
			snapshotRepo,
			paymentsRepo,
			scoresRepo,
		)
		err = reporter.Run()
		if err != nil {
			log.Error("can't run reporter: %v", err)
			return
		}
	},
}

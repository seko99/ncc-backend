package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/services/scheduler"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Scheduler",
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

		flagInformings, _ := cmd.Flags().GetBool("informings")
		flagExporter, _ := cmd.Flags().GetBool("exporter")

		log.Info("Starting Scheduler")

		schedulerService := scheduler.NewScheduler(log)

		if flagInformings {
			log.Info("Registering Informings")
			err = scheduler.RegisterInformings(cfg, log, schedulerService)
			if err != nil {
				log.Error("Can't register Informings in scheduler: %v", err)
			}
		}

		if flagExporter {
			log.Info("Registering Exporter")
			err = scheduler.RegisterExporter(cfg, log, schedulerService)
			if err != nil {
				log.Error("Can't register Exporter in scheduler: %v", err)
			}
		}

		err = schedulerService.Run()
		if err != nil {
			log.Error("Can't run scheduler: %v", err)
		}
	},
}

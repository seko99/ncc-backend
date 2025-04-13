package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	psql2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	dhcp2 "code.evixo.ru/ncc/ncc-backend/services/dhcp"
	"code.evixo.ru/ncc/ncc-backend/services/radius"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	dblogger "gorm.io/gorm/logger"
)

var coreCmd = &cobra.Command{
	Use:   "core",
	Short: "Core engine",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewConfig()
		if err != nil {
			panic(err)
		}
		log := zero.NewLogger()

		log.Info("Starting Core engine...")

		storage := psqlstorage.NewStorage(cfg, log,
			psqlstorage.WithLogLevel(dblogger.Error),
			psqlstorage.WithAppName("ncc:CORE"))
		err = storage.Connect()
		if err != nil {
			panic(fmt.Sprintf("can't connect to storage: %v", err))
		}

		broadcastEvents, err := events.NewBroadcastEvents(cfg, log, uuid.NewString(), events.BroadcastEvents)
		if err != nil {
			log.Error("Can't init event system: %v", err)
			return
		}

		customersRepo := psql2.NewCustomers(storage, broadcastEvents)
		nasesRepo := psql2.NewNases(storage, broadcastEvents)
		leasesRepo := psql2.NewDhcpLeases(storage, broadcastEvents)
		sessionsRepo := psql2.NewSessions(storage, broadcastEvents)
		sessionsLogRepo := psql2.NewSessionsLog(storage)
		serviceInternetRepo := psql2.NewServiceInternet(storage, broadcastEvents)
		poolsRepo := psql2.NewDhcpPools(storage, broadcastEvents)
		bindingsRepo := psql2.NewDhcpBindings(storage, broadcastEvents)

		rad, err := radius.NewRadiusEventHandler(
			cfg,
			log,
			customersRepo,
			nasesRepo,
			leasesRepo,
			sessionsRepo,
			sessionsLogRepo,
			serviceInternetRepo,
		)
		if err != nil {
			panic(fmt.Sprintf("can't create radius: %v", err))
		}

		dhcp, err := dhcp2.NewDhcpEventHandler(cfg, log, nasesRepo, leasesRepo, poolsRepo, bindingsRepo)
		if err != nil {
			panic(fmt.Sprintf("can't create dhcp: %v", err))
		}

		log.Info("Starting RadiusEventHandler...")
		go rad.Start()

		log.Info("Starting DhcpEventHandler...")
		go dhcp.Start()

		log.Info("Core engine ready")
		select {}
	},
}

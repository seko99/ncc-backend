package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/memory"
	psql2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"code.evixo.ru/ncc/ncc-backend/services/simulator"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/handlers"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	dblogger "gorm.io/gorm/logger"
)

var simulatorCmd = &cobra.Command{
	Use:   "simulator",
	Short: "Simulator",
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

		log.Info("Starting Simulator")

		daemonFlag, _ := cmd.Flags().GetBool("daemon")

		storage := psqlstorage.NewStorage(cfg, log,
			psqlstorage.WithLogLevel(dblogger.Error),
			psqlstorage.WithAppName("ncc:Simulator"),
		)
		err = storage.Connect()
		if err != nil {
			log.Error("Can't connect to psql storage: %v", err)
			return
		}

		broadcastEvents, err := events.NewBroadcastEvents(cfg, log, uuid.NewString(), events.BroadcastEvents)
		if err != nil {
			log.Error("Can't init event system: %v", err)
			return
		}

		customersRepo := psql2.NewCustomers(storage, broadcastEvents)
		customerGroupsRepo := psql2.NewCustomerGroups(storage)
		citiesRepo := psql2.NewCities(storage)
		streetsRepo := psql2.NewStreets(storage)
		serviceInternetRepo := psql2.NewServiceInternet(storage, broadcastEvents)
		vendorsRepo := psql2.NewVendors(storage)
		hardwareModelsRepo := psql2.NewHardwareModels(storage)
		paymentTypesRepo := psql2.NewPaymentTypes(storage)
		mapNodesRepo := psql2.NewMapNodes(storage)
		devicesRepo := psql2.NewDevices(storage)
		deviceStatesRepo := psql2.NewDeviceStates(storage)
		ifacesRepo := psql2.NewDeviceInterfaces(storage)
		ifaceStatesRepo := psql2.NewDeviceInterfaceStates(storage)
		dhcpBindingsRepo := psql2.NewDhcpBindings(storage, broadcastEvents)
		dhcpLeasesRepo := psql2.NewDhcpLeases(storage, broadcastEvents)
		sessionsRepo := psql2.NewSessions(storage, broadcastEvents)
		nasesRepo := psql2.NewNases(storage, broadcastEvents)
		nasTypesRepo := psql2.NewNasTypes(storage)
		paymentsRepo := psql2.NewPayments(storage)
		feesRepo := psql2.NewFees(storage)
		contractsRepo := psql2.NewContracts(storage)
		paymentSystemsRepo := psql2.NewPaymentSystems(storage)
		usersRepo := psql2.NewUsers(storage)
		radiusVendorsRepo := psql2.NewRadiusVendors(storage)
		radiusAttributesRepo := psql2.NewRadiusAttributes(storage)
		issueTypesRepo := psql2.NewIssueTypes(storage)
		issueUrgenciesRepo := psql2.NewIssueUrgencies(storage)
		issuesRepo := psql2.NewIssues(storage)
		issueActionsRepo := psql2.NewIssueActions(storage)
		dhcpPoolsRepo := psql2.NewDhcpPools(storage, broadcastEvents)
		ipPoolRepo := psql2.NewIpPools(storage, broadcastEvents)

		nasCache, err := memory.NewNases(log, nasesRepo, broadcastEvents)
		if err != nil {
			log.Error("Can't init NAS cache: %v", err)
			return
		}

		sessionCache, err := memory.NewSessions(log, sessionsRepo, nil)
		if err != nil {
			log.Error("Can't init session cache: %v", err)
			return
		}

		leasesCache, err := memory.NewDhcpLeases(log, dhcpLeasesRepo, broadcastEvents)
		if err != nil {
			log.Error("Can't init lease cache: %v", err)
			return
		}

		customerCache, err := memory.NewCustomers(log, customersRepo, broadcastEvents)
		if err != nil {
			log.Error("Can't init customer cache: %v", err)
			return
		}

		events, err := events.NewEvents(cfg, log, uuid.NewString(), "simulator.events")
		if err != nil {
			panic(fmt.Sprintf("Can't init event system: %v", err))
		}

		simulatorService := simulator.NewSimulator(cfg, log, storage, events,
			customerGroupsRepo,
			citiesRepo,
			streetsRepo,
			customersRepo,
			serviceInternetRepo,
			vendorsRepo,
			hardwareModelsRepo,
			paymentTypesRepo,
			mapNodesRepo,
			devicesRepo,
			deviceStatesRepo,
			ifacesRepo,
			ifaceStatesRepo,
			dhcpBindingsRepo,
			dhcpLeasesRepo,
			sessionsRepo,
			nasesRepo,
			nasTypesRepo,
			paymentsRepo,
			feesRepo,
			contractsRepo,
			paymentSystemsRepo,
			usersRepo,
			radiusVendorsRepo,
			radiusAttributesRepo,
			dhcpPoolsRepo,
			ipPoolRepo,
			issueTypesRepo,
			issueUrgenciesRepo,
			issuesRepo,
			issueActionsRepo,

			sessionCache,
			leasesCache,
			customerCache,
			nasCache,
		)

		if daemonFlag {
			g := gin.New()

			sessionsDropUsecase := usecases.NewSessionsDropUsecase(log, simulatorService)
			sessionsDropEndpoint := handlers.NewSessionsDrop(log, &sessionsDropUsecase)

			sessionsUpdateUsecase := usecases.NewSessionsUpdateUsecase(log, simulatorService)
			sessionsUpdateEndpoint := handlers.NewSessionsUpdateEndpoint(log, &sessionsUpdateUsecase)

			leasesUpdateUsecase := usecases.NewLeasesUpdateUsecase(log, simulatorService)
			leasesUpdateEndpoint := handlers.NewLeasesUpdateEndpoint(log, &leasesUpdateUsecase)

			issuesCreateUsecase := usecases.NewIssuesCreateUsecase(log, simulatorService)
			issuesCreateEndpoint := handlers.NewIssuesCreateEndpoint(log, &issuesCreateUsecase)

			issuesDeleteAllUsecase := usecases.NewIssuesDeleteAllUsecase(log, simulatorService)
			issuesDeleteAllEndpoint := handlers.NewIssuesDeleteAllEndpoint(log, &issuesDeleteAllUsecase)

			initDictionariesUsecase := usecases.NewInitDictionariesUsecase(log, simulatorService)
			initDictionariesEndpoint := handlers.NewInitDictionariesEndpoint(log, &initDictionariesUsecase)
			fakeDataCreateUsecase := usecases.NewFakeDataCreateUsecase(log, simulatorService)
			fakeDataCreateEndpoint := handlers.NewFakeDataCreateEndpoint(log, &fakeDataCreateUsecase)
			fakeDataClearUsecase := usecases.NewFakeDataClearUsecase(log, simulatorService)
			fakeDataClearEndpoint := handlers.NewFakeDataClearEndpoint(log, &fakeDataClearUsecase)

			radiusStartAllUsecase := usecases.NewRadiusStartAllUsecase(log, simulatorService)
			radiusStartAllEndpoint := handlers.NewRadiusStartAllEndpoint(cfg, log, &radiusStartAllUsecase)
			radiusStopAllUsecase := usecases.NewRadiusStopAllUsecase(log, simulatorService)
			radiusStopAllEndpoint := handlers.NewRadiusStopAllEndpoint(cfg, log, &radiusStopAllUsecase)
			radiusUpdateAllUsecase := usecases.NewRadiusUpdateAllUsecase(log, simulatorService)
			radiusUpdateAllEndpoint := handlers.NewRadiusUpdateAllEndpoint(cfg, log, &radiusUpdateAllUsecase)
			radiusKillSessionsUsecase := usecases.NewRadiusKillSessionsUsecase(log, simulatorService)
			radiusKillSessionsEndpoint := handlers.NewRadiusKillSessionsEndpoint(cfg, log, &radiusKillSessionsUsecase)

			brasGetSessionsUsecase := usecases.NewBrasGetSessionsUsecase(log, simulatorService)
			brasGetSessionsEndpoint := handlers.NewBrasGetSessionsEndpoint(cfg, log, &brasGetSessionsUsecase)

			brasGetStatUsecase := usecases.NewBrasGetStatUsecase(log, simulatorService)
			brasGetStatEndpoint := handlers.NewBrasGetStatEndpoint(cfg, log, &brasGetStatUsecase)

			brasSetParamsUsecase := usecases.NewBrasSetParamsUsecase(log, simulatorService)
			brasSetParamsEndpoint := handlers.NewBrasSetParamsEndpoint(cfg, log, &brasSetParamsUsecase)

			apiGroup := g.Group("/api")

			v1Group := apiGroup.Group("/v1")

			v1Group.POST("/dictionary/init", initDictionariesEndpoint.Execute())

			fakeGroup := v1Group.Group("/fake")
			{
				fakeGroup.POST("/create", fakeDataCreateEndpoint.Execute())
				fakeGroup.POST("/clear", fakeDataClearEndpoint.Execute())
			}

			brasGroup := v1Group.Group("/bras")
			{
				brasGroup.GET("/sessions", brasGetSessionsEndpoint.Execute())
				brasGroup.GET("/stat", brasGetStatEndpoint.Execute())
				brasGroup.PATCH("/params", brasSetParamsEndpoint.Execute())
			}

			radiusGroup := v1Group.Group("/radius")
			{
				//todo: rename and move to BRAS section
				radiusGroup.POST("/start", radiusStartAllEndpoint.Execute())
				radiusGroup.POST("/stop", radiusStopAllEndpoint.Execute())
				radiusGroup.POST("/update", radiusUpdateAllEndpoint.Execute())
				radiusGroup.POST("/kill", radiusKillSessionsEndpoint.Execute())
			}

			issuesGroup := v1Group.Group("/issues")
			{
				issuesGroup.POST("/create", issuesCreateEndpoint.Execute())
				issuesGroup.POST("/delete/all", issuesDeleteAllEndpoint.Execute())
			}

			sessionsGroup := v1Group.Group("/sessions")
			{
				sessionsGroup.POST("/drop", sessionsDropEndpoint.Execute())
				sessionsGroup.POST("/update", sessionsUpdateEndpoint.Execute())
				//sessionsGroup.POST("/disconnect", initDictionariesEndpoint.Execute())
			}

			leasesGroup := v1Group.Group("/leases")
			{
				leasesGroup.POST("/update", leasesUpdateEndpoint.Execute())
			}

			//mapGroup := v1Group.Group("/map")
			{
				//mapGroup.POST("/update", sessionsDropEndpoint.Execute())
			}

			err = g.Run(cfg.Listen)
			if err != nil {
				log.Error("can't listen: %v", err)
			}
		}
	},
}

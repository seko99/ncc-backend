package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	psql2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"code.evixo.ru/ncc/ncc-backend/services/fees"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var feeSchedulerCmd = &cobra.Command{
	Use:   "fee",
	Short: "Fee processor",
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

		log.Info("Starting Fee Processor")

		timeFlag, _ := cmd.Flags().GetString("time")
		dryRunFlag, _ := cmd.Flags().GetBool("dry")
		yesFlag, _ := cmd.Flags().GetBool("yes")
		ignoreCreditExpireFlag, _ := cmd.Flags().GetBool("ignore-credit-expire")
		ignoreServiceStateFlag, _ := cmd.Flags().GetBool("ignore-service-state")
		loginFlag, _ := cmd.Flags().GetString("login")
		maxBlocksFlag, _ := cmd.Flags().GetInt("max-blocks")
		uidFlag, _ := cmd.Flags().GetString("uid")
		zeroFeeFlag, _ := cmd.Flags().GetBool("zero-fee")
		sessionStopFlag, _ := cmd.Flags().GetString("session-stop")

		uids := strings.Split(uidFlag, ",")

		if dryRunFlag {
			log.Info("DRY RUN - NO CHANGES IN DATABASE!")
		}

		forTime, err := time.Parse("2006-01-02 15:04:05", timeFlag)
		if err != nil {
			log.Error("Invalid time: %v", err)
			return
		}

		log.Info("Processing fee for %s", timeFlag)

		storage := psqlstorage.NewStorage(cfg, log)
		err = storage.Connect()
		if err != nil {
			panic(fmt.Sprintf("can't connect to storage: %v", err))
		}

		feeRepo := psql2.NewFees(storage)
		customerRepo := psql2.NewCustomers(storage, nil)
		serviceInternetRepo := psql2.NewServiceInternet(storage, nil)
		ipPoolRepo := psql2.NewIpPools(storage, nil)
		preferencesRepo := psql2.NewPreferences(storage)

		feesService := fees.NewFees(log, feeRepo, customerRepo, serviceInternetRepo, ipPoolRepo)

		internets, err := serviceInternetRepo.Get()
		if err != nil {
			log.Error("Can't get serviceInternet: %v", err)
			return
		}
		internetMap := map[string]models2.ServiceInternetData{}
		for _, i := range internets {
			internetMap[i.Id] = i
		}

		var customers []models2.CustomerData

		if loginFlag != "" {
			customer, err := customerRepo.GetByLogin(loginFlag)
			if err != nil {
				log.Error("Can't get customer: %v", err)
				return
			}
			customers = append(customers, *customer)
		} else if uidFlag != "" {
			if len(uids) > 0 {
				for _, uuid := range uids {
					customer, err := customerRepo.GetByUid(uuid)
					if err != nil {
						log.Error("Can't get customer: %v", err)
						return
					}
					customers = append(customers, *customer)
				}
			} else {
				customer, err := customerRepo.GetByUid(uidFlag)
				if err != nil {
					log.Error("Can't get customer: %v", err)
					return
				}
				customers = append(customers, *customer)
			}
		} else if zeroFeeFlag {
			customers, err = customerRepo.GetByFeeAmountAndSessions(0.0, sessionStopFlag)
			if err != nil {
				log.Error("Can't get customers: %v", err)
				return
			}
		} else {
			customers, err = customerRepo.Get()
			if err != nil {
				log.Error("Can't get customers: %v", err)
				return
			}
		}

		todayFees, err := feeRepo.GetProcessedMap(forTime)
		if err != nil {
			log.Error("Can't get processed map: %v", err)
			return
		}

		customDataMap, err := serviceInternetRepo.GetCustomDataMap()
		if err != nil {
			log.Error("Can't get custom data map: %v", err)
			return
		}

		log.Info("Processing %d customers with already processed %d", len(customers), len(todayFees))

		if len(customers) > 1 {
			if !yesFlag && !askForConfirmation(fmt.Sprintf("Process %d customers?", len(customers))) {
				return
			}
		}

		err = preferencesRepo.SetFeeProcessingInProgress(true)
		if err != nil {
			log.Error("Can't set fee processing flag: %v", err)
			return
		}

		defer func() {
			err := preferencesRepo.SetFeeProcessingInProgress(false)
			if err != nil {
				log.Error("Can't reset fee processing flag: %v", err)
			}
		}()

		ts := time.Now()
		days := feesService.DaysIn(time.Now().Month(), time.Now().Year())

		feeDatas, err := feesService.Process(
			internetMap,
			customers,
			customDataMap,
			todayFees,
			days,
			forTime,
			dryRunFlag,
			maxBlocksFlag,
			ignoreCreditExpireFlag,
			ignoreServiceStateFlag,
		)
		if err != nil {
			log.Error("Can't process fees: %v", err)
		}

		log.Info("Processed: %d in %v", len(feeDatas), time.Since(ts))
	},
}

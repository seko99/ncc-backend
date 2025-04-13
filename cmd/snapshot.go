package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	psql2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"time"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Snapshot processor",
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

		loginFlag, _ := cmd.Flags().GetString("login")

		log.Info("Starting Snapshot Processor")

		storage := psqlstorage.NewStorage(cfg, log)
		err = storage.Connect()
		if err != nil {
			panic(fmt.Sprintf("can't connect to storage: %v", err))
		}

		broadcastEvents, err := events.NewEvents(cfg, log, uuid.NewString(), events.BroadcastEvents)
		if err != nil {
			log.Error("Can't init event system: %v", err)
			return
		}

		customersRepo := psql2.NewCustomers(storage, broadcastEvents)
		snapshotRepo := psql2.NewSnapshots(storage)
		serviceInternetRepo := psql2.NewServiceInternet(storage, broadcastEvents)
		sessionLogRepo := psql2.NewSessionsLog(storage)
		sessionsRepo := psql2.NewSessions(storage, broadcastEvents)

		var customers []models2.CustomerData

		if loginFlag != "" {
			customer, err := customersRepo.GetByLogin(loginFlag)
			if err != nil {
				log.Error("Can't get customer: %v", err)
				return
			}
			customers = append(customers, *customer)
		} else {
			customers, err = customersRepo.Get()
			if err != nil {
				log.Error("Can't get customers: %v", err)
				return
			}
		}

		customDataMap, err := serviceInternetRepo.GetCustomDataMap()
		if err != nil {
			log.Error("Can't get custom data map: %v", err)
			return
		}

		processed := 0
		for _, c := range customers {
			log.Info("[%s] Taking snapshot", c.Login)

			hasSessions := false

			sessions, err := sessionLogRepo.GetByCustomer(c.Id, repository.TimePeriod{
				In: time.Now(),
			})
			if err != nil {
				log.Error("[%s] Can't get session log: %v", c.Login, err)
			}
			if sessions != nil && len(sessions) > 0 {
				hasSessions = true
			}

			online, err := sessionsRepo.GetByCustomer(c.Id, repository.TimePeriod{
				End: time.Now(),
			})
			if err != nil {
				log.Error("[%s] Can't get sessions: %v", c.Login, err)
			}
			if online != nil && len(online) > 0 {
				hasSessions = true
			}

			snapshot := &models2.SnapshotData{
				Uid:             c.Uid,
				Login:           c.Login,
				Deposit:         c.Deposit,
				Credit:          c.Credit,
				Scores:          c.Scores,
				BlockingDate:    c.BlockingDate,
				BlockingTill:    c.BlockingTill,
				BlockingLastSet: c.BlockingLastSet,
				BlockingState:   c.BlockingState,
				CreditExpire:    c.CreditExpire.Time,
				CreditDaysLeft:  c.CreditDaysLeft,

				ServiceInternetId:    c.ServiceInternetId,
				ServiceInternetState: c.ServiceInternetState,
				ServiceInternetFee:   c.ServiceInternet.Fee,
				ServiceInternetIPFee: c.ServiceInternet.IpFee,

				ServiceIptvId:    c.ServiceIptvId,
				ServiceIptvState: c.ServiceIptvState,
				ServiceIptvFee:   c.ServiceIptv.Fee,

				ServiceCatvId:    c.ServiceCatvId,
				ServiceCatvState: c.ServiceCatvState,
				ServiceCatvFee:   c.ServiceCatv.Fee,

				HasSessions: hasSessions,
			}

			customData, ok := customDataMap[c.Id]
			if ok {
				snapshot.ServiceInternetIP = customData.Ip
				snapshot.ServiceInternetSpeedIn = customData.SpeedIn
				snapshot.ServiceInternetSpeedOut = customData.SpeedOut
			}

			err = snapshotRepo.Create(snapshot)
			if err != nil {
				log.Error("Can't create snapshot: %v", err)
				continue
			}

			processed++
		}

		log.Info("Processed %d snapshots of %d customers", processed, len(customers))
	},
}

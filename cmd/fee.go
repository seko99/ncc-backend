package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
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
		fixFlag, _ := cmd.Flags().GetBool("fix")
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
		sessionsRepo := psql2.NewSessionsLog(storage)
		ipPoolRepo := psql2.NewIpPools(storage, nil)

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

		pools, err := ipPoolRepo.Get()
		if err != nil {
			log.Error("Can't get IP pools: %v", err)
			return
		}

		if fixFlag {
			paymentsRepo := psql2.NewPayments(storage)

			start := time.Date(forTime.Year(), forTime.Month(), forTime.Day(), 0, 0, 0, 0, forTime.Location())
			log.Info("Fixing since %v", start)
			startTime := time.Now()

			for _, c := range customers {
				prevFee, err := feeRepo.GetByCustomer(c.Id, repository.TimePeriod{
					End: start,
				}, 1)
				if err != nil {
					log.Error("[%s] Can't get prev fee by customer: %v", c.Login, err)
					continue
				}

				if prevFee == nil || len(prevFee) == 0 {
					log.Error("[%s] Can't get prev fee", c.Login)
					continue
				}

				log.Info("[%s] Prev fee at %v deposit=%0.2f", c.Login, prevFee[0].FeeTimestamp, prevFee[0].NewDeposit)

				ts := forTime

				deposit := prevFee[0].NewDeposit
				credit := prevFee[0].Credit
				sumPayments := 0.0
				startPayment := ts.Add(-24 * time.Hour)

				timeShift := time.Time{}
				if prevFee[0].FeeTimestamp.Before(start.Add(-24*time.Hour)) || (deposit <= 0.0 && credit <= 0.0) {
					log.Info("[%s] No fee at %v or negative deposit", c.Login, forTime)

					sessions, err := sessionsRepo.GetByCustomer(c.Id, repository.TimePeriod{
						Start:       forTime,
						ExtraFields: []string{"stop_time"},
					}, 1)
					if err != nil {
						log.Error("[%s] Can't get sessions: %v", c.Login, err)
						continue
					}

					var firstSession time.Time
					if len(sessions) > 0 {
						if !firstSession.Before(forTime) {
							st := sessions[0].StartTime
							firstSession = st
							timeShift = firstSession
							ts = time.Date(st.Year(), st.Month(), st.Day(), forTime.Hour(), forTime.Minute(), forTime.Second(), 0, st.Location())
						} else {
							st := sessions[0].StartTime
							firstSession = st
						}
						log.Info("[%s] First session at %v", c.Login, firstSession)
					} else {
						log.Info("[%s] No sessions since %v", c.Login, forTime)
					}

					pStart := prevFee[0].FeeTimestamp
					pEnd := firstSession

					payments, err := paymentsRepo.GetPaymentsByCustomer(c.Id, repository.TimePeriod{
						Start: pStart,
						End:   pEnd,
					})
					if err != nil {
						log.Error("[%s] Can't get payments: %v", c.Login, err)
						continue
					}

					pSum := 0.0
					firstPayment := false
					for _, p := range payments {
						if p.PaymentTypeId.UUID.String() != "12349064-a1d9-156a-2550-316948ab13fe" &&
							p.PaymentTypeId.UUID.String() != "57bed6af-7bbc-888f-e7d5-e86c9410ee87" &&
							p.PaymentTypeId.UUID.String() != "612899be-22cd-fbdc-aa81-594d13f4d55d" &&
							p.PaymentTypeId.UUID.String() != "30ba6fe9-7205-015c-1080-22ee96b65556" &&
							p.PaymentTypeId.UUID.String() != "b9d3ae5b-0eb0-8934-3cc1-746704875b16" {
							continue
						}
						pSum += p.Amount
						deposit += p.Amount
						startPayment = p.Date
						if deposit > 0.0 && !firstPayment {
							ts = time.Date(p.Date.Year(), p.Date.Month(), p.Date.Day(), ts.Hour(), ts.Minute(), ts.Second(), 0, ts.Location()).Add(24 * time.Hour)
							timeShift = ts
							firstPayment = true
							log.Info("[%s] Setting start at %v", c.Login, ts)
						}
					}

					sumPayments += pSum

					log.Info("[%s] Adding beginning payments: pSum=%0.2f deposit=%0.2f start=%v", c.Login, pSum, deposit, startPayment)
				}

				startDeposit := deposit
				sumFee := 0.0

				fixedLog := []models2.FeeLogData{}

				singleStartTime := time.Now()
				for {
					log.Info("[%s] Processing at %v", c.Login, ts)
					if ts.After(time.Now()) {
						log.Info("[%s] Last day reached ts=%v", c.Login, ts)
						break
					}

					existingFees, err := feeRepo.GetByCustomer(c.Id, repository.TimePeriod{
						Start: time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location()),
						End:   time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location()).Add(23 * time.Hour),
					})

					var fee models2.FeeLogData

					if len(existingFees) > 0 {
						fee = existingFees[0]
					} else {
						fee = models2.FeeLogData{
							FeeTimestamp: ts,
						}
					}

					pStart := startPayment.Add(time.Minute)
					pEnd := fee.FeeTimestamp

					payments, err := paymentsRepo.GetPaymentsByCustomer(c.Id, repository.TimePeriod{
						Start: pStart,
						End:   pEnd,
					})
					if err != nil {
						log.Error("[%s] Can't get payments: %v", c.Login, err)
						ts = ts.Add(24 * time.Hour)
						continue
					}

					pSum := 0.0
					for _, p := range payments {
						if p.PaymentTypeId.UUID.String() != "12349064-a1d9-156a-2550-316948ab13fe" &&
							p.PaymentTypeId.UUID.String() != "57bed6af-7bbc-888f-e7d5-e86c9410ee87" &&
							p.PaymentTypeId.UUID.String() != "612899be-22cd-fbdc-aa81-594d13f4d55d" &&
							p.PaymentTypeId.UUID.String() != "30ba6fe9-7205-015c-1080-22ee96b65556" &&
							p.PaymentTypeId.UUID.String() != "b9d3ae5b-0eb0-8934-3cc1-746704875b16" {
							continue
						}
						pSum += p.Amount
						log.Info("[%s] Adding payment amount=%0.2f at %v", c.Login, p.Amount, p.Date)
					}

					startPayment = pEnd
					sumPayments += pSum

					deposit += pSum

					if deposit <= 0.0 && credit <= 0.0 {
						sessions, err := sessionsRepo.GetByCustomer(c.Id, repository.TimePeriod{
							In: ts,
						}, 1)
						if err != nil {
							log.Error("[%s] Can't get sessions: %v", c.Login, err)
							ts = ts.Add(24 * time.Hour)
							continue
						}
						if len(sessions) == 0 {
							log.Info("[%s] No sessions while negative deposit", c.Login)
							start := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location())
							end := time.Date(ts.Year(), ts.Month(), ts.Day(), 23, 59, 59, 0, ts.Location())
							log.Info("[%s] Deleting fees from %v to %v", c.Login, start, end)
							if !dryRunFlag {
								err := feeRepo.DeleteByCustomer(c.Id, repository.TimePeriod{
									Start: start,
									End:   end,
								})
								if err != nil {
									log.Error("Can't delete fee for %s: %v", c.Login, err)
								}
							}

							ts = ts.Add(24 * time.Hour)
							continue
						}
					}

					days := feesService.DaysIn(ts.Month(), ts.Year())

					services, err := feesService.GetFeeServices(internetMap, c, customDataMap, map[string]models2.CustomerData{}, days, pools, ignoreServiceStateFlag)
					if err != nil {
						log.Error("[%s] Can't get fee services: %v", c.Login, err)
						ts = ts.Add(24 * time.Hour)
						continue
					}

					c.Deposit = deposit
					c.Credit = credit
					if c.Deposit <= 0.0 {
						c.Credit = 0.0
						c.CreditExpire.Time = time.Now().Add(24 * time.Hour)
					}
					feeData, err := feesService.CreateFee(c, services, "", fee.FeeTimestamp, ignoreCreditExpireFlag)
					if err != nil {
						log.Error("[%s] Can't create fee: %v", c.Login, err)
						ts = ts.Add(24 * time.Hour)
						continue
					}
					feeData.FeeLog.Id = fee.Id

					log.Info("[%s] Created fee: %v %0.2f => %0.2f", c.Login, feeData.FeeLog.FeeTimestamp, feeData.FeeLog.FeeAmount, feeData.NewDeposit)

					fixedLog = append(fixedLog, feeData.FeeLog)

					sumFee += feeData.FeeLog.FeeAmount
					deposit = feeData.NewDeposit

					if !dryRunFlag {
						if len(existingFees) > 0 {
							err := feeRepo.Update(feeData.FeeLog)
							if err != nil {
								log.Error("[%s] Can't update fee: %v", c.Login, err)
								continue
							}
						} else {
							err := feeRepo.Create(feeData.FeeLog)
							if err != nil {
								log.Error("[%s] Can't create fee: %v", c.Login, err)
								continue
							}
						}

						err = customerRepo.SetDeposit(c.Id, float64(int(feeData.FeeLog.NewDeposit*100))/100)
						if err != nil {
							log.Error("[%s] Can't set deposit: %v", c.Login, err)
							continue
						}
					}

					ts = ts.Add(24 * time.Hour)
				}

				pStart := ts.Add(-24 * time.Hour)
				pEnd := ts

				payments, err := paymentsRepo.GetPaymentsByCustomer(c.Id, repository.TimePeriod{
					Start: pStart,
					End:   pEnd,
				})
				if err != nil {
					log.Error("[%s] Can't get last day payments: %v", c.Login, err)
				} else {
					if len(payments) > 0 {
						pSum := 0.0
						for _, p := range payments {
							if p.PaymentTypeId.UUID.String() != "12349064-a1d9-156a-2550-316948ab13fe" &&
								p.PaymentTypeId.UUID.String() != "57bed6af-7bbc-888f-e7d5-e86c9410ee87" &&
								p.PaymentTypeId.UUID.String() != "612899be-22cd-fbdc-aa81-594d13f4d55d" &&
								p.PaymentTypeId.UUID.String() != "30ba6fe9-7205-015c-1080-22ee96b65556" &&
								p.PaymentTypeId.UUID.String() != "b9d3ae5b-0eb0-8934-3cc1-746704875b16" {
								continue
							}
							pSum += p.Amount
						}

						if pSum > 0.0 {
							deposit = deposit + pSum
							sumPayments += pSum
							log.Info("[%s] Adding last day payments: pSum=%0.2f newDeposit=%0.2f", c.Login, pSum, deposit)
							if !dryRunFlag {
								err = customerRepo.SetDeposit(c.Id, float64(int(deposit*100))/100)
								if err != nil {
									log.Error("[%s] Can't set final deposit: %v", c.Login, err)
								}
							}
						}
					}
				}

				if !timeShift.IsZero() {
					start := prevFee[0].FeeTimestamp.Add(time.Hour)
					log.Info("[%s] Deleting incorrect fees from %v to %v", c.Login, start, timeShift)
					if !dryRunFlag {
						err := feeRepo.DeleteByCustomer(c.Id, repository.TimePeriod{
							Start: start,
							End:   timeShift.Add(-time.Hour),
						})
						if err != nil {
							log.Error("Can't delete fee for %s: %v", c.Login, err)
						}
					}
				}

				log.Info("[%s] days=%d startDeposit=%0.2f payments=%0.2f fee=%0.2f deposit=%0.2f in %v", c.Login, len(fixedLog), startDeposit, sumPayments, sumFee, deposit, time.Since(singleStartTime))
			}

			log.Info("Processed %d customers in %v", len(customers), time.Since(startTime))
			return
		}

		log.Info("Processing %d customers with already processed %d", len(customers), len(todayFees))

		if len(customers) > 1 {
			if !yesFlag && !askForConfirmation(fmt.Sprintf("Process %d customers?", len(customers))) {
				return
			}
		}

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

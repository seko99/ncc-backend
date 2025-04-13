package fees

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/domain"
	"code.evixo.ru/ncc/ncc-backend/pkg/helpers"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	repository2 "code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"fmt"
	"strconv"
	"time"
)

type Fees struct {
	log                 *zero.Logger
	feeRepo             repository2.Fees
	customerRepo        repository2.Customers
	serviceInternetRepo repository2.ServiceInternet
	ipPoolRepo          repository2.IpPools
}

func (s *Fees) GetFeeServices(
	internet map[string]models2.ServiceInternetData,
	c models2.CustomerData,
	customDataMap map[string]models2.ServiceInternetCustomData,
	todayFees map[string]models2.CustomerData,
	days int,
	pools []models2.IpPoolData,
	ignoreServiceState bool,
) ([]domain.FeeService, error) {
	services := []domain.FeeService{}
	totalFee := 0.0

	if (c.BlockingState > 0) && (c.BlockingState != models2.CustomerStateActive) {
		s.log.Debug("[%s] Customer is not in active state: %d", c.Login, c.BlockingState)
		return nil, nil
	}

	f, processed := todayFees[c.Login]
	if processed {
		s.log.Debug("[%s] Customer already processed at %v", c.Login, f.CreateTs)
		return nil, nil
	}

	if c.ServiceInternetId.Valid && c.ServiceInternetId.UUID.String() != "" {
		inet, ok := internet[c.ServiceInternetId.UUID.String()]
		if !ok {
			return nil, fmt.Errorf("[%s] Can't find serviceInternet", c.Login)
		}
		if inet.FeeType != 0 {
			if c.ServiceInternetState == models2.ServiceStateEnabled || ignoreServiceState {

				switch inet.FeeType {
				case models2.FeeTypeDaily:
					fee := float64(int(inet.Fee/float64(days)*100)) / 100

					customData, ok := customDataMap[c.Id]

					if ok {
						customFee, err := s.getCustomFee(customData, inet, days, pools)
						if err != nil {
							s.log.Error("[%s] CustomFee error: %v", c.Login, err)
						}
						if customFee != nil {
							fee += customFee.Fee
						}
					}

					totalFee += fee
					services = append(services, domain.FeeService{
						Service: inet,
						Fee:     fee,
					})
					s.log.Debug("[%s] Added service %s with fee %0.2f", c.Login, c.ServiceInternet.Name, fee)
				case models2.FeeTypeMonthly:
					if inet.FeeDate == time.Now().Day() {
						fee := inet.Fee

						customData, ok := customDataMap[c.Id]

						if ok {
							customFee, err := s.getCustomFee(customData, inet, days, pools)
							if err != nil {
								s.log.Error("[%s] CustomFee error: %v", c.Login, err)
							}
							if customFee != nil {
								fee += customFee.Fee
							}
						}

						totalFee += fee
						services = append(services, domain.FeeService{
							Service: inet,
							Fee:     fee,
						})
						s.log.Debug("[%s] Added service %s with fee %0.2f", c.Login, inet.Name, fee)
					}
				default:
					s.log.Warn("[%s] Unknown fee type: %v", c.Login, inet.FeeType)
				}
			} else {
				s.log.Debug("[%s] ServiceInternet disabled", c.Login)
			}
		} else {
			s.log.Debug("[%s] Zero fee type", c.Login)
		}
	} else {
		s.log.Debug("[%s] No ServiceInternet", c.Login)
	}

	return services, nil
}

func (s *Fees) getCustomFee(
	customData models2.ServiceInternetCustomData,
	service models2.ServiceInternetData,
	days int,
	pools []models2.IpPoolData,
) (*domain.FeeService, error) {
	if customData.Ip != "" {
		pool, err := helpers.GetPoolByIP(pools, customData.Ip)
		if err != nil {
			return nil, fmt.Errorf("can't get pool by IP: %w", err)
		}

		if !pool.IsPaid {
			s.log.Debug("[%s] Free IP pool for %s", customData.Customer.Login, customData.Ip)
			return nil, nil
		}

		ipFee := service.IpFee / float64(days)
		if customData.IpFee > 0 {
			ipFee = customData.IpFee / float64(days)
		}
		s.log.Debug("[%s] Added IP fee: %s %0.2f", customData.Customer.Login, customData.Ip, ipFee)
		return &domain.FeeService{
			Service: service,
			Fee:     ipFee,
		}, nil
	}

	return nil, nil
}

func (s *Fees) Process(
	internet map[string]models2.ServiceInternetData,
	customers []models2.CustomerData,
	customDataMap map[string]models2.ServiceInternetCustomData,
	todayFees map[string]models2.CustomerData,
	days int,
	forTime time.Time,
	dryRun bool,
	maxBlocks int,
	ignoreCreditExpire bool,
	ignoreServiceState bool,
) (map[string]domain.Fee, error) {

	days = s.DaysIn(forTime.Month(), forTime.Year())

	feeDatas := map[string]domain.Fee{}

	customersProcessed := 0

	pools, err := s.ipPoolRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get IP pools: %w", err)
	}

	blocks := 0
	for _, c := range customers {
		services, err := s.GetFeeServices(internet, c, customDataMap, todayFees, days, pools, ignoreServiceState)
		if err != nil {
			s.log.Error("Can't get fee services: %v", err)
			continue
		}

		if len(services) > 0 {
			feeData, err := s.CreateFee(c, services, "", forTime, ignoreCreditExpire)
			if err != nil {
				s.log.Error("Can't create fee: %v", err)
			}
			feeData.Login = c.Login
			feeData.ServiceInternetState = c.ServiceInternetState

			if feeData != nil {
				feeDatas[c.Login] = *feeData
			}

			s.log.Debug("[%s] Created fee: %v %0.2f => %0.2f", c.Login, feeData.FeeLog.FeeTimestamp, feeData.FeeLog.FeeAmount, feeData.NewDeposit)
		} else {
			s.log.Debug("[%s] No services found", c.Login)
		}

		customersProcessed++
	}

	fees := []models2.FeeLogData{}

	for _, feeData := range feeDatas {
		fees = append(fees, feeData.FeeLog)
		if !dryRun {
			err := s.feeRepo.Create(feeData.FeeLog)
			if err != nil {
				s.log.Error("[%s] Can't create fee log: %v", feeData.Login, err)
			}
		}

		s.log.Debug("[%s] Setting deposit: %0.2f", feeData.Login, feeData.NewDeposit)
		if !dryRun {
			err := s.customerRepo.SetDeposit(feeData.FeeLog.CustomerId.UUID.String(), float64(int(feeData.NewDeposit*100))/100)
			if err != nil {
				s.log.Error("[%s] Can't set deposit: %v", feeData.Login, err)
				continue
			}
		}

		if feeData.NewCreditDaysLeft >= 0 {
			s.log.Debug("[%s] Setting new credit days: %d", feeData.Login, feeData.NewCreditDaysLeft)
			if !dryRun {
				err := s.customerRepo.SetCreditDaysLeft(feeData.FeeLog.CustomerId.UUID.String(), feeData.NewCreditDaysLeft)
				if err != nil {
					s.log.Error("[%s] Can't set CreditDaysLeft: %v", feeData.Login, err)
				}
			}
		}

		if feeData.ServiceInternetState != feeData.NewState {
			if maxBlocks > 0 && blocks >= maxBlocks {
				s.log.Debug("[%s] Max blocks reached: %d/%d", feeData.Login, blocks, maxBlocks)
			} else {
				s.log.Debug("[%s] Setting ServiceInternetState: %d", feeData.Login, feeData.NewState)
				blocks++
				if !dryRun {
					err := s.customerRepo.SetServiceInternetState(feeData.FeeLog.CustomerId.UUID.String(), feeData.NewState)
					if err != nil {
						s.log.Error("can't set state: %v", err)
					}
				}
			}
		}
	}

	return feeDatas, nil
}

func (s *Fees) NonStrictLte(a, b float64) bool {
	sA := fmt.Sprintf("%0.2f", a)
	sB := fmt.Sprintf("%0.2f", b)
	if sA == sB {
		return true
	}
	fA, err := strconv.ParseFloat(sA, 64)
	if err != nil {
		return false
	}
	fB, err := strconv.ParseFloat(sB, 64)
	if err != nil {
		return false
	}
	if fA <= fB {
		return true
	}
	return false
}

func (s *Fees) CreateFee(c models2.CustomerData, services []domain.FeeService, descr string, forTime time.Time, ignoreCreditExpire bool) (*domain.Fee, error) {

	feeLog := models2.FeeLogData{}
	newDeposit := c.Deposit
	decreaseCreditDays := false

	for _, service := range services {

		fee := 0.0
		decreaseDeposit := true

		switch serv := service.Service.(type) {
		case models2.ServiceInternetData:
			fee = service.Fee
			if s.NonStrictLte(c.Deposit, 0.00) {
				s.log.Debug("[%s] Negative deposit: %0.2f", c.Login, c.Deposit)
				if !c.CreditExpire.Time.IsZero() {
					s.log.Debug("[%s] CreditExpire is set: %v", c.Login, c.CreditExpire)
					if forTime.After(c.CreditExpire.Time) {
						s.log.Debug("[%s] Credit expired, disabling service", c.Login)
						c.ServiceInternetState = models2.ServiceStateDisabled
						decreaseDeposit = false
					}
				} else if c.CreditDaysLeft > 0 {
					c.CreditDaysLeft--
					decreaseCreditDays = true
					s.log.Debug("[%s] CreditDaysLeft is set, decreasing: %d", c.Login, c.CreditDaysLeft)
				} else {
					creditExceeded := false
					if s.NonStrictLte(c.Deposit, -c.Credit) {
						creditExceeded = true
					}

					if creditExceeded {
						s.log.Debug("[%s] Credit exceeded: %0.2f/%d/%v, disabling service", c.Login, c.Credit, c.CreditDaysLeft, c.CreditExpire)
						c.ServiceInternetState = models2.ServiceStateDisabled
						decreaseDeposit = false
					}
				}
			}
			if decreaseDeposit {
				s.log.Debug("[%s] Decreasing deposit", c.Login)
				feeLog.ServiceInternetId = models2.NewNullUUID(serv.Id)
			}
		}

		if decreaseDeposit {
			feeLog.FeeAmount += fee
			newDeposit -= fee
		}
	}

	feeLog.CustomerId = models2.NewNullUUID(c.Id)
	feeLog.FeeTimestamp = forTime
	feeLog.PrevDeposit = float64(int(c.Deposit*100)) / 100
	feeLog.NewDeposit = float64(int(newDeposit*100)) / 100
	feeLog.Credit = c.Credit
	feeLog.Descr = descr
	if !c.CreditExpire.Time.IsZero() {
		feeLog.CreditExpire.Time = c.CreditExpire.Time
		feeLog.CreditExpire.Valid = true
	}

	feeLog.CreateTs = forTime

	feeData := &domain.Fee{
		FeeLog:     feeLog,
		NewDeposit: newDeposit,
		NewState:   c.ServiceInternetState,
	}

	if decreaseCreditDays {
		feeData.NewCreditDaysLeft = c.CreditDaysLeft
	} else {
		feeData.NewCreditDaysLeft = -1
	}

	return feeData, nil
}

func (s *Fees) DaysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func NewFees(
	log *zero.Logger,
	feeRepo repository2.Fees,
	customerRepo repository2.Customers,
	serviceInternetRepo repository2.ServiceInternet,
	ipPoolRepo repository2.IpPools,
) *Fees {
	return &Fees{
		log:                 log,
		feeRepo:             feeRepo,
		customerRepo:        customerRepo,
		serviceInternetRepo: serviceInternetRepo,
		ipPoolRepo:          ipPoolRepo,
	}
}

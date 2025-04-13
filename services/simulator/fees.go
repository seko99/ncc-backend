package simulator

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/fees"
	"fmt"
	"time"
)

func (ths *Simulator) createFees(startTime time.Time) error {
	ths.log.Info("Creating fees starting from %v", startTime)

	feesService := fees.NewFees(ths.log, ths.feesRepo, ths.customerRepo, ths.serviceInternetRepo, ths.ipPoolRepo)

	internets, err := ths.serviceInternetRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get serviceInternet: %w", err)
	}
	internetMap := map[string]models.ServiceInternetData{}
	for _, i := range internets {
		internetMap[i.Id] = i
	}

	currentTime := startTime
	for {
		ths.log.Info("Creating fees for %v...", currentTime)
		todayFees, err := ths.feesRepo.GetProcessedMap(currentTime)
		if err != nil {
			return fmt.Errorf("can't get processed fees: %w", err)
		}

		customDataMap, err := ths.serviceInternetRepo.GetCustomDataMap()
		if err != nil {
			return fmt.Errorf("can't get custom map: %w", err)
		}

		customers, err := ths.customerRepo.Get()
		if err != nil {
			return fmt.Errorf("can't get customers: %w", err)
		}

		days := feesService.DaysIn(time.Now().Month(), time.Now().Year())
		_, err = feesService.Process(internetMap, customers, customDataMap, todayFees, days, currentTime, false, 0, false, false)
		if err != nil {
			ths.log.Error("Can't process fees: %v", err)
		}

		currentTime = currentTime.Add(24 * time.Hour)
		if currentTime.After(time.Now()) {
			break
		}
	}

	return nil
}

package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
)

type LeasesUpdateUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *LeasesUpdateUsecase) Execute() error {
	ths.log.Info("Updating leases...")

	err := ths.simulator.UpdateLeases()
	if err != nil {
		return err
	}

	return nil
}

func NewLeasesUpdateUsecase(log logger.Logger, simulator interfaces.Simulator) LeasesUpdateUsecase {
	return LeasesUpdateUsecase{
		log:       log,
		simulator: simulator,
	}
}

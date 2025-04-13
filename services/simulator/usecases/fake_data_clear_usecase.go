package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
)

type FakeDataClearUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *FakeDataClearUsecase) Execute(req dto.FakeDataClearUsecaseRequest) error {
	ths.log.Info("Clearing fake data...")

	err := ths.simulator.ClearFakeData(req)

	if err != nil {
		ths.log.Error("Can't clear: %v", err)
		return err
	}

	return nil
}

func NewFakeDataClearUsecase(log logger.Logger, simulator interfaces.Simulator) FakeDataClearUsecase {
	return FakeDataClearUsecase{
		log:       log,
		simulator: simulator,
	}
}

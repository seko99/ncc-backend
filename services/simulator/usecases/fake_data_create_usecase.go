package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
)

type FakeDataCreateUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *FakeDataCreateUsecase) Execute(req dto.FakeDataCreateUsecaseRequest) error {
	ths.log.Info("Creating fake data...")

	go func() {
		err := ths.simulator.CreateFakeData(req)

		if err != nil {
			ths.log.Error("Can't create: %v", err)
			return
		}
	}()

	return nil
}

func NewFakeDataCreateUsecase(log logger.Logger, simulator interfaces.Simulator) FakeDataCreateUsecase {
	return FakeDataCreateUsecase{
		log:       log,
		simulator: simulator,
	}
}

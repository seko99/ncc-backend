package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
)

type InitDictionariesUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *InitDictionariesUsecase) Execute() error {
	ths.log.Info("Init dictionaries...")

	err := ths.simulator.InitDictionaries()
	if err != nil {
		ths.log.Error("Can't init: %v", err)
		return err
	}

	return nil
}

func NewInitDictionariesUsecase(log logger.Logger, simulator interfaces.Simulator) InitDictionariesUsecase {
	return InitDictionariesUsecase{
		log:       log,
		simulator: simulator,
	}
}

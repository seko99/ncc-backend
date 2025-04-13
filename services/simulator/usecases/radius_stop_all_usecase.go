package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"fmt"
)

type RadiusStopAllUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *RadiusStopAllUsecase) Execute(req dto.RadiusUsecaseRequest) error {
	err := ths.simulator.StopRadiusSessions(req)
	if err != nil {
		return fmt.Errorf("can't stop radius sessions: %w", err)
	}
	return nil
}

func NewRadiusStopAllUsecase(log logger.Logger, simulator interfaces.Simulator) RadiusStopAllUsecase {
	return RadiusStopAllUsecase{
		log:       log,
		simulator: simulator,
	}
}

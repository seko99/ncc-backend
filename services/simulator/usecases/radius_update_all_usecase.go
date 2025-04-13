package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"fmt"
)

type RadiusUpdateAllUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *RadiusUpdateAllUsecase) Execute(req dto.RadiusUsecaseRequest) error {
	err := ths.simulator.UpdateRadiusSessions(req)
	if err != nil {
		return fmt.Errorf("can't update radius sessions: %w", err)
	}
	return nil
}

func NewRadiusUpdateAllUsecase(log logger.Logger, simulator interfaces.Simulator) RadiusUpdateAllUsecase {
	return RadiusUpdateAllUsecase{
		log:       log,
		simulator: simulator,
	}
}

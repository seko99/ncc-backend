package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"fmt"
)

type RadiusStartAllUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *RadiusStartAllUsecase) Execute(req dto.RadiusUsecaseRequest) (dto.RadiusUsecaseResponse, error) {
	response, err := ths.simulator.StartRadiusSessions(req)
	if err != nil {
		return response, fmt.Errorf("can't start radius sessions: %w", err)
	}
	return response, nil
}

func NewRadiusStartAllUsecase(log logger.Logger, simulator interfaces.Simulator) RadiusStartAllUsecase {
	return RadiusStartAllUsecase{
		log:       log,
		simulator: simulator,
	}
}

package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"fmt"
)

type RadiusKillSessionsUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *RadiusKillSessionsUsecase) Execute(req dto.RadiusKillSessionsUsecaseRequest) (dto.RadiusKillSessionsUsecaseResponse, error) {
	response, err := ths.simulator.KillRadiusSessions(req)
	if err != nil {
		return dto.RadiusKillSessionsUsecaseResponse{}, fmt.Errorf("can't kill radius sessions: %w", err)
	}
	return response, nil
}

func NewRadiusKillSessionsUsecase(log logger.Logger, simulator interfaces.Simulator) RadiusKillSessionsUsecase {
	return RadiusKillSessionsUsecase{
		log:       log,
		simulator: simulator,
	}
}

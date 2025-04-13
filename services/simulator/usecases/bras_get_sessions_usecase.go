package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"fmt"
)

type BrasGetSessionsUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *BrasGetSessionsUsecase) Execute(req dto.BrasGetSessionsUsecaseRequest) (dto.BrasGetSessionsUsecaseResponse, error) {
	cache, err := ths.simulator.GetSessionCache()
	if err != nil {
		return dto.BrasGetSessionsUsecaseResponse{}, fmt.Errorf("can't get session: %w", err)
	}

	return dto.BrasGetSessionsUsecaseResponse{
		Sessions: cache,
	}, nil
}

func NewBrasGetSessionsUsecase(log logger.Logger, simulator interfaces.Simulator) BrasGetSessionsUsecase {
	return BrasGetSessionsUsecase{
		log:       log,
		simulator: simulator,
	}
}

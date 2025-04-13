package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
)

type BrasSetParamsUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *BrasSetParamsUsecase) Execute(req dto.BrasSetParamsUsecaseRequest) (dto.BrasSetParamsUsecaseResponse, error) {
	ths.simulator.SetBRASParams(simulator.BrasParams{
		SendInterims: req.SendInterims,
	})
	return dto.BrasSetParamsUsecaseResponse{}, nil
}

func NewBrasSetParamsUsecase(log logger.Logger, simulator interfaces.Simulator) BrasSetParamsUsecase {
	return BrasSetParamsUsecase{
		log:       log,
		simulator: simulator,
	}
}

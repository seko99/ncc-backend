package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
)

type IssuesCreateUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *IssuesCreateUsecase) Execute() error {
	err := ths.simulator.CreateIssues()
	if err != nil {
		return err
	}

	return nil
}

func NewIssuesCreateUsecase(log logger.Logger, simulator interfaces.Simulator) IssuesCreateUsecase {
	return IssuesCreateUsecase{
		log:       log,
		simulator: simulator,
	}
}

package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
)

type IssuesDeleteAllUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *IssuesDeleteAllUsecase) Execute() error {
	err := ths.simulator.DeleteIssues()
	if err != nil {
		return err
	}

	return nil
}

func NewIssuesDeleteAllUsecase(log logger.Logger, simulator interfaces.Simulator) IssuesDeleteAllUsecase {
	return IssuesDeleteAllUsecase{
		log:       log,
		simulator: simulator,
	}
}

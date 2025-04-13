package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
)

type SessionsUpdateUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *SessionsUpdateUsecase) Execute() error {
	ths.log.Info("Updating sessions...")

	err := ths.simulator.UpdateSessions()
	if err != nil {
		return err
	}

	return nil
}

func NewSessionsUpdateUsecase(log logger.Logger, simulator interfaces.Simulator) SessionsUpdateUsecase {
	return SessionsUpdateUsecase{
		log:       log,
		simulator: simulator,
	}
}

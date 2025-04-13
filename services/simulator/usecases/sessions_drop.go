package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"fmt"
)

type SessionsDropUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *SessionsDropUsecase) Execute() error {
	ths.log.Info("Dropping sessions...")

	err := ths.simulator.DropSessions()
	if err != nil {
		return fmt.Errorf("can't drop sessions: %w", err)
	}

	return nil
}

func NewSessionsDropUsecase(log logger.Logger, simulator interfaces.Simulator) SessionsDropUsecase {
	return SessionsDropUsecase{
		log:       log,
		simulator: simulator,
	}
}

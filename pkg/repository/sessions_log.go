package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_sessions_log_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository SessionsLog

type SessionsLog interface {
	Create(session models.SessionsLogData) error
	GetByCustomer(id string, period TimePeriod, limit ...int) ([]models.SessionsLogData, error)
	GetBySessionId(sessionId string) (models.SessionsLogData, error)
	DeleteById(id string) error
}

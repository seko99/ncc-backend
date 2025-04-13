package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_sessions_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Sessions

const (
	SessionCreatedEvent    = "sessionCreated"
	SessionUpdatedEvent    = "sessionUpdated"
	SessionDeletedEvent    = "sessionDeleted"
	SessionAllDeletedEvent = "sessionAllDeleted"
)

type Sessions interface {
	Get() ([]models.SessionData, error)
	GetById(id string) (models.SessionData, error)
	GetBySessionId(sessionId string) (models.SessionData, error)
	GetByIP(ip string) (models.SessionData, error)
	GetByLogin(login string) ([]models.SessionData, error)

	Create(data []models.SessionData) error
	Update(data models.SessionData) error
	UpdateBySessionId(data models.SessionData) error
	GetByCustomer(id string, period TimePeriod, limit ...int) ([]models.SessionData, error)
	Delete(id string) error
	DeleteAll() error
	DeleteBySessionId(data models.SessionData) error
}

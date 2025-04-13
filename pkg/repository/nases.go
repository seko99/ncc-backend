package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_nases_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Nases

const (
	NASCreatedEvent    = "nasCreated"
	NASUpdatedEvent    = "nasUpdated"
	NASDeletedEvent    = "nasDeleted"
	NASAllDeletedEvent = "nasAllDeleted"
)

type Nases interface {
	Get() ([]models.NasData, error)
	GetByIP(ip string) (models.NasData, error)

	Create(data []models.NasData) error
	Update(data models.NasData) error
	Delete(id string) error
	DeleteAll() error
	Upsert(data models.NasData) error
}

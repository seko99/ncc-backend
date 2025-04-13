package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_streets_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Streets

type Streets interface {
	Create(data models.StreetData) error
	Upsert(data models.StreetData) error
	Update(data models.StreetData) error
	Delete(id string) error
	Get() ([]models.StreetData, error)
}

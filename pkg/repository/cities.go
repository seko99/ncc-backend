package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_cities_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Cities

type Cities interface {
	Create(data models.CityData) error
	Upsert(data models.CityData) error
	Update(data models.CityData) error
	Delete(id string) error
	Get() ([]models.CityData, error)
}

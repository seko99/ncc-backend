package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_sorm_customers_errors_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository SormCustomersErrors

type SormCustomersErrors interface {
	Create(data []models.SormCustomersErrorsData) error
	Upsert(data []models.SormCustomersErrorsData) error
	Update(data models.SormCustomersErrorsData) error
	Delete(id string) error
	DeleteAll() error
	Get() ([]models.SormCustomersErrorsData, error)
}

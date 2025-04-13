package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_sorm_customers_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository SormCustomers

type SormCustomers interface {
	Create(data []models.SormCustomersData) error
	Upsert(data []models.SormCustomersData) error
	Update(data models.SormCustomersData) error
	Delete(id string) error
	Get() ([]models.SormCustomersData, error)
}

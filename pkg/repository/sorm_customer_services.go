package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_sorm_customer_services_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository SormCustomerServices

type SormCustomerServices interface {
	Create(data []models.SormCustomerServiceData) error
	Upsert(data []models.SormCustomerServiceData) error
	Update(data models.SormCustomerServiceData) error
	Delete(id string) error
	Get() ([]models.SormCustomerServiceData, error)
}

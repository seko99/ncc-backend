package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_sorm_customer_services_errors_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository SormCustomerServicesErrors

type SormCustomerServicesErrors interface {
	Create(data []models.SormCustomerServicesErrorsData) error
	Upsert(data []models.SormCustomerServicesErrorsData) error
	Update(data models.SormCustomerServicesErrorsData) error
	Delete(id string) error
	DeleteAll() error
	Get() ([]models.SormCustomerServicesErrorsData, error)
}

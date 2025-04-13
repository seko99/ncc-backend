package repository

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_service_internet_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository ServiceInternet

type ServiceInternet interface {
	Get() ([]models2.ServiceInternetData, error)
	GetById(id string) (*models2.ServiceInternetData, error)
	Create(service models2.ServiceInternetData) error
	Upsert(service models2.ServiceInternetData) error
	GetCustomDataByCustomer(customer models2.CustomerData) (*models2.ServiceInternetCustomData, error)
	GetCustomDataMap() (map[string]models2.ServiceInternetCustomData, error)
}

package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_informings_test_customers_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository InformingsTestCustomers

type InformingsTestCustomers interface {
	Get() ([]models.InformingTestCustomerData, error)
}

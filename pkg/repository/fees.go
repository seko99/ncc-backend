package repository

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"time"
)

//go:generate mockgen -destination=mocks/mock_fees_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Fees

type Fees interface {
	Create(fee models2.FeeLogData) error
	Update(fee models2.FeeLogData) error
	CreateBatch(fees []models2.FeeLogData) error
	GetByCustomer(id string, period TimePeriod, limit ...int) ([]models2.FeeLogData, error)
	Get(period TimePeriod) ([]models2.FeeLogData, error)
	GetProcessed(t ...time.Time) ([]models2.CustomerData, error)
	GetProcessedMap(t ...time.Time) (map[string]models2.CustomerData, error)
}

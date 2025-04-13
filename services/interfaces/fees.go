package interfaces

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/domain"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"time"
)

//go:generate mockgen -destination=mocks/mock_fees_service.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/interfaces FeesProcessor

type FeesProcessor interface {
	Process(customers []models.CustomerData, todayFees map[string]models.CustomerData, days int, forTime time.Time) (map[string]domain.Fee, error)
	CreateFee(c models.CustomerData, services []domain.FeeService, descr string, forTime time.Time) (*domain.Fee, error)
	DaysIn(m time.Month, year int) int
}

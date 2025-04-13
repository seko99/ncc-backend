package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_payments_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Payments

type Payments interface {
	Create(data models.PaymentData) error
	GetPaymentsByCustomer(id string, period TimePeriod) ([]models.PaymentData, error)
	GetPayments(period TimePeriod) ([]models.PaymentData, error)
}

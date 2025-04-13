package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_payment_types_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository PaymentTypes

type PaymentTypes interface {
	Create(data models.PaymentTypeData) error
	Upsert(data models.PaymentTypeData) error
	Update(data models.PaymentTypeData) error
	Delete(id string) error
	Get() ([]models.PaymentTypeData, error)
}

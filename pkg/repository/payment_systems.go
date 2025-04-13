package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_payment_systems_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository PaymentSystems

type PaymentSystems interface {
	Create(data models.PaymentSystemData) error
	Upsert(data models.PaymentSystemData) error
	Update(data models.PaymentSystemData) error
	Delete(id string) error
	Get() ([]models.PaymentSystemData, error)
}
